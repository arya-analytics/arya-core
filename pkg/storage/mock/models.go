package mock

type ModelA struct {
	ID                    int
	Name                  string
	InnerModel            *ModelB
	ChainInnerModel       []*ModelB
	CommonChainInnerModel []*ModelB
	RefObj                *RefObj
}

type ModelB struct {
	ID                    int
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
	ID         int
	Name       int
	InnerModel *ModelC
}

type ModelD struct {
	ID         int
	Name       int
	InnerModel *ModelB
}
