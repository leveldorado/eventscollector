package config

import (
	"os"
)

func MustGetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Empty PORT")
	}
	return port
}

func MustGetCounterMgoHost() string {
	counterMgoHost := os.Getenv("COUNTER_MONGODB_HOST")
	if counterMgoHost == "" {
		panic("Empty COUNTER_MONGODB_HOST")
	}
	return counterMgoHost
}

func MustGetEventMgoHost() string {
	eventMgoHost := os.Getenv("EVENT_MONGODB_HOST")
	if eventMgoHost == "" {
		panic("Empty EVENT_MONGODB_HOST")
	}
	return eventMgoHost
}
