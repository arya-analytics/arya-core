package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
)

type baseQuery struct {
	storageWrapper *storage.ModelWrapper
	roachWrapper   *storage.ModelWrapper
}

func (b *baseQuery) bindWrappers(m interface{}) {
	b.storageWrapper = storage.NewModelWrapper(m)
	roachModel := roachModelFromStorage(m)
	b.roachWrapper = storage.NewModelWrapper(roachModel)
}
