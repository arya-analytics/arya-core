package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/minio/minio-go/v7"
)

func newErrorConvert() errutil.ConvertChain {
	return query.NewErrorConvertChain(errorConvertDefault)
}

func errorConvertDefault(err error) (error, bool) {
	mErr := minio.ToErrorResponse(err)
	t, ok := minioErrors()[mErr.Code]
	return query.NewSimpleError(t, err), ok
}

func minioErrors() map[string]query.ErrorType {
	return map[string]query.ErrorType{
		"NoSuchKey": query.ErrorTypeItemNotFound,
	}
}
