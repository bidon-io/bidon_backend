// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package config

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"sync"
)

// Ensure, that AppDemandProfileFetcherMock does implement AppDemandProfileFetcher.
// If this is not the case, regenerate this file with moq.
var _ AppDemandProfileFetcher = &AppDemandProfileFetcherMock{}

// AppDemandProfileFetcherMock is a mock implementation of AppDemandProfileFetcher.
//
//	func TestSomethingThatUsesAppDemandProfileFetcher(t *testing.T) {
//
//		// make and configure a mocked AppDemandProfileFetcher
//		mockedAppDemandProfileFetcher := &AppDemandProfileFetcherMock{
//			FetchFunc: func(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]AppDemandProfile, error) {
//				panic("mock out the Fetch method")
//			},
//		}
//
//		// use mockedAppDemandProfileFetcher in code that requires AppDemandProfileFetcher
//		// and then make assertions.
//
//	}
type AppDemandProfileFetcherMock struct {
	// FetchFunc mocks the Fetch method.
	FetchFunc func(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]AppDemandProfile, error)

	// calls tracks calls to the methods.
	calls struct {
		// Fetch holds details about calls to the Fetch method.
		Fetch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// AppID is the appID argument value.
			AppID int64
			// AdapterKeys is the adapterKeys argument value.
			AdapterKeys []adapter.Key
		}
	}
	lockFetch sync.RWMutex
}

// Fetch calls FetchFunc.
func (mock *AppDemandProfileFetcherMock) Fetch(ctx context.Context, appID int64, adapterKeys []adapter.Key) ([]AppDemandProfile, error) {
	if mock.FetchFunc == nil {
		panic("AppDemandProfileFetcherMock.FetchFunc: method is nil but AppDemandProfileFetcher.Fetch was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		AppID       int64
		AdapterKeys []adapter.Key
	}{
		Ctx:         ctx,
		AppID:       appID,
		AdapterKeys: adapterKeys,
	}
	mock.lockFetch.Lock()
	mock.calls.Fetch = append(mock.calls.Fetch, callInfo)
	mock.lockFetch.Unlock()
	return mock.FetchFunc(ctx, appID, adapterKeys)
}

// FetchCalls gets all the calls that were made to Fetch.
// Check the length with:
//
//	len(mockedAppDemandProfileFetcher.FetchCalls())
func (mock *AppDemandProfileFetcherMock) FetchCalls() []struct {
	Ctx         context.Context
	AppID       int64
	AdapterKeys []adapter.Key
} {
	var calls []struct {
		Ctx         context.Context
		AppID       int64
		AdapterKeys []adapter.Key
	}
	mock.lockFetch.RLock()
	calls = mock.calls.Fetch
	mock.lockFetch.RUnlock()
	return calls
}