package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/osipchuk/eventscollector/collector"
	"github.com/osipchuk/eventscollector/config"
	"github.com/osipchuk/eventscollector/httphandler"
	"github.com/osipchuk/eventscollector/store"
	"gopkg.in/mgo.v2"
)

func Run() {
	//IDEA is to keep counter in one instance for all backends and shards
	counterSess, err := mgo.Dial(config.MustGetCounterMgoHost())
	if err != nil {
		panic(err)
	}
	defer counterSess.Close()
	eventSess, err := mgo.Dial(config.MustGetEventMgoHost())
	if err != nil {
		panic(err)
	}
	defer eventSess.Close()

	counterStore := store.NewMgoEventCounterStore(counterSess)
	eventStore := store.NewMgoEventStore(eventSess)

	eventsCollector := collector.NewCollector(eventStore, eventStore, counterStore)

	handler := httphandler.NewHTTPHandler(eventsCollector)

	router := gin.New()
	router.Use(gin.Recovery())

	handler.RegisterRoutes(router)

	listener := &http.Server{
		Addr:         `:` + config.MustGetPort(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router,
	}
	err = listener.ListenAndServe()
	log.Println("SERVER SHUTDOWN")
	if err != nil {
		log.Println(err)
	}
}
