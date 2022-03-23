package mock

import (
	"github.com/google/uuid"
)

type ModelA struct {
	ID                    int    `model:"role:pk"`
	Name                  string `model:"role:str" randomcat:"random:hello"`
	InnerModel            *ModelB
	ChainInnerModel       []*ModelB
	CommonChainInnerModel []*ModelB
	RefObj                *RefObj
}

type ModelB struct {
	ID                    int `model:"role:pk"`
	Name                  string
	InnerModel            *ModelA
	ChainInnerModel       []*ModelA
	CommonChainInnerModel []*ModelB
	RefObj                *RefObj
}

type RefObj struct {
	ID int
}

type ModelC struct {
	ID                     int `model:"role:pk"`
	Name                   int
	ModelIncompatible      *ModelC
	ChainModelIncompatible []*ModelC
}

type ModelD struct {
	ID                     int `model:"role:pk"`
	Name                   int
	ModelIncompatible      *ModelB
	ChainModelIncompatible []*ModelB
}

type ModelE struct {
	ID                  int `model:"role:pk"`
	PointerIncompatible *map[string]string
}

type ModelF struct {
	ID                  int `model:"role:pk"`
	PointerIncompatible *map[int]int
}

type ModelG struct {
	ID uuid.UUID `model:"role:pk,"`
}

type ModelH struct {
	ID string `model:"role:pk,"`
}
