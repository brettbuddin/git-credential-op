package op

import (
	"io"
	"os/exec"
)

const op = "op"

// InPath checks that 1Password command line tool "op" is in the path.
func InPath() bool {
	_, err := exec.LookPath(op)
	return err == nil
}

func (r *Runner) execItemCommand(out ExecutorOutput, subcmd string, extraArgs ...string) error {
	args := r.itemCmd(subcmd, extraArgs...)
	if err := r.Executor.ExecuteCommand(out, op, args...); err != nil {
		return err
	}
	return nil
}

func (r *Runner) execOutput(stdout io.Writer) ExecutorOutput {
	return ExecutorOutput{
		Stdout: stdout,
		Stderr: r.Stderr,
	}
}

func (r *Runner) itemCmd(subcmd string, extraArgs ...string) []string {
	args := []string{"item", subcmd}
	if r.Account != "" {
		args = append(args, "--account", r.Account)
	}
	if r.Vault != "" {
		args = append(args, "--vault", r.Vault)
	}
	return append(args, extraArgs...)
}
