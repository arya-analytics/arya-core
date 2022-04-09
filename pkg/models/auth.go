package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	model.Base `storage:"engines:md,"`
	ID         uuid.UUID `model:"role:pk," bun:"type:uuid,pk"`
	Username   string
	Password   string   `bun:"type:varchar(60)"`
	Groups     []*Group `bun:"m2m:user_to_group,join:User=Group"`
}

type UserToGroup struct {
	bun.BaseModel `bun:"table:user_to_group"`
	model.Base    `storage:"engines:md," `
	ID            uuid.UUID `model:"role:pk," bun:"type:uuid,pk"`
	UserID        uuid.UUID `bun:"type:uuid"`
	User          *User     `model:"rel:belongs-to,join:UserID=ID" bun:"rel:belongs-to,join:user_id=id"`
	GroupID       uuid.UUID `bun:"type:uuid"`
	Group         *Group    `model:"rel:belongs-to,join:GroupID=ID" bun:"rel:belongs-to,join:group_id=id"`
}

type Group struct {
	model.Base `storage:"engines:md,"`
	ID         uuid.UUID `model:"role:pk," bun:"type:uuid,pk"`
	Name       string
}
