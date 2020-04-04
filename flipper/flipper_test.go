package flipper_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vickleford/promex/flipper"
)

func TestOneFlip(t *testing.T) {
	fl := flipper.New()
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
	if content := string(body); content != "flip" {
		t.Errorf("Unexpected content. Read '%s' Wanted '%s'", content, "flip")
	}
}

func TestFiveFlips(t *testing.T) {
	fl := flipper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flips=5")
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
	if content := string(body); content != "flipflipflipflipflip" {
		t.Errorf("Unexpected content. Read '%s' Wanted '%s'", content, "flipflipflipflipflip")
	}
}

func TestZeroFlips(t *testing.T) {
	fl := flipper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flips=0")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestNegativeFlips(t *testing.T) {
	fl := flipper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flips=-1")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestNonNumericFlips(t *testing.T) {
	fl := flipper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flips=lots")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestTooManyFlips(t *testing.T) {
	fl := flipper.New()
	ts := httptest.NewServer(fl)
	defer ts.Close()

	c := http.Client{}
	resp, err := c.Get(ts.URL + "?flips=500")
	if err != nil {
		t.Errorf("That request should have worked, but: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
