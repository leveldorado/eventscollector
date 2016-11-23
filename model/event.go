package model

type Event struct {
	ID     string `json:"id" bson:"_id"`
	Type   string `json:"type" bson:"type"`
	Ts     uint64 `json:"ts" bson:"ts"`
	Params Params `json:"params" bson:"params"`
}

type Params map[string]interface{}
