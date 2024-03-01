package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Set("p", "8080")
	flag.Parse()
	os.Exit(m.Run())
}

func TestPort(t *testing.T) {
	port := flag.Int("p", 1080, "the port to listen on")

	flag.Set("p", "8080")
	flag.Parse()

	if *port != 8080 {
		t.Errorf("expected port to be 8080, but got %d", *port)
	}
}

func TestOptions(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if allowOrigin := rr.Header().Get("Access-Control-Allow-Origin"); allowOrigin != "*" {
		t.Errorf("handler returned wrong Access-Control-Allow-Origin header: got %v want %v",
			allowOrigin, "*")
	}

	if allowMethods := rr.Header().Get("Access-Control-Allow-Methods"); allowMethods != "*" {
		t.Errorf("handler returned wrong Access-Control-Allow-Methods header: got %v want %v",
			allowMethods, "*")
	}

	if allowHeaders := rr.Header().Get("Access-Control-Allow-Headers"); allowHeaders != "*" {
		t.Errorf("handler returned wrong Access-Control-Allow-Headers header: got %v want %v",
			allowHeaders, "*")
	}
}

func TestInvalidPort(t *testing.T) {
	req, err := http.NewRequest("GET", "/1/package.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if body := rr.Body.String(); body != "目标端口号到少2位\n" {
		t.Errorf("handler returned wrong body: got %v want %v",
			body, "目标端口号到少2位\n")
	}
}
