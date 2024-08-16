package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/brettbuddin/git-credential-op/internal/subcmd"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
)

func main() {
	err := run(os.Args[1:])
	if err == nil {
		return
	}
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func run(args []string) error {
	var (
		titleFormat string
		locatorTag  string
		account     string
		vault       string
	)
	fs := flag.NewFlagSet("git-credential-op", flag.ExitOnError)
	fs.StringVar(&titleFormat, "title", subcmd.DefaultTitleFormat, "item title format")
	fs.StringVar(&locatorTag, "locator-tag", op.DefaultLocatorTag, "locator tag value")
	fs.StringVar(&account, "account", "", "account URL (e.g. mycompany.1password.com")
	fs.StringVar(&vault, "vault", "", "vault name (e.g. Private)")
	fs.Usage = usageFn(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	args = fs.Args()
	if len(args) != 1 {
		fs.Usage()
		return flag.ErrHelp
	}

	if !op.InPath() {
		return fmt.Errorf(`op (1Password) not found in PATH`)
	}

	runner := &op.Runner{
		LocatorTag: locatorTag,
		Account:    account,
		Vault:      vault,
		Executor:   op.DefaultExecutor(),
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}
	switch args[0] {
	case "get":
		return subcmd.Get(runner)
	case "store":
		return subcmd.Store(runner, titleFormat)
	case "erase":
		return subcmd.Erase(runner)
	default:
		// To fully conform with the git-credential contract we must ignore any subcommands we don't
		// recognize. See https://git-scm.com/docs/gitcredentials#_custom_helpers
		return nil
	}
}

var usage = `git-credential-op is a git-credential helper for storing credentials in 1Password.

Usage:

	git-credential-op (get|store|erase)

Setup with Git:

	# .gitconfig OR .git/config
	[credential]
            helper = op
`

func usageFn(fs *flag.FlagSet) func() {
	return func() {
		out := fs.Output()
		fmt.Fprintln(out, usage)
		fmt.Fprint(out, "Flags:\n\n")
		fs.PrintDefaults()
		fmt.Fprintln(out)
	}
}
