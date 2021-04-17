package jsonexec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// Run run `name arg...` command and treats its stdout as a JSON,
// unmarshaling it into provided dest object.
// Returns ErrorExecution if execution failed and ErrorUnmarshal if
// JSON unmarshal is failed.
func Run(dest interface{}, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ErrorExecution{
			errorWithOutput{
				err:    fmt.Errorf("run command: %w", err),
				output: stderr.String(),
			},
		}
	}

	if err := json.Unmarshal(stdout.Bytes(), dest); err != nil {
		return ErrorUnmarshal{
			errorWithOutput{
				err:    fmt.Errorf("unmarshal command output: %w", err),
				output: stdout.String(),
			},
		}
	}

	return nil
}

type errorWithOutput struct {
	err    error
	output string
}

func (e errorWithOutput) Error() string {
	return e.err.Error()
}

// Details error details
func (e errorWithOutput) Details() string {
	return e.output
}

// ErrorExecution ошибка исполнения команды
type ErrorExecution struct {
	errorWithOutput
}

// ErrorUnmarshal ошибка декодирования выхлопа команды трактуемого как JSON
type ErrorUnmarshal struct {
	errorWithOutput
}

// HandleError jsonexec internal unified handling
func HandleError(err error, details func(string)) {
	var errExec ErrorExecution
	if errors.As(err, &errExec) {
		details(errExec.output)
		return
	}

	var errUnmarshal ErrorUnmarshal
	if errors.As(err, &errUnmarshal) {
		details(errUnmarshal.output)
		return
	}
}
