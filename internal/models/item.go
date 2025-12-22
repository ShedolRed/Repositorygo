package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	CreatedAt int64                  `bson:"created_at_unix" json:"created_at_unix"`
	UpdatedAt int64                  `bson:"updated_at_unix" json:"updated_at_unix"`
}
