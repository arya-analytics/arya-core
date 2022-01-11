package storage

type BaseEngine interface {
	Role() EngineRole
	Type() EngineType
}

type EngineType int

// TODO: Generate stringers for enums
const (
	EngineTypeRoach EngineType = iota
	EngineTypeMinio
	EngineTypeRedisTS
	EngineTypeMDStub
	EngineTypeBulkStub
	EngineTypeCacheStub
)

type EngineRole int

const (
	EngineRoleMetaData EngineRole = iota
	EngineRoleBulk
	EngineRoleCache
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
