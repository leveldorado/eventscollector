package store

import (
	"net/http"

	"github.com/osipchuk/eventscollector/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoEventStore struct {
	*MgoStore
}

const (
	EVENT_COLLECTOR_DATABASE_NAME = "event_collector"
	EVENT_COLLECTION_NAME         = "event"
)

func NewMgoEventStore(sess *mgo.Session) *MgoEventStore {
	return &MgoEventStore{NewMgoStore(sess, EVENT_COLLECTOR_DATABASE_NAME, EVENT_COLLECTION_NAME)}
}

func (s *MgoEventStore) CreateIndexesIfNotExists() {
	s.collection.EnsureIndexKey("type")
	s.collection.EnsureIndexKey("-ts")
	s.collection.EnsureIndexKey("ts")
}

func (s *MgoEventStore) Save(doc *model.Event) StoreChannel {
	ch := make(StoreChannel, 1)
	go func() {
		res := &StoreResult{}
		defer func() {
			ch <- res
			close(ch)
		}()
		doc.PreSave()
		if err := doc.IsValid(); err != nil {
			res.Err = err
			return
		}
		if err := s.collection.Insert(doc); err != nil {
			res.Err = model.NewAppError(model.ERR_ID_STORE_ERROR, err.Error(), http.StatusInternalServerError, s.Save)
		}
		res.Data = doc
	}()
	return ch
}

type SelectParam struct {
	Type      string
	From      uint64
	To        uint64
	Limit     uint
	Offset    uint
	Direction Direction
}

type Direction int

const (
	DescDirection Direction = iota
	AscDirection
)

func (p SelectParam) getBsonQuery() bson.M {
	return bson.M{
		"type": p.Type,
		"ts": bson.M{
			"gt": p.From,
			"lt": p.From,
		},
	}
}

func (s *MgoEventStore) Select(p SelectParam) StoreChannel {
	ch := make(StoreChannel, 1)
	go func() {
		res := &StoreResult{}
		events := []*model.Event{}
		sort := "-ts"
		if p.Direction == AscDirection {
			sort = sort[1:]
		}
		query := p.getBsonQuery()
		if err := s.collection.Find(query).Sort(sort).Skip(int(p.Offset)).Limit(int(p.Limit)).All(&events); err != nil {
			res.Err = model.NewAppError(model.ERR_ID_STORE_ERROR, err.Error(), http.StatusInternalServerError, s.Select)
		}
		res.Data = events
		ch <- res
		close(ch)
	}()
	return ch
}
