package subcmd_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/mockutil"
	"github.com/brettbuddin/git-credential-op/internal/subcmd"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op/opmock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	executor := opmock.NewExecutor(t)
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

	mockutil.InOrder(listCall, getCall)

	err := subcmd.Get(runner)
	require.NoError(t, err)

	b, err := io.ReadAll(stdout)
	require.NoError(t, err)
	require.Equal(t, "username=brettbuddin\npassword=gizmo\n", string(b))
}
