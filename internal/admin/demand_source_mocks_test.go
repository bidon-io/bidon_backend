// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package admin

import (
	"context"
	"sync"
)

// Ensure, that DemandSourceRepoMock does implement DemandSourceRepo.
// If this is not the case, regenerate this file with moq.
var _ DemandSourceRepo = &DemandSourceRepoMock{}

// DemandSourceRepoMock is a mock implementation of DemandSourceRepo.
//
//	func TestSomethingThatUsesDemandSourceRepo(t *testing.T) {
//
//		// make and configure a mocked DemandSourceRepo
//		mockedDemandSourceRepo := &DemandSourceRepoMock{
//			CreateFunc: func(ctx context.Context, attrs *DemandSourceAttrs) (*DemandSource, error) {
//				panic("mock out the Create method")
//			},
//			DeleteFunc: func(ctx context.Context, id int64) error {
//				panic("mock out the Delete method")
//			},
//			FindFunc: func(ctx context.Context, id int64) (*DemandSource, error) {
//				panic("mock out the Find method")
//			},
//			ListFunc: func(contextMoqParam context.Context) ([]DemandSource, error) {
//				panic("mock out the List method")
//			},
//			UpdateFunc: func(ctx context.Context, id int64, attrs *DemandSourceAttrs) (*DemandSource, error) {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockedDemandSourceRepo in code that requires DemandSourceRepo
//		// and then make assertions.
//
//	}
type DemandSourceRepoMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, attrs *DemandSourceAttrs) (*DemandSource, error)

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, id int64) error

	// FindFunc mocks the Find method.
	FindFunc func(ctx context.Context, id int64) (*DemandSource, error)

	// ListFunc mocks the List method.
	ListFunc func(contextMoqParam context.Context) ([]DemandSource, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, id int64, attrs *DemandSourceAttrs) (*DemandSource, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Attrs is the attrs argument value.
			Attrs *DemandSourceAttrs
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
		}
		// Find holds details about calls to the Find method.
		Find []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
		}
		// List holds details about calls to the List method.
		List []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
			// Attrs is the attrs argument value.
			Attrs *DemandSourceAttrs
		}
	}
	lockCreate sync.RWMutex
	lockDelete sync.RWMutex
	lockFind   sync.RWMutex
	lockList   sync.RWMutex
	lockUpdate sync.RWMutex
}

// Create calls CreateFunc.
func (mock *DemandSourceRepoMock) Create(ctx context.Context, attrs *DemandSourceAttrs) (*DemandSource, error) {
	if mock.CreateFunc == nil {
		panic("DemandSourceRepoMock.CreateFunc: method is nil but DemandSourceRepo.Create was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Attrs *DemandSourceAttrs
	}{
		Ctx:   ctx,
		Attrs: attrs,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, attrs)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedDemandSourceRepo.CreateCalls())
func (mock *DemandSourceRepoMock) CreateCalls() []struct {
	Ctx   context.Context
	Attrs *DemandSourceAttrs
} {
	var calls []struct {
		Ctx   context.Context
		Attrs *DemandSourceAttrs
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *DemandSourceRepoMock) Delete(ctx context.Context, id int64) error {
	if mock.DeleteFunc == nil {
		panic("DemandSourceRepoMock.DeleteFunc: method is nil but DemandSourceRepo.Delete was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int64
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, id)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedDemandSourceRepo.DeleteCalls())
func (mock *DemandSourceRepoMock) DeleteCalls() []struct {
	Ctx context.Context
	ID  int64
} {
	var calls []struct {
		Ctx context.Context
		ID  int64
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// Find calls FindFunc.
func (mock *DemandSourceRepoMock) Find(ctx context.Context, id int64) (*DemandSource, error) {
	if mock.FindFunc == nil {
		panic("DemandSourceRepoMock.FindFunc: method is nil but DemandSourceRepo.Find was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int64
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockFind.Lock()
	mock.calls.Find = append(mock.calls.Find, callInfo)
	mock.lockFind.Unlock()
	return mock.FindFunc(ctx, id)
}

// FindCalls gets all the calls that were made to Find.
// Check the length with:
//
//	len(mockedDemandSourceRepo.FindCalls())
func (mock *DemandSourceRepoMock) FindCalls() []struct {
	Ctx context.Context
	ID  int64
} {
	var calls []struct {
		Ctx context.Context
		ID  int64
	}
	mock.lockFind.RLock()
	calls = mock.calls.Find
	mock.lockFind.RUnlock()
	return calls
}

// List calls ListFunc.
func (mock *DemandSourceRepoMock) List(contextMoqParam context.Context) ([]DemandSource, error) {
	if mock.ListFunc == nil {
		panic("DemandSourceRepoMock.ListFunc: method is nil but DemandSourceRepo.List was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(contextMoqParam)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockedDemandSourceRepo.ListCalls())
func (mock *DemandSourceRepoMock) ListCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *DemandSourceRepoMock) Update(ctx context.Context, id int64, attrs *DemandSourceAttrs) (*DemandSource, error) {
	if mock.UpdateFunc == nil {
		panic("DemandSourceRepoMock.UpdateFunc: method is nil but DemandSourceRepo.Update was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		ID    int64
		Attrs *DemandSourceAttrs
	}{
		Ctx:   ctx,
		ID:    id,
		Attrs: attrs,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(ctx, id, attrs)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedDemandSourceRepo.UpdateCalls())
func (mock *DemandSourceRepoMock) UpdateCalls() []struct {
	Ctx   context.Context
	ID    int64
	Attrs *DemandSourceAttrs
} {
	var calls []struct {
		Ctx   context.Context
		ID    int64
		Attrs *DemandSourceAttrs
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}