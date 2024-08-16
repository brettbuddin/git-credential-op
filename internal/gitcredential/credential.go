package gitcredential

import (
	"fmt"
	"strings"
)

// Credential is a subset of the attributes in the git-credential input/output format.
//
// See https://git-scm.com/docs/git-credential#IOFMT for details.
type Credential struct {
	// Protocol is the protocol in which the credential will be used (e.g. "https").
	Protocol string
	// Host is the remote hostname for the credential. Includes the port if one was specified.
	Host string
	// Path is the remote server path for the credential.
	Path string
	// Username is the username for the credential if we already have one.
	Username string
	// Password is the password for the credential if we're asking for it to be stored.
	Password string
	// RawURL is the full URL for the credential. Alternative to passing in the parts.
	RawURL string
}

// URL returns as a full URL for the credential's remote upstream.
func (c Credential) URL() string {
	// Don't include the trailing slash if we don't have a path.
	path := c.Path
	if path != "" {
		path = "/" + path
	}
	return fmt.Sprintf("%s://%s%s", c.Protocol, c.Host, path)
}

func (c Credential) String() string {
	var s strings.Builder
	if c.Protocol != "" {
		s.WriteString("protocol=" + c.Protocol + "\n")
	}
	if c.Host != "" {
		s.WriteString("host=" + c.Host + "\n")
	}
	if c.Path != "" {
		s.WriteString("path=" + c.Path + "\n")
	}
	if c.Username != "" {
		s.WriteString("username=" + c.Username + "\n")
	}
	if c.Password != "" {
		s.WriteString("password=" + c.Password + "\n")
	}
	return s.String()
}
