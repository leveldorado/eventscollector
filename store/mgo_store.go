package store

import (
	"gopkg.in/mgo.v2"
)

type MgoStore struct {
	collection *mgo.Collection
}

func NewMgoStore(sess *mgo.Session, databaseName, collectionName string) *MgoStore {
	return &MgoStore{collection: sess.DB(databaseName).C(collectionName)}
}
