package subcmd

import (
	"fmt"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
)

// Get is the "get" subcommand of the git-credential helper contract.
//
// See https://git-scm.com/docs/gitcredentials#Documentation/gitcredentials.txt-codegetcode
func Get(runner *op.Runner) error {
	input, err := gitcredential.Parse(runner.Stdin)
	if err != nil {
		return stdinError(err)
	}
	item, err := runner.FindItem(input.URL())
	if err != nil {
		return opError(err)
	}
	_, werr := fmt.Fprint(runner.Stdout, gitcredential.Credential{
		Username: item.Username,
		Password: item.Password,
	})
	return werr
}
