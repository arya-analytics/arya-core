package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func parseMinioErr(err error) error {
	if err == nil {
		return nil
	}
	mErr := minio.ToErrorResponse(err)
	log.Info(err)
	return storage.Error{Base: err, Type: _minioErrors[mErr.Code],
		Message: mErr.Error()}
}

var _minioErrors = map[string]storage.ErrorType{
	"NoSuchKey": storage.ErrTypeItemNotFound,
}
