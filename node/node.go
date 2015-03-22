package node

import (
	"time"
)

type Node struct {
	Id string `bson:"_id"`
	Type string `bson:"type"`
	Timestamp time.Time `bson:"timestamp"`
	Address string `bson:"address"`
}
