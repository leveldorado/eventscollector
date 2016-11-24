package store

import (
	"github.com/osipchuk/eventscollector/model"
)

type StoreResult struct {
	Err  *model.AppError
	Data interface{}
}

type StoreChannel chan *StoreResult
