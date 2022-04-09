package internal

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"strings"
)

type EngineRole string

const (
	TagCat                      = "storage"
	EngineKey                   = "engines"
	EngineRoleMD     EngineRole = "md"
	EngineRoleObject EngineRole = "obj"
	EngineRoleCache  EngineRole = "cache"
)

func RequiresEngine(m interface{}, engine Engine) bool {
	rfl := model.NewReflect(m)
	bt, ok := rfl.StructTagChain().RetrieveBase()
	if !ok {
		panic("model must have a base struct tag")
	}
	v, ok := bt.Retrieve(TagCat, EngineKey)
	if !ok {
		panic("model must have an engine role specified")
	}
	switch engine.(type) {
	case EngineMD:
		return strings.Contains(v, string(EngineRoleMD))
	case EngineObject:
		return strings.Contains(v, string(EngineRoleObject))
	case EngineCache:
		return strings.Contains(v, string(EngineRoleCache))
	default:
		panic("invalid engine specified")
	}
}
