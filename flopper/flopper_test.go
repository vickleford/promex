package flopper_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/vickleford/promex/flopper"
)

func TestOneFlop(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL)
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wanted status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unexpected error reading body: %s", err.Error())
	}
	if content := string(body); content != "flop" {
		t.Errorf("Unexpected content. Read '%s' Wanted '%s'", content, "flop")
	}
}

func TestFiveFlops(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=5")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wanted status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unexpected error reading body: %s", err.Error())
	}
	if content := string(body); content != "flopflopflopflopflop" {
		t.Errorf("Unexpected content. Read '%s' Wanted '%s'", content, "flopflopflopflopflop")
	}
}

func TestZeroFlops(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=0")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestNegativeFlops(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=-1")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestNonNumericFlops(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=lots")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestTooManyFlops(t *testing.T) {
	fl := flopper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=500")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestMetricRegistration(t *testing.T) {
	fl := flopper.New()
	fl.RegisterMetrics()
}

func TestSingleFlopCounterIncrements(t *testing.T) {
	reg := prometheus.NewRegistry()
	fl := flopper.New()
	fl.RegisterMetricsTo(reg)

	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL)
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	expected := strings.NewReader(`
# HELP flops_total a counter of successful flops
# TYPE flops_total counter
flops_total 1
# HELP illegal_flops_total a counter of times hops that were requested but not executed
# TYPE illegal_flops_total counter
illegal_flops_total 0
`)
	err = testutil.GatherAndCompare(reg, expected)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestManyFlopsCounterIncrements(t *testing.T) {
	reg := prometheus.NewRegistry()
	fl := flopper.New()
	fl.RegisterMetricsTo(reg)

	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=40")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	expected := strings.NewReader(`
# HELP flops_total a counter of successful flops
# TYPE flops_total counter
flops_total 40
# HELP illegal_flops_total a counter of times hops that were requested but not executed
# TYPE illegal_flops_total counter
illegal_flops_total 0
`)
	err = testutil.GatherAndCompare(reg, expected)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestSingleRequestForManyIllegalFlopsOnlyCountsAsOneIllegalFlop(t *testing.T) {
	reg := prometheus.NewRegistry()
	fl := flopper.New()
	fl.RegisterMetricsTo(reg)

	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flops=516")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	expected := strings.NewReader(`
# HELP flops_total a counter of successful flops
# TYPE flops_total counter
flops_total 0
# HELP illegal_flops_total a counter of times hops that were requested but not executed
# TYPE illegal_flops_total counter
illegal_flops_total 1
`)
	err = testutil.GatherAndCompare(reg, expected)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}
