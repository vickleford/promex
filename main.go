package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vickleford/promex/flipper"
	"github.com/vickleford/promex/flopper"
)

func main() {
	flipper := flipper.New()

	flopper := flopper.New()
	flopper.RegisterMetricsTo(prometheus.DefaultRegisterer)

	mux := http.NewServeMux()
	mux.Handle("/flipper", flipper)
	mux.Handle("/flopper", flopper)
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:         "0.0.0.0:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		Handler:      mux,
	}

	log.Printf("Starting server...")
	err := server.ListenAndServe()
	log.Fatalln(err)
}
