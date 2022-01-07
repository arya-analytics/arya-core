package storage

type BaseEngine interface {
	Type() EngineType
}

type EngineType int

const (
	MetaData EngineType = iota
	Bulk
	Cache
)

type MetaDataEngine interface {
	BaseEngine
	NewRetrieve() MetaDataRetrieve
	NewUpdate() MetaDataCreate
}

type MetaDataRetrieve interface {
	Exec() error
	Where() MetaDataRetrieve
	Model(interface{}) MetaDataCreate
}

type MetaDataCreate interface {
	Exec() error
	Model(interface{}) MetaDataCreate
}