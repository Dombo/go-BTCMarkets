package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type request struct {
	Body                []byte
	ContentType, Method string
	URL                 *url.URL
}

var Request request

var server *httptest.Server

var responseBody []byte
var status int

func EnableServer(m *testing.M) {
	initAndStartServer()
	exitCode := m.Run()
	closeServer()

	os.Exit(exitCode)
}

func initAndStartServer() {
	server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Request = request{
			ContentType: r.Header.Get("Content-Type"),
			Method:      r.Method,
			URL:         r.URL,
		}

		var err error

		// Reading from the request body is fine, as it's not used elsewhere.
		// Server always returns fake data/testdata.
		Request.Body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		// status and responseBody are defined in returns.go.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if _, err := w.Write(responseBody); err != nil {
			panic(err.Error())
		}
	}))
}

func closeServer() {
	server.Close()
}

func WillReturn(b []byte, s int) {
	responseBody = b
	status = s
}

func WillReturnTestdata(t *testing.T, relativePath string, s int) {
	WillReturn(Testdata(t, relativePath), s)
}