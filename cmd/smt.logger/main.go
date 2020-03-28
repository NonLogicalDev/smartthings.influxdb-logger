package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/NonLogicalDev/smartthings.influxdb-logger/pkg/httpserver"
	"go.uber.org/zap"
)

var (
	listenURL string
	influxURL string
)

func init() {
	flag.StringVar(&listenURL,
		"listen", "", // "0.0.0.0:5555"
		"HOST:PORT on which to listen. env:(SMT_LISTEN)")
	flag.StringVar(&influxURL,
		"influx", "", // "http://0.0.0.0:8086"
		"URL location of the influxdb server. env:(SMT_INFLUX_URL)")
	flag.Parse()

	var ok bool
	if listenURL == "" {
		listenURL, ok = os.LookupEnv("SMT_LISTEN")
		if !ok {
			listenURL = "0.0.0.0:5555"
		}
	}

	if influxURL == "" {
		influxURL, ok  = os.LookupEnv("SMT_INFLUX_URL")
		if !ok {
			influxURL = "http://0.0.0.0:8086"
		}
	}
}

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	httpserver.RegisterHandlers(log, "/", influxURL, http.DefaultServeMux)

	log.Info(fmt.Sprint("Listening on: ", listenURL))
	log.Info(fmt.Sprint("InfluxDB URL: ", influxURL))

	if err := http.ListenAndServe(listenURL, nil); err != nil {
		panic(err)
	}
}
