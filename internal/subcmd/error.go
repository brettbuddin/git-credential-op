package subcmd

import "fmt"

func opError(err error) error {
	return fmt.Errorf("op: %w", err)
}

func stdinError(err error) error {
	return fmt.Errorf("read stdin: %w", err)
}

func stdoutError(err error) error {
	return fmt.Errorf("write stdout: %w", err)
}
