package op

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpItemCommand(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		args := new(Runner).itemCmd("list", "--categories", "login")
		require.Equal(t, []string{
			"item",
			"list",
			"--categories",
			"login",
		}, args)
	})

	t.Run("with account", func(t *testing.T) {
		r := &Runner{Account: "companyaccount.1password.com"}
		args := r.itemCmd("list", "--categories", "login")
		require.Equal(t, []string{
			"item",
			"list",
			"--account", "companyaccount.1password.com",
			"--categories",
			"login",
		}, args)
	})

	t.Run("with vault", func(t *testing.T) {
		r := &Runner{Vault: "Team"}
		args := r.itemCmd("list", "--categories", "login")
		require.Equal(t, []string{
			"item",
			"list",
			"--vault", "Team",
			"--categories",
			"login",
		}, args)
	})
}
