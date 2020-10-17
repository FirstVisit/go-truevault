// Code generated by mockery v2.1.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gotruevault "github.com/FirstVisit/go-truevault"
	document "github.com/FirstVisit/go-truevault/document"

	mock "github.com/stretchr/testify/mock"
)

// Document is an autogenerated mock type for the Document type
type Document struct {
	mock.Mock
}

// SearchDocument provides a mock function with given fields: ctx, vaultID, filter
func (_m *Document) SearchDocument(ctx context.Context, vaultID string, filter gotruevault.SearchFilter) (document.SearchDocumentResult, error) {
	ret := _m.Called(ctx, vaultID, filter)

	var r0 document.SearchDocumentResult
	if rf, ok := ret.Get(0).(func(context.Context, string, gotruevault.SearchFilter) document.SearchDocumentResult); ok {
		r0 = rf(ctx, vaultID, filter)
	} else {
		r0 = ret.Get(0).(document.SearchDocumentResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, gotruevault.SearchFilter) error); ok {
		r1 = rf(ctx, vaultID, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
