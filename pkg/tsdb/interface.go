package tsdb

import (
	"time"

	influxdb "github.com/influxdata/influxdb1-client/v2"
)

type Client struct {
	defaultPrecision string
	influxClient     influxdb.Client
}

type Metric struct {
	TS   time.Time
	Name string
	Tags map[string]string

	Values map[string]interface{}
}

func NewTSDBClient(url string, defaultPrecision string) (*Client, error) {
	var err error

	c := Client{
		defaultPrecision: defaultPrecision,
	}

	c.influxClient, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: url,
	})

	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (ic *Client) WriteMetrics(db string, metrics ...Metric) error {
	bp, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Precision: ic.defaultPrecision,
		Database:  db,
	})

	var points []*influxdb.Point
	for _, m := range metrics {
		ipoint, err := influxdb.NewPoint(
			m.Name, m.Tags, m.Values, m.TS,
		)
		if err != nil {
			return nil
		}
		points = append(points, ipoint)
	}

	bp.AddPoints(points)
	return ic.influxClient.Write(bp)
}
