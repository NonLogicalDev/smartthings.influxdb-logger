package httpserver

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	SMTMessage struct {
		Type string          `json:""`
		Data json.RawMessage `json:"data"`
	}

	SMTMessageEventData struct {
		TS   int64  `json:"ts"`
		Id   string `json:"id"`
		Date string `json:"date"`

		Mechanism string `json:"mechanism"`

		Device SMTDeviceData            `json:"device"`
		Data   SMTMessageEventValueData `json:"data"`
	}

	SMTDeviceData struct {
		Id string `json:"id"`

		Label       string `json:"label"`
		Name        string `json:"name"`
		DisplayName string `json:"dname"`

		Make  string `json:"make"`
		Model string `json:"model"`

		GroupID string `json:"group"`
	}

	SMTMessageEventValueData struct {
		Metric   string `json:"metric"`
		Unit     string `json:"unit"`
		StrValue string `json:"strValue"`

		RawType  string      `json:""`
		RawValue interface{} `json:"value"`

		DecodedType  string      `json:"decodedType"`
		DecodedValue interface{} `json:"decodedValue"`
	}
)

func (m SMTMessage) Populate(out interface{}) error {
	return json.Unmarshal(m.Data, out)
}

func DecodeValueToFields(evt SMTMessageEventData) map[string]interface{} {
	values := map[string]interface{}{
		"type": evt.Mechanism,
		"raw":  evt.Data.StrValue,
	}

	eData := evt.Data
	switch eData.RawType {
	case "numeric":
		v, err := ReadFloat(eData.RawValue)
		if err == nil {
			values["raw-numeric"] = v
		} else {
			values["err-raw"] = err.Error()
		}
	}

	switch eData.DecodedType {
	case "numeric":
		v, err := ReadFloat(eData.DecodedValue)
		if err == nil {
			values["numeric"] = v
		} else {
			values["err-decode"] = err.Error()
		}
	case "logical":
		v, err := ReadFloat(eData.DecodedValue)
		if err == nil {
			values["bool"] = v >= 1
			values["numeric"] = int64(v)
		} else {
			values["err-decode"] = err.Error()
		}
	}

	return values
}

func ReadFloat(val interface{}) (float64, error) {
	return strconv.ParseFloat(
		fmt.Sprintf("%v", val), 64,
	)
}
