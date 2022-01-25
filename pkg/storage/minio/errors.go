package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

func parseMinioErr(err error) error {
	mErr := minio.ToErrorResponse(err)
	return storage.Error{Base: err, Type: _minioErrors[mErr.Code],
		Message: mErr.Error()}
}

var _minioErrors = map[string]storage.ErrorType{
	"NoSuchKey": storage.ErrTypeItemNotFound,
}
