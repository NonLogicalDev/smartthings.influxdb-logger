package httpserver

import (
	"strings"
	"time"

	"github.com/NonLogicalDev/smartthings.influxdb-logger/pkg/tsdb"
)

type SMTHandler struct {
	dbName    string
	tsdb *tsdb.TSDBClient
}

func NewSMTHandler(influxUrl string) *SMTHandler {
	return &SMTHandler{
		dbName:    "smt",
		tsdb:      tsdb.NewTSDBClient(influxUrl),
	}
}

func (h *SMTHandler) HandleRequest(x *RequestContext) {
	h.handleSmtData(x)
}

func (h *SMTHandler) handleSmtData(x *RequestContext) {
	var msg SMTMessage
	x.OnError(400, x.ReadJSON(&msg))

	switch msg.Type {
	case "event":
		var event SMTMessageEventData
		x.Log("t:data:evt", "%s", string(msg.Data))
		x.OnError(400, msg.Populate(&event))
		h.handleSmtEventData(x, event)
		return
	case "state":
		var state string
		x.OnError(400, msg.Populate(&state))
		x.Log("t:state", "%s", state)
		return
	case "subscribe":
		var state []string
		x.OnError(400, msg.Populate(&state))
		x.Log("t:subscribe", "\n%s", strings.Join(state, "\n"))
		return
	}

	x.WriteStatus(400)
}

func (h *SMTHandler) handleSmtEventData(x *RequestContext, evt SMTMessageEventData) {
	// Chop off extra digits from data returned from SMT
	// Example: 1556002854(442)
	date := time.Unix(evt.TS/1000, evt.TS%1000)

	switch evt.Data.Metric {
	case "status":
	case "primaryStatus":
	case "secondaryStatus":
		return
	}

	x.Log("t:metric:evt", "label: %s, m: %v, v: %v, rv: %+v, date: (%v)",
		evt.Device.Label,
		evt.Data.Metric,
		evt.Data.StrValue,
		evt.Data,
		date,
	)

	fields := DecodeValueToFields(evt)
	x.Log("t:metric:fields", "%+v", fields)

	var err error
	err = h.tsdb.WriteMetrics(h.dbName,
		tsdb.Metric{
			TS:     date,
			Name: evt.Data.Metric,
			Tags: map[string]string{
				"device":    evt.Device.Label,
				"device-id": evt.Device.Id,
			},

			Values: fields,
		},
	)
	if err != nil {
		x.Log("t:metric:error", "%v", err)
		x.OnError(500, err)
	}

	x.Log("t:metric:status", "OK")
	x.WriteStatus(200)
}

