package subcmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/brettbuddin/git-credential-op/internal/gitcredential"
	"github.com/brettbuddin/git-credential-op/internal/mockutil"
	"github.com/brettbuddin/git-credential-op/internal/subcmd"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op"
	"github.com/brettbuddin/git-credential-op/internal/subcmd/op/opmock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStore_UpdateItem_PasswordChanged(t *testing.T) {
	executor := opmock.NewExecutor(t)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	input := gitcredential.Credential{
		Protocol: "http",
		Host:     "localhost:10000",
		Username: "brettbuddin",
		Password: "newpassword",
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
		fmt.Fprintln(out.Stdout, `[{"label":"username", "value":"brettbuddin"}, {"label":"password", "value":"oldpassword"}]`)
		return nil
	}).Once()

	editCall := exp.ExecuteCommand(
		mock.Anything,
		"op",
		"item",
		"edit",
		"id1",
		"username=brettbuddin",
		"password=newpassword",
	).Return(nil).Once()

	mockutil.InOrder(listCall, getCall, editCall)

	err := subcmd.Store(runner, subcmd.DefaultTitleFormat, nil)
	require.NoError(t, err)
}

func TestStore_UpdateItem_PasswordUnchanged(t *testing.T) {
	executor := opmock.NewExecutor(t)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	input := gitcredential.Credential{
		Protocol: "http",
		Host:     "localhost:10000",
		Username: "brettbuddin",
		Password: "gizmo",
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

	// No update will take place, because the password hasn't changed.

	err := subcmd.Store(runner, subcmd.DefaultTitleFormat, nil)
	require.NoError(t, err)
}

func TestStore_CreateItem(t *testing.T) {
	executor := opmock.NewExecutor(t)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	input := gitcredential.Credential{
		Protocol: "http",
		Host:     "localhost:10000",
		Username: "brettbuddin",
		Password: "gizmo",
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
		fmt.Fprintln(out.Stdout, `[]`)
		return nil
	}).Once()

	createCall := exp.ExecuteCommand(
		mock.Anything,
		"op",
		"item",
		"create",
		"--category", "login",
		"--title", "localhost:10000 (git-credential-op)",
		"--url", "http://localhost:10000",
		"--tags", op.DefaultLocatorTag,
		"username=brettbuddin",
		"password=gizmo",
	).Return(nil).Once()

	mockutil.InOrder(listCall, createCall)

	err := subcmd.Store(runner, subcmd.DefaultTitleFormat, nil)
	require.NoError(t, err)
}
