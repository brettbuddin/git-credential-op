package mockutil

import "github.com/stretchr/testify/mock"

// InOrder ensures each of the calls occurs in the specified order.
func InOrder(calls ...*mock.Call) {
	for i := len(calls) - 1; i > 0; i-- {
		calls[i].NotBefore(calls[0:i]...)
	}
}
