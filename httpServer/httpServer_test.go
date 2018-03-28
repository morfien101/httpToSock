package httpServer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	// Tracer is used to see if the functions are called correctly.
	tracer := make(map[string]bool)

	svrconfig := &ServerConfig{
		ListenAddress: ":8888",
		Routes: map[string]func() http.Handler{
			"/_status": func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					tracer["/_status"] = true
				})
			},
			"/healthz": func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					tracer["/healthz"] = true
				})
			},
		},
	}

	svr := NewServer(svrconfig)
	testStrings := []string{"/_status", "/healthz"}
	for _, v := range testStrings {
		request := httptest.NewRequest(http.MethodGet, v, nil)
		recorder := httptest.NewRecorder()
		svr.serveHTTP(recorder, request)
		if tracer[v] != true {
			t.Errorf("Tracer flag for %s is not true", v)
		}
	}
}
