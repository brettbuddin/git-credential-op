package ioutil

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixWriter(t *testing.T) {
	var buf bytes.Buffer
	p := PrefixWriter(&buf, "example: ")

	fmt.Fprintln(p, "hello")
	fmt.Fprintln(p, "world")

	require.Equal(t, "example: hello\nexample: world\n", buf.String())
}
