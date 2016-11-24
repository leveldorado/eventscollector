package model

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID        string    `json:"id" bson:"_id"`
	Type      string    `json:"type" bson:"type"`
	Ts        uint64    `json:"ts" bson:"ts"`
	Params    Params    `json:"params" bson:"params"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Params map[string]interface{}

func (doc *Event) PreSave() {
	doc.ID = bson.NewObjectId().Hex()
	doc.CreatedAt = time.Now()
}

func (doc Event) IsValid() *AppError {
	if !bson.IsObjectIdHex(doc.ID) {
		return NewAppError(ERR_ID_INTERNAL_SERVER_ERROR, "Invalid id", http.StatusInternalServerError, doc.IsValid)
	}
	if doc.Type == "" {
		return NewAppError(ERR_ID_INVALID_DATA, "Empty type", http.StatusBadRequest, doc.IsValid)
	}
	if doc.Ts == 0 {
		return NewAppError(ERR_ID_INVALID_DATA, "Empty ts", http.StatusBadRequest, doc.IsValid)
	}
	if doc.CreatedAt.IsZero() {
		return NewAppError(ERR_ID_INTERNAL_SERVER_ERROR, "Empty created_at", http.StatusBadRequest, doc.IsValid)
	}
	return nil
}
