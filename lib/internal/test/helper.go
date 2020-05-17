package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const testdataDir = "testdata"
const apiVersion = "v3"

func Testdata(t *testing.T, relativePath string) []byte {
	path := filepath.Join(testdataDir, relativePath)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("%s", err)
	}

	return b
}

func AssertEndpointCalled(t *testing.T, method, path string) {
	if Request.Method != method {
		t.Fatalf("expected %s, got %s", method, Request.Method)
	}

	path = fmt.Sprintf("/%s/%s", apiVersion, path)
	if escapedPath := Request.URL.EscapedPath(); escapedPath != path {
		t.Fatalf("expected %s, got %s", path, escapedPath)
	}
}
