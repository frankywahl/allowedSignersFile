// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package github_test

import (
	"github.com/frankywahl/allowedSignatures/internal/github"
	"sync"
)

// Ensure, that LoggerMock does implement github.Logger.
// If this is not the case, regenerate this file with moq.
var _ github.Logger = &LoggerMock{}

// LoggerMock is a mock implementation of github.Logger.
//
//	func TestSomethingThatUsesLogger(t *testing.T) {
//
//		// make and configure a mocked github.Logger
//		mockedLogger := &LoggerMock{
//			InfofFunc: func(format string, opts ...interface{})  {
//				panic("mock out the Infof method")
//			},
//		}
//
//		// use mockedLogger in code that requires github.Logger
//		// and then make assertions.
//
//	}
type LoggerMock struct {
	// InfofFunc mocks the Infof method.
	InfofFunc func(format string, opts ...interface{})

	// calls tracks calls to the methods.
	calls struct {
		// Infof holds details about calls to the Infof method.
		Infof []struct {
			// Format is the format argument value.
			Format string
			// Opts is the opts argument value.
			Opts []interface{}
		}
	}
	lockInfof sync.RWMutex
}

// Infof calls InfofFunc.
func (mock *LoggerMock) Infof(format string, opts ...interface{}) {
	if mock.InfofFunc == nil {
		panic("LoggerMock.InfofFunc: method is nil but Logger.Infof was just called")
	}
	callInfo := struct {
		Format string
		Opts   []interface{}
	}{
		Format: format,
		Opts:   opts,
	}
	mock.lockInfof.Lock()
	mock.calls.Infof = append(mock.calls.Infof, callInfo)
	mock.lockInfof.Unlock()
	mock.InfofFunc(format, opts...)
}

// InfofCalls gets all the calls that were made to Infof.
// Check the length with:
//
//	len(mockedLogger.InfofCalls())
func (mock *LoggerMock) InfofCalls() []struct {
	Format string
	Opts   []interface{}
} {
	var calls []struct {
		Format string
		Opts   []interface{}
	}
	mock.lockInfof.RLock()
	calls = mock.calls.Infof
	mock.lockInfof.RUnlock()
	return calls
}
