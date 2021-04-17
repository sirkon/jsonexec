// +build development

package jsonexec

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		dest     map[string]interface{}
		cmd      string
		arg      []string
		checkOut func(v map[string]interface{}) error
		wantErr  bool
	}{
		{
			name: "go list",
			cmd:  "go",
			arg:  []string{"list", "-m", "-json", "all"},
			checkOut: func(v map[string]interface{}) error {
				const curPath = "github.com/sirkon/jsonexec"
				if v["Path"] != curPath {
					return fmt.Errorf("path %s expected got %s", curPath, v["Path"])
				}

				return nil
			},
			wantErr: false,
		},
		{
			name:     "command execute error",
			cmd:      "aaaaaa",
			arg:      nil,
			checkOut: nil,
			wantErr:  true,
		},
		{
			name:     "not json output",
			cmd:      "ls",
			arg:      nil,
			checkOut: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(&tt.dest, tt.cmd, tt.arg...)
			if err != nil {
				t.Log(err)
				HandleError(err, func(d string) {
					t.Log(d)
				})
			}
			switch {
			case err == nil && !tt.wantErr:
				if err := tt.checkOut(tt.dest); err != nil {
					t.Errorf("check output: %s", err)
				}
			case err == nil && tt.wantErr:
				t.Errorf("missing expected error")
			case err != nil && !tt.wantErr:
				t.Errorf("unexpected error: %s", err)
			case err != nil && tt.wantErr:
			}
		})
	}
}
