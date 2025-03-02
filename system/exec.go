package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tuxounet/k2-sdk/types"
)

type CmdCall struct {
	Command string
	Args    []string
	Cwd     *string
	log     types.ILogger
}

func NewCmdCall(log types.ILogger, command string, args ...string) *CmdCall {

	return &CmdCall{
		Command: command,
		Args:    args,
		Cwd:     nil,
		log:     log,
	}
}

func (c *CmdCall) String() string {

	out := c.Command + " "
	for _, arg := range c.Args {
		out += arg + " "
	}
	out = strings.TrimSpace(out)

	return out
}

func OsExec(query *CmdCall) error {

	code, err := OsExecWithExitCode(query)
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New("command failed with exit code " + fmt.Sprint(code))
	}
	return nil

}

func OsExecWithExitCode(query *CmdCall) (int, error) {
	query.log.TraceF("executing command: %s", query.String())

	cmd := exec.Command(query.Command, query.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if query.Cwd != nil {
		cmd.Dir = *query.Cwd
	}

	err := cmd.Start()
	if err != nil {
		query.log.ErrorF("command %s failed to start command: %s", query.String(), err)
		return -1, err
	}

	cmd.Wait()
	status := cmd.ProcessState.ExitCode()
	query.log.TraceF("command %s exited with status %d", query.String(), status)

	return status, nil

}

type cmdExecTail struct {
	log    types.ILogger
	stream string
	buffer []byte
}

func (c *cmdExecTail) Write(p []byte) (n int, err error) {
	c.buffer = append(c.buffer, p...)
	log := c.log.CreateSubLogger(c.stream)
	log.Trace(string(p))
	return len(p), nil
}

func OsExecAndTailToLog(query *CmdCall) ([]byte, error) {

	query.log.TraceF("executing command: %s", query.String())

	cmd := exec.Command(query.Command, query.Args...)

	strErrTailer := &cmdExecTail{log: query.log, stream: "stderr"}
	stdOutTailer := &cmdExecTail{log: query.log, stream: "stdout"}
	cmd.Stderr = strErrTailer
	cmd.Stdout = stdOutTailer

	if query.Cwd != nil {
		cmd.Dir = *query.Cwd
	}

	err := cmd.Start()
	if err != nil {
		query.log.ErrorF("command %s failed to start command: %s", query.String(), err)
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		query.log.ErrorF("command %s failed to complete: %s", query.String(), err)
		return nil, err
	}

	return stdOutTailer.buffer, nil

}

type cmdStartTail struct {
	log    types.ILogger
	stream string
}

func (c *cmdStartTail) Write(p []byte) (n int, err error) {
	log := c.log.CreateSubLogger(c.stream)
	log.Trace(string(p))
	return len(p), nil
}

func OsStartAndTailToLog(query *CmdCall) (*exec.Cmd, error) {

	query.log.TraceF("executing command: %s", query.String())

	cmd := exec.Command(query.Command, query.Args...)
	strErrTailer := &cmdStartTail{log: query.log, stream: "stderr"}
	stdOutTailer := &cmdStartTail{log: query.log, stream: "stdout"}
	cmd.Stderr = strErrTailer
	cmd.Stdout = stdOutTailer

	if query.Cwd != nil {
		cmd.Dir = *query.Cwd
	}

	err := cmd.Start()
	if err != nil {
		query.log.ErrorF("command %s failed to start command: %s", query.String(), err)
		return nil, err
	}

	return cmd, nil

}

func JsonCommandOutput[R any](query *CmdCall) (R, error) {
	var result R
	out, err := RawCommandOutput(query)
	if err != nil {
		return result, err
	}

	result, err = LoadJSONFromString[R](out)
	if err != nil {
		return result, err
	}

	return result, nil

}

func RawCommandOutput(query *CmdCall) (string, error) {
	query.log.TraceF("executing raw command: %s", query.String())

	cmd := exec.Command(query.Command, query.Args...)
	if query.Cwd != nil {
		cmd.Dir = *query.Cwd
	}
	ret, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output := string(ret)
	query.log.TraceF("command %s returned: %s", query.String(), output)

	return output, nil

}
