package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type VerticeInfoSchema struct {
	Url  string `json:"url" bson:"url"`
	Date int64  `json:"date,omitempty" bson:"date,omitempty"`
}

type VerticeSchema struct {
	ID    primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	From  string              `json:"from" bson:"from"`
	To    string              `json:"to" bson:"to"`
	Infos []VerticeInfoSchema `json:"infos" bson:"infos"`
}
