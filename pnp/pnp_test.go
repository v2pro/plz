package pnp

import (
	"net/http"
	"testing"
	"time"
)

func Test_pnp(t *testing.T) {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(respWriter http.ResponseWriter, req *http.Request) {
		respWriter.Write([]byte(req.URL.Path))
	})
	Start("http://127.0.0.1:8318/ping", map[string]interface{}{}, mux)
	time.Sleep(time.Hour)
}
