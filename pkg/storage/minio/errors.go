package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

func newErrorHandler() storage.ErrorHandler {
	return storage.NewErrorHandler(errorTypeConverterDefault)
}

func errorTypeConverterDefault(err error) (storage.ErrorType, bool) {
	mErr := minio.ToErrorResponse(err)
	errT, ok := _minioErrors[mErr.Code]
	return errT, ok
}

var _minioErrors = map[string]storage.ErrorType{
	"NoSuchKey": storage.ErrorTypeItemNotFound,
}
