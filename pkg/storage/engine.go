package storage

type Engine interface {
	Type() EngineType
}

type EngineType int

const (
	MetaData EngineType = iota
	Bulk
	Cache
)

type MetaDataEngine interface {
	Engine
}