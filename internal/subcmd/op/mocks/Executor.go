// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	op "github.com/brettbuddin/git-credential-op/internal/subcmd/op"
	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

type Executor_Expecter struct {
	mock *mock.Mock
}

func (_m *Executor) EXPECT() *Executor_Expecter {
	return &Executor_Expecter{mock: &_m.Mock}
}

// ExecuteCommand provides a mock function with given fields: out, name, args
func (_m *Executor) ExecuteCommand(out op.ExecutorOutput, name string, args ...string) error {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, out, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteCommand")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(op.ExecutorOutput, string, ...string) error); ok {
		r0 = rf(out, name, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Executor_ExecuteCommand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteCommand'
type Executor_ExecuteCommand_Call struct {
	*mock.Call
}

// ExecuteCommand is a helper method to define mock.On call
//   - out op.ExecutorOutput
//   - name string
//   - args ...string
func (_e *Executor_Expecter) ExecuteCommand(out interface{}, name interface{}, args ...interface{}) *Executor_ExecuteCommand_Call {
	return &Executor_ExecuteCommand_Call{Call: _e.mock.On("ExecuteCommand",
		append([]interface{}{out, name}, args...)...)}
}

func (_c *Executor_ExecuteCommand_Call) Run(run func(out op.ExecutorOutput, name string, args ...string)) *Executor_ExecuteCommand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(op.ExecutorOutput), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *Executor_ExecuteCommand_Call) Return(_a0 error) *Executor_ExecuteCommand_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Executor_ExecuteCommand_Call) RunAndReturn(run func(op.ExecutorOutput, string, ...string) error) *Executor_ExecuteCommand_Call {
	_c.Call.Return(run)
	return _c
}

// NewExecutor creates a new instance of Executor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *Executor {
	mock := &Executor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
