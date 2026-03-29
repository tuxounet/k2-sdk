package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
)

// writeTestPlaybooks writes minimal playbook YAML files for all 4 verbs
// so checksum computation can read them.
func writeTestPlaybooks(t *testing.T, svc *Service) {
	t.Helper()
	rootStore, err := svc.getRootStore()
	require.NoError(t, err)

	for _, verb := range []types.RunnerVerb{
		types.RunnerVerbProvision,
		types.RunnerVerbStart,
		types.RunnerVerbStop,
		types.RunnerVerbTeardown,
	} {
		paths := svc.getPathsService()
		path := paths.CominePath("etc", "compute", fmt.Sprintf("%s.yaml", verb))
		content := fmt.Sprintf("- name: %s\n  hosts: all\n  gather_facts: no\n  tasks: []\n", verb)
		err := rootStore.WriteObject(path, []byte(content))
		require.NoError(t, err)
	}
}

func TestChecksumCache_LoadEmpty(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	cache, err := svc.loadChecksumCache()
	require.NoError(t, err)
	assert.Empty(t, cache)
}

func TestChecksumCache_SaveAndLoad(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	cache := types.ChecksumCache{
		"provision": {
			Verb:       types.RunnerVerbProvision,
			Checksum:   "abc123",
			ExecutedAt: time.Now().Truncate(time.Second),
		},
	}

	err := svc.saveChecksumCache(cache)
	require.NoError(t, err)

	loaded, err := svc.loadChecksumCache()
	require.NoError(t, err)
	assert.Equal(t, cache["provision"].Checksum, loaded["provision"].Checksum)
	assert.Equal(t, cache["provision"].Verb, loaded["provision"].Verb)
}

func TestChecksumCache_CorruptedReturnsEmpty(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	rootStore, err := svc.getRootStore()
	require.NoError(t, err)
	err = rootStore.WriteObject(checksumCachePath, []byte("not valid json"))
	require.NoError(t, err)

	cache, err := svc.loadChecksumCache()
	require.NoError(t, err)
	assert.Empty(t, cache)
}

func TestComputePlaybookChecksum_Deterministic(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	checksum1, err := svc.computePlaybookChecksum(types.RunnerVerbProvision)
	require.NoError(t, err)
	assert.NotEmpty(t, checksum1)

	checksum2, err := svc.computePlaybookChecksum(types.RunnerVerbProvision)
	require.NoError(t, err)

	assert.Equal(t, checksum1, checksum2)
}

func TestComputePlaybookChecksum_DiffersByVerb(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	provisionChecksum, err := svc.computePlaybookChecksum(types.RunnerVerbProvision)
	require.NoError(t, err)

	startChecksum, err := svc.computePlaybookChecksum(types.RunnerVerbStart)
	require.NoError(t, err)

	// Different verbs have different playbook names so checksums should differ
	assert.NotEqual(t, provisionChecksum, startChecksum)
}

func TestShouldExecPlaybook_Force(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	should, err := svc.shouldExecPlaybook(types.RunnerVerbProvision, true)
	require.NoError(t, err)
	assert.True(t, should)
}

func TestShouldExecPlaybook_NoCache(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	should, err := svc.shouldExecPlaybook(types.RunnerVerbProvision, false)
	require.NoError(t, err)
	assert.True(t, should, "should execute when no cache exists")
}

func TestShouldExecPlaybook_CachedAndUnchanged(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	// Mark as executed
	err := svc.markPlaybookExecuted(types.RunnerVerbProvision)
	require.NoError(t, err)

	// Should now skip
	should, err := svc.shouldExecPlaybook(types.RunnerVerbProvision, false)
	require.NoError(t, err)
	assert.False(t, should, "should skip when cached and unchanged")
}

func TestShouldExecPlaybook_ExpiredTTL(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	// Compute current checksum
	currentChecksum, err := svc.computePlaybookChecksum(types.RunnerVerbProvision)
	require.NoError(t, err)

	// Manually write an expired cache entry
	cache := types.ChecksumCache{
		"provision": {
			Verb:       types.RunnerVerbProvision,
			Checksum:   currentChecksum,
			ExecutedAt: time.Now().Add(-25 * time.Hour),
		},
	}
	data, _ := json.Marshal(cache)
	rootStore, _ := svc.getRootStore()
	err = rootStore.WriteObject(checksumCachePath, data)
	require.NoError(t, err)

	should, err := svc.shouldExecPlaybook(types.RunnerVerbProvision, false)
	require.NoError(t, err)
	assert.True(t, should, "should execute when TTL expired")
}

func TestShouldExecPlaybook_ChecksumChanged(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	// Write a cache entry with a different checksum
	cache := types.ChecksumCache{
		"provision": {
			Verb:       types.RunnerVerbProvision,
			Checksum:   "outdated_checksum_value",
			ExecutedAt: time.Now(),
		},
	}
	data, _ := json.Marshal(cache)
	rootStore, _ := svc.getRootStore()
	err := rootStore.WriteObject(checksumCachePath, data)
	require.NoError(t, err)

	should, err := svc.shouldExecPlaybook(types.RunnerVerbProvision, false)
	require.NoError(t, err)
	assert.True(t, should, "should execute when checksum changed")
}

func TestMarkPlaybookExecuted(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	err := svc.markPlaybookExecuted(types.RunnerVerbProvision)
	require.NoError(t, err)

	cache, err := svc.loadChecksumCache()
	require.NoError(t, err)

	entry, exists := cache["provision"]
	assert.True(t, exists)
	assert.Equal(t, types.RunnerVerbProvision, entry.Verb)
	assert.NotEmpty(t, entry.Checksum)
	assert.WithinDuration(t, time.Now(), entry.ExecutedAt, 5*time.Second)
}

func TestExecVerb_Disabled(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, false)
	defer os.RemoveAll(tmpDir)

	err := svc.ExecVerb(types.RunnerVerbProvision)
	require.NoError(t, err)
}

func TestShouldExecPlaybook_TeardownAlwaysExecutes(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	// Mark teardown as executed
	err := svc.markPlaybookExecuted(types.RunnerVerbTeardown)
	require.NoError(t, err)

	// Teardown should always execute regardless of cache
	should, err := svc.shouldExecPlaybook(types.RunnerVerbTeardown, false)
	require.NoError(t, err)
	assert.True(t, should, "teardown should always execute")
}

func TestShouldExecPlaybook_StopAlwaysExecutes(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	err := svc.markPlaybookExecuted(types.RunnerVerbStop)
	require.NoError(t, err)

	should, err := svc.shouldExecPlaybook(types.RunnerVerbStop, false)
	require.NoError(t, err)
	assert.True(t, should, "stop should always execute")
}

func TestNukeChecksumCache(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, true)
	defer os.RemoveAll(tmpDir)

	writeTestPlaybooks(t, svc)

	// Create cache
	err := svc.markPlaybookExecuted(types.RunnerVerbProvision)
	require.NoError(t, err)

	// Verify cache exists
	cache, err := svc.loadChecksumCache()
	require.NoError(t, err)
	assert.NotEmpty(t, cache)

	// Nuke cache
	err = svc.nukeChecksumCache()
	require.NoError(t, err)

	// Verify cache is gone
	cache, err = svc.loadChecksumCache()
	require.NoError(t, err)
	assert.Empty(t, cache)
}
