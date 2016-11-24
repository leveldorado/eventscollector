package httphandler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/osipchuk/eventscollector/collector"
	"github.com/osipchuk/eventscollector/model"
	"github.com/osipchuk/eventscollector/store"
)

type HTTPHandler struct {
	collector *collector.Collector
}

func NewHTTPHandler(c *collector.Collector) *HTTPHandler {
	return &HTTPHandler{collector: c}
}

const (
	eventTypeParam = "event_type"
)

func getEventTypeFromContext(cx *gin.Context) (eventType string, ok bool) {
	eventType = cx.Param(eventTypeParam)
	if eventType == "" {
		WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
			"Empty event_type", http.StatusBadRequest, getEventTypeFromContext))
		return
	}
	ok = true
	return
}

func (h *HTTPHandler) save(cx *gin.Context) {
	eventType, ok := getEventTypeFromContext(cx)
	if !ok {
		return
	}
	doc := &model.Event{}
	if err := cx.BindJSON(doc); err != nil {
		WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
			err.Error(), http.StatusBadRequest, h.save))
		return
	}
	doc.Type = eventType
	if err := h.collector.Save(doc); err != nil {
		WriteError(cx, err)
	}
	WriteSuccess(cx, nil)
}

const (
	rangeQueryParamName     = "range"
	intervalQueryParamName  = "interval"
	directionQueryParamName = "direction"
	limitQueryParamName     = "limit"
	offsetQueryParamName    = "offset"
)

const (
	maxLimit = 1000
)

func (h *HTTPHandler) selectEvents(cx *gin.Context) {
	eventType, ok := getEventTypeFromContext(cx)
	if !ok {
		return
	}
	var direction store.Direction
	direction.FromString(cx.Query(directionQueryParamName))
	from, to, ok := getRangeOrIntervalFromContext(cx)
	if !ok {
		return
	}
	limit, ok := parseUintAndWriteErrorIfNeed(cx, cx.Query(limitQueryParamName))
	if !ok {
		return
	}
	if limit > maxLimit {
		WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
			"limit is greater than max limit "+strconv.Itoa(maxLimit), http.StatusBadRequest, h.selectEvents))
		return
	}
	offset, _ := strconv.ParseUint(offsetQueryParamName, 10, 64)
	p := store.SelectParam{
		Type: eventType, Direction: direction,
		From: from, To: to,
		Limit: uint(limit), Offset: uint(offset),
	}
	resp, err := h.collector.Select(p)
	if err != nil {
		WriteError(cx, err)
		return
	}
	WriteSuccess(cx, resp)
}

func getRangeOrIntervalFromContext(cx *gin.Context) (from, to uint64, ok bool) {
	rangeParam := cx.Query(rangeQueryParamName)
	if rangeParam != "" {
		parts := strings.Split(rangeParam, ",")
		if len(parts) != 2 {
			WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
				"Invalid range param. Must be two uints separated by ,", http.StatusBadRequest,
				getRangeOrIntervalFromContext))
			return
		}
		from, ok = parseUintAndWriteErrorIfNeed(cx, parts[0])
		if !ok {
			return
		}
		to, ok = parseUintAndWriteErrorIfNeed(cx, parts[0])
		return
	}
	interval, err := time.ParseDuration(cx.Query(intervalQueryParamName))
	if err != nil {
		WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
			err.Error(), http.StatusBadRequest, getRangeOrIntervalFromContext))
		return
	}
	from, ok = uint64(time.Now().Add(-1*interval).Unix()), true
	return
}

func parseUintAndWriteErrorIfNeed(cx *gin.Context, param string) (value uint64, ok bool) {
	var err error
	value, err = strconv.ParseUint(param, 10, 64)
	if err != nil {
		WriteError(cx, model.NewAppError(model.ERR_ID_INVALID_DATA,
			err.Error(), http.StatusBadRequest, parseUintAndWriteErrorIfNeed))
		return
	}
	ok = true
	return
}

func WriteSuccess(cx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]bool{"success": true}
	}
	cx.JSON(http.StatusOK, data)
}

func WriteError(cx *gin.Context, err *model.AppError) {
	cx.JSON(err.Code, err)
}

func (h *HTTPHandler) RegisterRoutes(router gin.IRouter) {
	event := router.Group("/api/events-collector/:"+eventTypeParam, func(cx *gin.Context) {
		//TODO Check access
	})
	event.POST("", h.save)
	event.GET("", h.selectEvents)
}
