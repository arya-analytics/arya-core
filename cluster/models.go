package cluster

import "github.com/google/uuid"

type Node struct {
	ID uuid.UUID `bun:"type:uuid,default:gen_random_uuid()"`
	Addr string `bun:"type:inet,unique"`
}
