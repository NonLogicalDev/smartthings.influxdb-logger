package tsdb

import (
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"time"
)

type TSDBClient struct {
	influxClient influxdb.Client
}

type Metric struct {
	TS time.Time
	Name string
	Tags map[string]string

	Values map[string]interface{}
}

func NewTSDBClient(url string) *TSDBClient {
	c := TSDBClient{}

	var err error
	c.influxClient, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: url,
	})

	if err != nil {
		panic(err)
	}

	return &c
}

func (ic *TSDBClient) WriteMetrics(db string, metrics ...Metric) (error) {
	bp, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Precision: "s",
		Database: db,
	})
	var ipoints []*influxdb.Point
	for _, m := range metrics {
		ipoint, err := influxdb.NewPoint(
			m.Name, m.Tags, m.Values, m.TS,
		)
		if err != nil {
			return nil
		}
		ipoints = append(ipoints, ipoint)
	}

	bp.AddPoints(ipoints)
	return ic.influxClient.Write(bp)
}
