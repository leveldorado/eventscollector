package store

import (
	"net/http"

	"github.com/osipchuk/eventscollector/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoEventCounterStore struct {
	*MgoStore
}

const (
	EVENT_COUNTER_COLLECTION_NAME = "event_counter"
)

func NewMgoEventCounterStore(sess *mgo.Session) *MgoEventCounterStore {
	return &MgoEventCounterStore{NewMgoStore(sess, EVENT_COLLECTOR_DATABASE_NAME, EVENT_COUNTER_COLLECTION_NAME)}
}

func (s *MgoEventCounterStore) Increase(eventType string) StoreChannel {
	ch := make(StoreChannel, 1)
	go func() {
		query := bson.M{"_id": eventType}
		inc := bson.M{"$inc": bson.M{"counter": 1}}
		res := &StoreResult{Data: eventType}
		if _, err := s.collection.Upsert(query, inc); err != nil {
			res.Err = model.NewAppError(model.ERR_ID_STORE_ERROR, err.Error(), http.StatusInternalServerError, s.Increase)
		}
		ch <- res
		close(ch)
	}()
	return ch
}
