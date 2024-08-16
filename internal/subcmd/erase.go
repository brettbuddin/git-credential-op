package subcmd

import (
	"errors"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
)

// Erase is the "erase" subcommand of the git-credential helper contract.
//
// See https://git-scm.com/docs/gitcredentials#Documentation/gitcredentials.txt-codeerasecode
func Erase(runner *op.Runner) error {
	input, err := gitcredential.Parse(runner.Stdin)
	if err != nil {
		return stdinError(err)
	}
	item, err := runner.FindItem(input.URL())
	if err != nil {
		if errors.Is(err, op.ErrNotFound) {
			return nil
		}
		return opError(err)
	}
	if err := runner.DeleteItem(item.ID); err != nil {
		return opError(err)
	}
	return nil
}
