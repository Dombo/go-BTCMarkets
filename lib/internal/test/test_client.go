package test

import (
	"BTCMarkets/lib"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"testing"
)

type testWriter struct {
	t *testing.T
}

func (w testWriter) Write(p []byte) (int, error) {
	w.t.Logf("%s", p)

	return len(p), nil
}

func Client(t *testing.T) *BTCMarkets.Client {
	return client(t, "", "")
}

func client(t *testing.T, accessKey string, privateKey string) *BTCMarkets.Client {
	transport := &http.Transport{
		DialTLS: func(network, _ string) (net.Conn, error) {
			addr := server.Listener.Addr().String()
			return tls.Dial(network, addr, &tls.Config{
				InsecureSkipVerify: true,
			})
		},
	}
	client := BTCMarkets.New(accessKey, privateKey)
	client.HTTPClient.Transport = transport
	//client.DebugLog = testLogger(t)

	return client
}

func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t: t}, "", 0)
}
