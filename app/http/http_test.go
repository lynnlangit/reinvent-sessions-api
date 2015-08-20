package http

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRequestGetParam(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo")

	actual, found := RequestGetParam(req, "q")
	if !found {
		t.Error("Could not find the parameter.")
		return
	}
	expected := "foo"
	if actual != expected {
		t.Errorf("Unexpected result. Expected %v, but got %v", expected, actual)
		return
	}
}
