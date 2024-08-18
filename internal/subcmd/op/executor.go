package op

import (
	"io"
	"os/exec"
)

// DefaultExecutor returns a [Executor] that actually runs commands.
func DefaultExecutor() Executor {
	return ExecutorFunc(func(out ExecutorOutput, name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmd.Stdout = out.Stdout
		cmd.Stderr = out.Stderr
		return cmd.Run()
	})
}

//go:generate mockery --outpkg opmock --output opmock --with-expecter --name Executor
type Executor interface {
	ExecuteCommand(out ExecutorOutput, name string, args ...string) error
}

// ExecutorFunc is a function that implements [Executor].
type ExecutorFunc func(out ExecutorOutput, name string, args ...string) error

func (f ExecutorFunc) ExecuteCommand(out ExecutorOutput, name string, args ...string) error {
	return f(out, name, args...)
}

// ExecutorOutput holds standard IO streams.
type ExecutorOutput struct {
	Stdout io.Writer
	Stderr io.Writer
}
