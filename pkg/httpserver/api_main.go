package httpserver

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/NonLogicalDev/smartthings.influxdb-logger/pkg/tsdb"
	"go.uber.org/zap"
)

type SMTHandler struct {
	tsdb *tsdb.Client
}

func NewSMTHandler(influxUrl string) (*SMTHandler, error) {
	tsdbClient, err := tsdb.NewTSDBClient(influxUrl, "s")
	if err != nil {
		return nil, err
	}

	return &SMTHandler{
		tsdb: tsdbClient,
	}, nil
}

func (h *SMTHandler) HandleRequest(x *RequestContext) {
	h.handleSmtData(x)
}

func (h *SMTHandler) handleSmtData(x *RequestContext) {
	var msg SMTMessage
	x.OnError(400, x.ReadJSON(&msg))

	switch msg.Type {
	case "event":
		x.Log = x.Log.With(zap.String("message-type", "event"))
		x.Log.Info("received event request")

		var event SMTMessageEventData
		x.OnError(400, msg.Populate(&event))

		x.Log.Debug("received event", zap.Any("event", event))

		h.handleSmtEventData(x, event)

		return
	case "state":
		x.Log = x.Log.With(zap.String("message-type", "state"))
		x.Log.Info("received state request")

		var state string
		x.OnError(400, msg.Populate(&state))

		x.Log.Info("received state", zap.Any("state", state))
		return
	case "subscribe":
		x.Log = x.Log.With(zap.String("message-type", "subscribe"))
		x.Log.Info("received subscription request")

		var subscriptions []interface{}
		x.OnError(400, msg.Populate(&subscriptions))

		x.Log.Info("received subscriptions", zap.Any("subscriptions", subscriptions))
		return
	}

	x.WriteStatus(400)
}

func (h *SMTHandler) handleSmtEventData(x *RequestContext, evt SMTMessageEventData) {
	// Fetch dbName from the path.
	dbName := filepath.Base(x.Req.URL.Path)

	// Chop off extra digits from data returned from SMT
	// Example: 1556002854(442)
	date := time.Unix(evt.TS/1000, evt.TS%1000)

	x.Log = x.Log.Named("EventDataHandler").With(
		zap.String("db-name", dbName),
		zap.String("device-label", evt.Device.Label),
		zap.String("device-id", evt.Device.Id),
		zap.String("metric-name", evt.Metric.Name),
		zap.String("metric-value", evt.Metric.StrValue),
		zap.Time("date", date),
	)

	x.Log.Info("processing event request")

	switch evt.Metric.Name {
	case "status":
	case "primaryStatus":
	case "secondaryStatus":
		return
	}

	fields := DecodeValueToFields(evt)
	x.Log.Debug("parsed request", zap.Any("fields", fields))




	var err error
	err = h.tsdb.WriteMetrics(dbName,
		tsdb.Metric{
			TS:   date,
			Name: evt.Metric.Name,
			Tags: map[string]string{
				"device":    evt.Device.Label,
				"device-id": evt.Device.Id,
			},
			Values: fields.Export(),
		},
	)
	if err != nil {
		x.OnError(500, fmt.Errorf("failed writing to influxdb: %w", err))
	}
	x.WriteStatus(200)
}
