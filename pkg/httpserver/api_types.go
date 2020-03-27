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

		Kind string `json:"kind"`

		Device SMTDeviceData            `json:"device"`
		Metric SMTMessageEventValueData `json:"metric"`
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
		Name string `json:"name"`

		Unit  string      `json:"unit"`
		Type  interface{} `json:"type"`
		Value interface{} `json:"value"`

		StrValue string `json:"str-value"`
	}
)

func (m SMTMessage) Populate(out interface{}) error {
	return json.Unmarshal(m.Data, out)
}

func DecodeValueToFields(evt SMTMessageEventData) *ValueOutput {
	values := ValueOutput{}
	values.SetType(evt.Kind)
	values.SetString(evt.Metric.StrValue)
	values.SetRaw(evt.Metric.Value)

	eData := evt.Metric
	switch eData.Type {
	case "numeric":
		v, err := ReadFloat(eData.Value)
		if err == nil {
			values.SetFloat(v)
		}
	case "boolean":
		v, ok := eData.Value.(bool)
		if ok {
			values.SetBool(v)
		}
	}
	return &values
}

func ReadFloat(val interface{}) (float64, error) {
	return strconv.ParseFloat(
		fmt.Sprintf("%v", val), 64,
	)
}

type ValueOutput struct {
	Type *string `json:"type,omitempty"`

	Raw *string `json:"raw,omitempty"`
	Str *string `json:"str,omitempty"`
	RawValue interface{} `json:"raw-value,omitempty"`

	RawNumeric *float64 `json:"raw-numeric,omitempty"`
	Numeric *float64 `json:"numeric,omitempty"`
	Bool *float64 `json:"bool,omitempty"`
}

func (v *ValueOutput) SetBool(in bool) {
	boolVal := float64(0)
	if in {
		boolVal = 1
	}

	v.Bool = &boolVal

	v.RawNumeric = &boolVal
	v.Numeric = &boolVal
}

func (v *ValueOutput) SetFloat(in float64) {
	boolVal := float64(0)
	if in > 0 {
		boolVal = 1
	}

	v.Bool = &boolVal

	v.RawNumeric = &in
	v.Numeric = &in
}

func (v *ValueOutput) SetString(in string) {
	v.Raw = &in
	v.Str = &in
}

func (v *ValueOutput) SetRaw(in interface{}) {
	//v.RawValue = in
}

func (v *ValueOutput) SetType(in string) {
	v.Type = &in
}

func (v *ValueOutput) Export() (out map[string]interface{}) {
	data, _ := json.Marshal(v)
	_ = json.Unmarshal(data, &out)
	return out
}
