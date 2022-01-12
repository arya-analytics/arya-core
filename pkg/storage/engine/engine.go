package engine

import "github.com/google/uuid"

type Adapter interface {
	ID() uuid.UUID
	Role() Role
}

type Role int

const (
	RoleMetaData = iota
	RoleCache
	RoleBulk
)

type Base interface {
	NewAdapter() Adapter
	Role() Role
}
