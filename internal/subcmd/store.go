package subcmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
)

// DefaultTitleFormat is the default format used for item titles in 1Password.
const DefaultTitleFormat = "{{.Host}} ({{.LocatorTag}})"

// Store is the "store" subcommand of the git-credential helper contract.
//
// See https://git-scm.com/docs/gitcredentials#Documentation/gitcredentials.txt-codestorecode
func Store(runner *op.Runner, titleFormat string, additionalTags []string) error {
	input, err := gitcredential.Parse(runner.Stdin)
	if err != nil {
		return stdinError(err)
	}

	item, err := runner.FindItem(input.URL())
	if err == nil {
		if item.Username == input.Username && item.Password == input.Password {
			return nil
		}
		item.Username = input.Username
		item.Password = input.Password
		if err := runner.UpdateItem(item); err != nil {
			return opError(err)
		}
		return nil
	}

	if !errors.Is(err, op.ErrNotFound) {
		return opError(err)
	}

	title, err := renderTitle(input, runner.LocatorTag, titleFormat)
	if err != nil {
		return fmt.Errorf("render title: %w", err)
	}
	cErr := runner.CreateItem(op.CreateRequest{
		Title:          title,
		AdditionalTags: additionalTags,
		Username:       input.Username,
		Password:       input.Password,
		URL:            input.URL(),
	})
	if cErr != nil {
		return opError(cErr)
	}
	return nil
}

func renderTitle(
	c gitcredential.Credential,
	locatorTag string,
	titleFormat string,
) (string, error) {
	titleTemplate, err := template.New("").Parse(titleFormat)
	if err != nil {
		return "", fmt.Errorf("parse title format: %w", err)
	}

	type titleTemplateData struct {
		Protocol   string
		Host       string
		LocatorTag string
	}

	var title bytes.Buffer
	err = titleTemplate.Execute(&title, titleTemplateData{
		Protocol:   c.Protocol,
		Host:       c.Host,
		LocatorTag: locatorTag,
	})
	if err != nil {
		return "", err
	}
	return title.String(), nil
}
