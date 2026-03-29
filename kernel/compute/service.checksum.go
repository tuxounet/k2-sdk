package compute

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tuxounet/k2-sdk/kernel/compute/types"
)

const checksumCachePath = "etc/compute/checksums.json"
const defaultChecksumTTLMinutes = 1440 // 24 hours

func (s *Service) getChecksumTTL() time.Duration {
	ttlMinutes, err := s.getConfigService().GetAsIntOrDefault("host.compute.cache.ttl", defaultChecksumTTLMinutes)
	if err != nil {
		return time.Duration(defaultChecksumTTLMinutes) * time.Minute
	}
	return time.Duration(ttlMinutes) * time.Minute
}

func (s *Service) computePlaybookChecksum(verb types.RunnerVerb) (string, error) {
	paths := s.getPathsService()
	playbookPath := paths.CominePath("etc", "compute", fmt.Sprintf("%s.yaml", verb))

	rootStore, err := s.getRootStore()
	if err != nil {
		return "", fmt.Errorf("failed to get root store: %s", err.Error())
	}

	playbookContent, err := rootStore.ReadObject(playbookPath)
	if err != nil {
		return "", fmt.Errorf("failed to read playbook %s: %s", verb, err.Error())
	}

	configMap := s.getConfigService().GetCurrent()
	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to serialize config: %s", err.Error())
	}

	hasher := sha256.New()
	hasher.Write(playbookContent)
	hasher.Write([]byte("||"))
	hasher.Write(configBytes)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (s *Service) loadChecksumCache() (types.ChecksumCache, error) {
	rootStore, err := s.getRootStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get root store: %s", err.Error())
	}

	exists, err := rootStore.Exists(checksumCachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check checksum cache existence: %s", err.Error())
	}
	if !exists {
		return make(types.ChecksumCache), nil
	}

	data, err := rootStore.ReadObject(checksumCachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checksum cache: %s", err.Error())
	}

	var cache types.ChecksumCache
	err = json.Unmarshal(data, &cache)
	if err != nil {
		s.GetLogger().WarnF("corrupted checksum cache, resetting: %s", err.Error())
		return make(types.ChecksumCache), nil
	}

	return cache, nil
}

func (s *Service) saveChecksumCache(cache types.ChecksumCache) error {
	rootStore, err := s.getRootStore()
	if err != nil {
		return fmt.Errorf("failed to get root store: %s", err.Error())
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to serialize checksum cache: %s", err.Error())
	}

	err = rootStore.WriteObject(checksumCachePath, data)
	if err != nil {
		return fmt.Errorf("failed to write checksum cache: %s", err.Error())
	}

	return nil
}

func (s *Service) shouldExecPlaybook(verb types.RunnerVerb, force bool) (bool, error) {
	if verb == types.RunnerVerbTeardown || verb == types.RunnerVerbStop {
		s.GetLogger().DebugF("%s verb always executes (no cache)", verb)
		return true, nil
	}

	if force {
		s.GetLogger().DebugF("force flag set, will execute %s", verb)
		return true, nil
	}

	currentChecksum, err := s.computePlaybookChecksum(verb)
	if err != nil {
		s.GetLogger().WarnF("failed to compute checksum for %s, will execute: %s", verb, err.Error())
		return true, nil
	}

	cache, err := s.loadChecksumCache()
	if err != nil {
		s.GetLogger().WarnF("failed to load checksum cache, will execute: %s", err.Error())
		return true, nil
	}

	entry, exists := cache[string(verb)]
	if !exists {
		s.GetLogger().DebugF("no cached checksum for %s, will execute", verb)
		return true, nil
	}

	if entry.Checksum != currentChecksum {
		s.GetLogger().DebugF("checksum changed for %s (cached=%s, current=%s), will execute", verb, entry.Checksum[:8], currentChecksum[:8])
		return true, nil
	}

	elapsed := time.Since(entry.ExecutedAt)
	ttl := s.getChecksumTTL()
	if elapsed > ttl {
		s.GetLogger().DebugF("checksum TTL expired for %s (last run: %s ago), will execute", verb, elapsed.Round(time.Minute))
		return true, nil
	}

	s.GetLogger().InfoF("[CACHED] skipping %s playbook: no changes detected since last run %s ago (checksum: %s)", verb, elapsed.Round(time.Minute), currentChecksum[:8])
	return false, nil
}

func (s *Service) markPlaybookExecuted(verb types.RunnerVerb) error {
	currentChecksum, err := s.computePlaybookChecksum(verb)
	if err != nil {
		return fmt.Errorf("failed to compute checksum for %s: %s", verb, err.Error())
	}

	cache, err := s.loadChecksumCache()
	if err != nil {
		cache = make(types.ChecksumCache)
	}

	cache[string(verb)] = types.VerbChecksum{
		Verb:       verb,
		Checksum:   currentChecksum,
		ExecutedAt: time.Now(),
	}

	return s.saveChecksumCache(cache)
}

func (s *Service) nukeChecksumCache() error {
	rootStore, err := s.getRootStore()
	if err != nil {
		return fmt.Errorf("failed to get root store: %s", err.Error())
	}

	exists, err := rootStore.Exists(checksumCachePath)
	if err != nil {
		return fmt.Errorf("failed to check checksum cache existence: %s", err.Error())
	}
	if !exists {
		return nil
	}

	err = rootStore.DeleteObject(checksumCachePath)
	if err != nil {
		return fmt.Errorf("failed to delete checksum cache: %s", err.Error())
	}

	s.GetLogger().DebugF("checksum cache deleted after teardown")
	return nil
}

func (s *Service) invalidateVerbCache(verb types.RunnerVerb) error {
	cache, err := s.loadChecksumCache()
	if err != nil {
		return fmt.Errorf("failed to load checksum cache: %s", err.Error())
	}

	if _, exists := cache[string(verb)]; !exists {
		return nil
	}

	delete(cache, string(verb))
	s.GetLogger().DebugF("invalidated %s cache", verb)

	return s.saveChecksumCache(cache)
}
