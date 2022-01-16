package roach

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"reflect"
	"time"
)

type JSONB map[string]interface{}

func allModelTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(Node{}),
		reflect.TypeOf(Range{}),
		reflect.TypeOf(RangeReplicaToNode{}),
		reflect.TypeOf(ChannelConfig{}),
		reflect.TypeOf(ChannelChunk{}),
		reflect.TypeOf(GossipNode{}),
		reflect.TypeOf(GossipLiveness{}),
	}
}

func newRoachModelFromStorage(m interface{}) interface{} {
	for _, rm := range allModelTypes() {
		rmName := rm.Name()
		mName := reflect.TypeOf(m).Elem().Name()
		if rmName == mName {
			return reflect.New(rm).Interface()
		}
	}
	return fmt.Errorf("roach baseModel could not be found")
}

type Node struct {
	bun.BaseModel `bun:"select:nodes_w_gossip,table:nodes"`
	ID            int
	GossipNode
	GossipLiveness
}

type Range struct {
	ID                int `bun:",pk"`
	LeaseHolderNodeID int
	LeaseHolderNode   *Node   `bun:"rel:belongs-to,join:lease_holder_node_id=id"`
	ReplicaNodes      []*Node `bun:"m2m:range_replica_to_nodes,join:Range=Node"`
}

type RangeReplicaToNode struct {
	ID      uuid.UUID `bun:"type:uuid,default:gen_random_uuid()"`
	RangeID int
	Range   *Range `bun:"rel:belongs-to,join:range_id=id"`
	NodeID  int
	Node    *Node `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelConfig struct {
	ID     int `bun:",pk"`
	Name   string
	NodeID int
	Node   *Node `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelChunk struct {
	ID              int `bun:",pk"`
	RangeID         int
	Range           *Range `bun:"rel:belongs-to,join:range_id=id"`
	ChannelConfigID int
	ChannelConfig   *Range `bun:"rel:belongs-to,join:channel_config_id=id"`
}

// || ROACH INTERNAL MODELS ||

// GossipNode lives in crdb's internal schema and tracks the nodes in the roach cluster
type GossipNode struct {
	Address   string    `bun:"type:text,scanonly"`
	StartedAt time.Time `bun:"type:timestamp,scanonly"`
	IsLive    bool      `bun:"type:boolean,scanonly"`
}

// GossipLiveness lives in crdb's internal schema and tracks the health of nodes in
//the roach cluster
type GossipLiveness struct {
	Epoch           int       `bun:"type:bigint,scanonly"`
	Expiration      string    `bun:"type:text,scanonly"`
	Draining        bool      `bun:"type:boolean,scanonly"`
	Decommissioning bool      `bun:"type:boolean,scanonly"`
	Membership      string    `bun:"type:text,scanonly"`
	UpdatedAt       time.Time `bun:"type:timestamp,scanonly"`
}
