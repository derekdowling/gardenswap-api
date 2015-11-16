package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

func combineURL(baseURL string, path string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	path = strings.TrimPrefix(path, "/")
	return strings.Join([]string{baseURL, path}, "/")
}

func testRequest(data []byte) *http.Request {
	// build the fandangled io.ReadCloser
	reader := bytes.NewReader(data)
	closer := ioutil.NopCloser(reader)
	return &http.Request{Body: closer}
}
