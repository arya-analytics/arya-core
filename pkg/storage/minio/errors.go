package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var _errTypeConverterChain = storage.ErrorTypeConverterChain{}

var _defaultConverter = errConverterDefault

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(_errTypeConverterChain, _defaultConverter)
}

func errConverterDefault(err error) (storage.ErrorType, bool) {
	mErr := minio.ToErrorResponse(err)
	errT, ok := _minioErrors[mErr.Code]
	return errT, ok
}

var _minioErrors = map[string]storage.ErrorType{
	"NoSuchKey": storage.ErrTypeItemNotFound,
}
