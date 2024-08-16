package subcmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/mockutil"
	"github.com/brettbuddin/git-credential-op/internal/subcmd"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestErase(t *testing.T) {
	t.Run("successfully delete item", func(t *testing.T) {
		executor := mocks.NewExecutor(t)
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		input := gitcredential.Credential{
			Protocol: "http",
			Host:     "localhost:10000",
		}
		runner := &op.Runner{
			LocatorTag: op.DefaultLocatorTag,
			Executor:   executor,
			Stdin:      bytes.NewBufferString(input.String()),
			Stdout:     stdout,
			Stderr:     stderr,
		}

		exp := executor.EXPECT()
		listCall := exp.ExecuteCommand(
			mock.Anything,
			"op",
			"item",
			"list",
			"--categories", "login",
			"--format", "json",
			"--tags", op.DefaultLocatorTag,
		).RunAndReturn(func(out op.ExecutorOutput, name string, args ...string) error {
			fmt.Fprintln(out.Stdout, `[{"id":"id1", "urls":[{"href":"http://localhost:10000"}]}]`)
			return nil
		}).Once()

		getCall := exp.ExecuteCommand(
			mock.Anything,
			"op",
			"item",
			"get",
			"--format", "json",
			"--fields", "username,password",
			"id1",
		).RunAndReturn(func(out op.ExecutorOutput, name string, args ...string) error {
			fmt.Fprintln(out.Stdout, `[{"label":"username", "value":"brettbuddin"}, {"label":"password", "value":"gizmo"}]`)
			return nil
		}).Once()

		deleteCall := exp.ExecuteCommand(
			mock.Anything,
			"op",
			"item",
			"delete",
			"id1",
		).
			Return(nil).
			Once()

		mockutil.InOrder(listCall, getCall, deleteCall)

		err := subcmd.Erase(runner)
		require.NoError(t, err)
	})

	t.Run("item not found", func(t *testing.T) {
		executor := mocks.NewExecutor(t)
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		input := gitcredential.Credential{
			Protocol: "http",
			Host:     "localhost:10000",
		}
		runner := &op.Runner{
			LocatorTag: op.DefaultLocatorTag,
			Executor:   executor,
			Stdin:      bytes.NewBufferString(input.String()),
			Stdout:     stdout,
			Stderr:     stderr,
		}

		executor.EXPECT().ExecuteCommand(
			mock.Anything,
			"op",
			"item",
			"list",
			"--categories", "login",
			"--format", "json",
			"--tags", op.DefaultLocatorTag,
		).RunAndReturn(func(out op.ExecutorOutput, name string, args ...string) error {
			fmt.Fprintln(out.Stdout, `[]`)
			return nil
		}).Once()

		err := subcmd.Erase(runner)
		require.NoError(t, err)
	})
}
