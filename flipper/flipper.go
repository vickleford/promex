package flipper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	flips = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "flips_total",
			Help: "a counter of successful flips",
		},
	)
	illegalFlips = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "illegal_flips_total",
			Help: "a counter of times hops that were requested but not executed",
		},
	)
)

type flipper struct{}

func New() http.Handler {
	return &flipper{}
}

func (f *flipper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	flipsRequested := qs.Get("flips")
	if flipsRequested == "" {
		flipsRequested = "1"
	}
	n, err := strconv.Atoi(flipsRequested)
	if err != nil {
		illegalFlips.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n <= 0 || n >= 500 {
		illegalFlips.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i := 0; i < n; i++ {
		flips.Inc()
		fmt.Fprintf(w, "flip")
	}
}

func init() {
	prometheus.Register(flips)
	prometheus.Register(illegalFlips)
}
