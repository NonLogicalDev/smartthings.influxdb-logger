package httpserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var ExampleRequest = []byte(`
{
    "data": {
      "data": {
        "name": "tamper-alert.tamper",
        "str-value": "clear",
        "type": "boolean",
        "unit": null,
        "value": false
      },
      "date": "2020-03-27T02:42:42.876Z",
      "device": {
        "display": "Office :: ZMSensor",
        "groupId": null,
        "id": "a3f380f1-6ad9-4faa-967a-bb68a58466fa",
        "label": "Office :: ZMSensor",
        "make": null,
        "model": null,
        "name": "Zooz 4-in-1 Sensor",
        "type": "4-in-1-sensor-rev1"
      },
      "id": "f365fe1a-0eac-49de-9553-5157c82f28fc",
      "kind": "poll",
      "ts": 1585276962876
    },
    "type": "event"
  }
`)


func TestName(t *testing.T) {
	log, _ := zap.NewDevelopment()
	x := RequestContext{
		Context:      context.TODO(),
		Id:           1000,
		Log:          log,
		cacheBody:    ExampleRequest,
	}

	var msg SMTMessage
	err := x.ReadJSON(&msg)
	require.NoError(t, err)

	var evt SMTMessageEventData
	err = msg.Populate(&evt)
	require.NoError(t, err)

	out := DecodeValueToFields(evt).Export()
	require.NotNil(t, out)
}