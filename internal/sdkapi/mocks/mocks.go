// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/sdkapi"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"sync"
)

// Ensure, that AppFetcherMock does implement sdkapi.AppFetcher.
// If this is not the case, regenerate this file with moq.
var _ sdkapi.AppFetcher = &AppFetcherMock{}

// AppFetcherMock is a mock implementation of sdkapi.AppFetcher.
//
//	func TestSomethingThatUsesAppFetcher(t *testing.T) {
//
//		// make and configure a mocked sdkapi.AppFetcher
//		mockedAppFetcher := &AppFetcherMock{
//			FetchFunc: func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
//				panic("mock out the Fetch method")
//			},
//		}
//
//		// use mockedAppFetcher in code that requires sdkapi.AppFetcher
//		// and then make assertions.
//
//	}
type AppFetcherMock struct {
	// FetchFunc mocks the Fetch method.
	FetchFunc func(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error)

	// calls tracks calls to the methods.
	calls struct {
		// Fetch holds details about calls to the Fetch method.
		Fetch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// AppKey is the appKey argument value.
			AppKey string
			// AppBundle is the appBundle argument value.
			AppBundle string
		}
	}
	lockFetch sync.RWMutex
}

// Fetch calls FetchFunc.
func (mock *AppFetcherMock) Fetch(ctx context.Context, appKey string, appBundle string) (sdkapi.App, error) {
	if mock.FetchFunc == nil {
		panic("AppFetcherMock.FetchFunc: method is nil but AppFetcher.Fetch was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		AppKey    string
		AppBundle string
	}{
		Ctx:       ctx,
		AppKey:    appKey,
		AppBundle: appBundle,
	}
	mock.lockFetch.Lock()
	mock.calls.Fetch = append(mock.calls.Fetch, callInfo)
	mock.lockFetch.Unlock()
	return mock.FetchFunc(ctx, appKey, appBundle)
}

// FetchCalls gets all the calls that were made to Fetch.
// Check the length with:
//
//	len(mockedAppFetcher.FetchCalls())
func (mock *AppFetcherMock) FetchCalls() []struct {
	Ctx       context.Context
	AppKey    string
	AppBundle string
} {
	var calls []struct {
		Ctx       context.Context
		AppKey    string
		AppBundle string
	}
	mock.lockFetch.RLock()
	calls = mock.calls.Fetch
	mock.lockFetch.RUnlock()
	return calls
}

// Ensure, that GeocoderMock does implement sdkapi.Geocoder.
// If this is not the case, regenerate this file with moq.
var _ sdkapi.Geocoder = &GeocoderMock{}

// GeocoderMock is a mock implementation of sdkapi.Geocoder.
//
//	func TestSomethingThatUsesGeocoder(t *testing.T) {
//
//		// make and configure a mocked sdkapi.Geocoder
//		mockedGeocoder := &GeocoderMock{
//			LookupFunc: func(ctx context.Context, ipString string) (geocoder.GeoData, error) {
//				panic("mock out the Lookup method")
//			},
//		}
//
//		// use mockedGeocoder in code that requires sdkapi.Geocoder
//		// and then make assertions.
//
//	}
type GeocoderMock struct {
	// LookupFunc mocks the Lookup method.
	LookupFunc func(ctx context.Context, ipString string) (geocoder.GeoData, error)

	// calls tracks calls to the methods.
	calls struct {
		// Lookup holds details about calls to the Lookup method.
		Lookup []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IpString is the ipString argument value.
			IpString string
		}
	}
	lockLookup sync.RWMutex
}

// Lookup calls LookupFunc.
func (mock *GeocoderMock) Lookup(ctx context.Context, ipString string) (geocoder.GeoData, error) {
	if mock.LookupFunc == nil {
		panic("GeocoderMock.LookupFunc: method is nil but Geocoder.Lookup was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		IpString string
	}{
		Ctx:      ctx,
		IpString: ipString,
	}
	mock.lockLookup.Lock()
	mock.calls.Lookup = append(mock.calls.Lookup, callInfo)
	mock.lockLookup.Unlock()
	return mock.LookupFunc(ctx, ipString)
}

// LookupCalls gets all the calls that were made to Lookup.
// Check the length with:
//
//	len(mockedGeocoder.LookupCalls())
func (mock *GeocoderMock) LookupCalls() []struct {
	Ctx      context.Context
	IpString string
} {
	var calls []struct {
		Ctx      context.Context
		IpString string
	}
	mock.lockLookup.RLock()
	calls = mock.calls.Lookup
	mock.lockLookup.RUnlock()
	return calls
}