package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/testutils"
)

func TestNewCmdCall(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", "hello", "world")

	assert.Equal(t, "echo", cmd.Command)
	assert.Equal(t, []string{"hello", "world"}, cmd.Args)
	assert.Nil(t, cmd.Cwd)
}

func TestCmdCall_String_WithArgs(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "git", "commit", "-m", "message")

	assert.Equal(t, "git commit -m message", cmd.String())
}

func TestCmdCall_String_NoArgs(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "ls")

	assert.Equal(t, "ls", cmd.String())
}

func TestCmdCall_String_SingleArg(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "cat", "/etc/hosts")

	assert.Equal(t, "cat /etc/hosts", cmd.String())
}

func TestCmdCall_Cwd(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "ls")
	cwd := "/tmp"
	cmd.Cwd = &cwd

	assert.Equal(t, "/tmp", *cmd.Cwd)
}

func TestOsExec_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "true")

	err := OsExec(cmd)
	assert.NoError(t, err)
}

func TestOsExec_Failure(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "false")

	err := OsExec(cmd)
	assert.Error(t, err)
}

func TestOsExecWithExitCode_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "true")

	code, err := OsExecWithExitCode(cmd)
	assert.NoError(t, err)
	assert.Equal(t, 0, code)
}

func TestOsExecWithExitCode_NonZero(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "false")

	code, err := OsExecWithExitCode(cmd)
	assert.NoError(t, err)
	assert.Equal(t, 1, code)
}

func TestOsExecWithExitCode_InvalidCommand(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "nonexistent_command_xyz")

	_, err := OsExecWithExitCode(cmd)
	assert.Error(t, err)
}

func TestOsExecAndTailToLog_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", "hello world")

	output, err := OsExecAndTailToLog(cmd)
	assert.NoError(t, err)
	assert.Contains(t, string(output), "hello world")
}

func TestOsExecAndTailToLog_Failure(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "nonexistent_command_xyz")

	_, err := OsExecAndTailToLog(cmd)
	assert.Error(t, err)
}

func TestOsExecAndTailToLog_WithCwd(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "pwd")
	cwd := "/tmp"
	cmd.Cwd = &cwd

	output, err := OsExecAndTailToLog(cmd)
	assert.NoError(t, err)
	assert.Contains(t, string(output), "/tmp")
}

func TestOsStartAndTailToLog_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", "started")

	process, err := OsStartAndTailToLog(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, process)

	process.Wait()
}

func TestOsStartAndTailToLog_InvalidCommand(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "nonexistent_command_xyz")

	_, err := OsStartAndTailToLog(cmd)
	assert.Error(t, err)
}

func TestRawCommandOutput_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", "raw output")

	output, err := RawCommandOutput(cmd)
	assert.NoError(t, err)
	assert.Contains(t, output, "raw output")
}

func TestRawCommandOutput_Failure(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "nonexistent_command_xyz")

	_, err := RawCommandOutput(cmd)
	assert.Error(t, err)
}

func TestRawCommandOutput_WithCwd(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "pwd")
	cwd := "/tmp"
	cmd.Cwd = &cwd

	output, err := RawCommandOutput(cmd)
	assert.NoError(t, err)
	assert.Contains(t, output, "/tmp")
}

func TestJsonCommandOutput_Success(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", `{"name":"test","value":42}`)

	type Result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	result, err := JsonCommandOutput[Result](cmd)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

func TestJsonCommandOutput_InvalidJSON(t *testing.T) {
	log := testutils.NewMockLogger("test")
	cmd := NewCmdCall(log, "echo", "not json")

	type Result struct {
		Name string `json:"name"`
	}

	_, err := JsonCommandOutput[Result](cmd)
	assert.Error(t, err)
}
