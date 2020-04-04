package flopper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Flopper interface {
	http.Handler
	RegisterMetrics()
	Flops() int
}

type flopper struct {
	flops        prometheus.Counter
	illegalFlops prometheus.Counter
}

func New() Flopper {
	f := &flopper{}
	f.flops = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "flops_total",
			Help: "a counter of successful flops",
		},
	)
	f.illegalFlops = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "illegal_flops_total",
			Help: "a counter of times hops that were requested but not executed",
		},
	)
	return f
}

func (f *flopper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	flops := qs.Get("flops")
	if flops == "" {
		flops = "1"
	}
	n, err := strconv.Atoi(flops)
	if err != nil {
		f.illegalFlops.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n <= 0 || n >= 500 {
		f.illegalFlops.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i := 0; i < n; i++ {
		f.flops.Inc()
		fmt.Fprintf(w, "flop")
	}
}

func (f *flopper) RegisterMetrics() {
	prometheus.Register(f.flops)
	prometheus.Register(f.illegalFlops)
}

func (f *flopper) Flops() int {
	return 0 // how do i reach in?
}
