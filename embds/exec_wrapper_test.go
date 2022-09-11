package embds

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_execWrapper_run(t *testing.T) {
	tests := []struct {
		name       string
		cmd        string
		args       []string
		timeout    time.Duration
		initErr    bool
		runtimeErr bool
	}{
		{name: "sleep 1 (ok)", cmd: "sleep", args: []string{"1"}, timeout: 1010 * time.Millisecond, runtimeErr: false},
		{name: "sleep 1 (killed by timeout)", cmd: "sleep", args: []string{"1"}, timeout: 990 * time.Millisecond, runtimeErr: true},
		{name: "wrong command name", cmd: "nonExistingCmd", initErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Init command
			ew, err := newExecWrapper(tt.cmd)
			if tt.initErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			ew.enableStderrForwarding()

			// Run command
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			err = ew.execute(ctx, tt.args...)
			if tt.runtimeErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
