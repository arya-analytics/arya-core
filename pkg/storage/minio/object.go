package minio

import (
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

type object struct {
	*minio.Object
}

func (o *object) Size() int64 {
	stat, err := o.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	return stat.Size
}
