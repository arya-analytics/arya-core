package minio

import (
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

type Object struct {
	*minio.Object
}

func (o *Object) Size() int64 {
	stat, err := o.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	return stat.Size
}
