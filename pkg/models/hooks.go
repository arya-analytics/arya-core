package models

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryHookRunner struct {
	rfl *model.Reflect
	*errutil.CatchSimple
}

func (qhr *queryHookRunner) Exec(action func(rfl *model.Reflect) error) {
	qhr.CatchSimple.Exec(func() error { return action(qhr.rfl) })
}

func BindHooks(s storage.Storage) {
	hooks := []storage.QueryHook{
		&NodeQueryHook{},
	}
	for _, hook := range hooks {
		s.AddQueryHook(hook)
	}
}
