package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/NonLogicalDev/smartthings.influxdb-logger/pkg/httpserver"
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
	httpserver.RegisterHandlers("", influxURL, http.DefaultServeMux)

	fmt.Println("Listening on: ", listenURL)
	fmt.Println("InfluxDB URL: ", influxURL)
	if err := http.ListenAndServe(listenURL, nil); err != nil {
		panic(err)
	}
}
