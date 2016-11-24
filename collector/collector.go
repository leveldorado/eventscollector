package collector

import (
	"github.com/osipchuk/eventscollector/model"
	"github.com/osipchuk/eventscollector/store"
)

type Collector struct {
	eventSaver     EventSaver
	eventsSelector EventsSelector
	increaser      Increaser
}

func NewCollector(eventSaver EventSaver, eventsSelector EventsSelector, increaser Increaser) *Collector {
	return &Collector{eventSaver: eventSaver, eventsSelector: eventsSelector, increaser: increaser}
}

type EventSaver interface {
	Save(doc *model.Event) store.StoreChannel
}

type EventsSelector interface {
	Select(p store.SelectParam) store.StoreChannel
}

type Increaser interface {
	Increase(eventType string) store.StoreChannel
}

func (c *Collector) Save(doc *model.Event) *model.AppError {
	incChan := c.increaser.Increase(doc.Type)
	saveRes := <-c.eventSaver.Save(doc)
	if saveRes.Err != nil {
		go func() { _ = <-incChan }()
		return saveRes.Err
	}
	return (<-incChan).Err
}

func (c *Collector) Select(p store.SelectParam) ([]*model.Event, *model.AppError) {
	res := <-c.eventsSelector.Select(p)
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Data.([]*model.Event), nil
}
