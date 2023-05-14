package jsonexec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirkon/errors"
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
		out := strings.TrimSpace(stderr.String())
		if out != "" {
			out = "run command"
		}

		return errors.New(out).Str("error-status", err.Error())
	}

	if err := json.Unmarshal(stdout.Bytes(), dest); err != nil {
		return errors.Wrap(err, "unmarshal command output").Str("error-output", stdout.String())
	}

	return nil
}

// Rung generic version of Run. Can be a bit more convenient.
func Rung[T any](name string, args ...string) (res T, err error) {
	if err := Run(&res, name, args...); err != nil {
		return res, err
	}

	return res, nil
}

type errorWithOutput struct {
	err    error
	output string
}

func (e errorWithOutput) Error() string {
	return fmt.Sprintf("%s: %s", e.output, e.err)
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
