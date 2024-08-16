package gitcredential

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// ErrNoKVPairFound is returned if no key-value pair is found on a line of input.
var ErrNoKVPairFound = fmt.Errorf("no key-value pair found")

// Parse scans attribute key-value pairs from an [io.Reader] and returns a Credential that contains
// them.
//
// The git-credential system ignores all attributes it doesn't support. We do the same here.
//
// An error is returned if any of the input lines do not contain a key-value pair.
func Parse(r io.Reader) (Credential, error) {
	var c Credential
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		k, v, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			return c, ErrNoKVPairFound
		}
		switch k {
		case "protocol":
			c.Protocol = v
		case "host":
			c.Host = v
		case "path":
			c.Path = v
		case "username":
			c.Username = v
		case "password":
			c.Password = v
		case "url":
			c.RawURL = v
		}
	}

	if c.RawURL != "" {
		parsed, err := url.Parse(c.RawURL)
		if err != nil {
			return c, fmt.Errorf("invalid url: %w", err)
		}
		c.Protocol = parsed.Scheme
		c.Host = parsed.Host
		c.Path = strings.TrimLeft(parsed.Path, "/")

		if parsed.User != nil {
			c.Username = parsed.User.Username()
		}
	}

	return c, nil
}
