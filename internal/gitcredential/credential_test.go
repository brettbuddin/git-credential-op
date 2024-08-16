package gitcredential

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialURL(t *testing.T) {
	c := Credential{
		Protocol: "https",
		Host:     "example.com",
		Path:     "foo.git",
	}
	require.Equal(t, "https://example.com/foo.git", c.URL())
}

func TestCredentialStringFmt(t *testing.T) {
	c := Credential{
		Protocol: "https",
		Host:     "example.com",
		Path:     "foo.git",
		Username: "brettbuddin",
		Password: "gizmo",
	}
	require.Equal(t, "protocol=https\nhost=example.com\npath=foo.git\nusername=brettbuddin\npassword=gizmo\n", c.String())
}
