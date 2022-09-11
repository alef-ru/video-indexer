package embds

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type execWrapper struct {
	executablePath string
	forwardStderr  bool
}

func (ew *execWrapper) execute(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, ew.executablePath, args...)
	if ew.forwardStderr {
		return runForwardingStderr(cmd)
	} else {
		return cmd.Run()
	}
}

func runForwardingStderr(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start ffmpeg : %w", err)
	}

	buf := make([]byte, 8)
	for {
		n, err := stderr.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read stderr : %w", err)
		}
		_, err = fmt.Fprint(os.Stderr, string(buf[:n]))
		if err != nil {
			return fmt.Errorf("failed to write stderr : %w", err)
		}
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func newExecWrapper(cmdName string) (*execWrapper, error) {
	path, err := exec.LookPath(cmdName)
	if err != nil {
		return nil, fmt.Errorf("command %s not found: %w ", cmdName, err)
	}
	log.Printf("INFO %s is available at %s", cmdName, path)
	cmd := execWrapper{
		executablePath: path,
		forwardStderr:  false,
	}
	return &cmd, nil
}

func (ew *execWrapper) enableStderrForwarding() {
	ew.forwardStderr = true
}

func (ew *execWrapper) disableStderrForwarding() {
	ew.forwardStderr = false
}
