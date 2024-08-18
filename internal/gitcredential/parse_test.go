package gitcredential

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	r := strings.NewReader("protocol=https\nhost=example.com\npath=foo.git\nusername=foo\npassword=bar")
	c, err := Parse(r)
	require.NoError(t, err)
	expect := Credential{
		Protocol: "https",
		Host:     "example.com",
		Path:     "foo.git",
		Username: "foo",
		Password: "bar",
	}
	require.Equal(t, expect, c)
}

func TestParse_RawURL(t *testing.T) {
	t.Run("only url; username", func(t *testing.T) {
		r := strings.NewReader("url=https://user@example.com/foo.git")
		c, err := Parse(r)
		require.NoError(t, err)
		expect := Credential{
			Protocol: "https",
			Host:     "example.com",
			Path:     "foo.git",
			Username: "user",
			RawURL:   "https://user@example.com/foo.git",
		}
		require.Equal(t, expect, c)
	})

	t.Run("url override; username", func(t *testing.T) {
		r := strings.NewReader("url=https://baz@example.com/foo.git\nusername=foo")
		c, err := Parse(r)
		require.NoError(t, err)
		expect := Credential{
			Protocol: "https",
			Host:     "example.com",
			Path:     "foo.git",
			Username: "baz",
			RawURL:   "https://baz@example.com/foo.git",
		}
		require.Equal(t, expect, c)
	})
}
