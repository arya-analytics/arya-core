package roach

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"reflect"
	"time"
)

type JSONB map[string]interface{}

func models() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(ChannelConfig{}),
	}
}

func roachModelFromStorage(m interface{}) interface{} {
	for _, rm := range models() {
		rmName := rm.Name()
		mName := reflect.TypeOf(m).Elem().Name()
		if rmName == mName {
			return reflect.New(rm).Interface()
		}
	}
	return fmt.Errorf("roach model could not be found")
}

type Node struct {
	ID uuid.UUID `bun:"type:uuid,default:gen_random_uuid(),pk"`
	GossipNodeID int64
	GossipNode *GossipNode `bun:"rel:belongs-to,join:gossip_node_id=id"`
}

type Range struct {
	ID int32 `bun:",pk"`
	string
	LeaseHolderNodeId uuid.UUID
	LeaseHolderNode   *Node  `bun:"rel:belongs-to,join:lease_holder_node_id=id"`
	ReplicaNodes      []Node `bun:"m2m:range_replica_to_nodes,join:Range=Node"`
}

type RangeReplicaToNode struct {
	ID      uuid.UUID `bun:"type:uuid,default:gen_random_uuid()"`
	RangeID int32
	Range   *Range    `bun:"rel:belongs-to,join:range_id=id"`
	NodeID  uuid.UUID `bun:"type:uuid"`
	Node    *Node     `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelConfig struct {
	ID     int32 `bun:",pk"`
	Name   string
	//NodeId uuid.UUID `bun:"type:uuid"`
	//Node   *Node     `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelChunk struct {
	ID              int64 `bun:",pk"`
	RangeID         int32
	Range           *Range `bun:"rel:belongs-to,join:range_id=id"`
	ChannelConfigID int32
	ChannelConfig   *Range `bun:"rel:belongs-to,join:channel_config_id=id"`
}

// || ROACH INTERNAL MODELS ||

// GossipNode lives in crdb's internal schema and tracks the nodes in the roach cluster
type GossipNode struct {
	bun.BaseModel       `bun:"table:crdb_internal.gossip_nodes"`
	NodeID              int       `bun:"type:bigint,pk"`
	Network             string    `bun:"type:text"`
	Address             string    `bun:"type:text"`
	AdvertiseAddress    string    `bun:"type:text"`
	SQLNetwork          string    `bun:"type:text"`
	AdvertiseSQLAddress string    `bun:"type:type"`
	Attrs               JSONB     `bun:"type:jsonb"`
	Locality            string    `bun:"type:text"`
	ClusterName         string    `bun:"type:text"`
	ServerVersion       string    `bun:"type:text"`
	BuildTag            string    `bun:"type:text"`
	StartedAt           time.Time `bun:"type:timestamp"`
	IsLive              bool      `bun:"type:boolean"`
	Ranges              int       `bun:"type:bigint"`
	Leases              int       `bun:"type:bigint"`
}

// GossipLiveness lives in crdb's internal schema and tracks the health of nodes in
//the roach cluster
type GossipLiveness struct {
	bun.BaseModel   `bun:"table:crdb_internal.gossip_liveness"`
	NodeID          int       `bun:"type:bigint"`
	Epoch           int       `bun:"type:bigint"`
	Expiration      string    `bun:"type:text"`
	Draining        bool      `bun:"type:boolean"`
	Decommissioning bool      `bun:"type:boolean"`
	Membership      string    `bun:"type:text"`
	UpdatedAt       time.Time `bun:"type:timestamp"`
}
