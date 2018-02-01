package plz

import (
	"github.com/v2pro/plz/counselor"
	"github.com/v2pro/plz/pnp"
	"github.com/v2pro/plz/witch"
	"net/http"
)

// who am i, set externally
var PrimaryIP string
var PrimaryPort int
var Service string
var Cluster string

// to be set externally, additional info about this process
var ProcessInfo = map[string]interface{}{}

var PingUrl = ""

// will make counselor externally available after PlugAndPlay
var ExportCounselor = true

// will make witch externally available after PlugAndPlay
var ExportWitch = false

// PlugAndPlay will register the process into the grid
func PlugAndPlay() {
	ProcessInfo["Service"] = Service
	ProcessInfo["Cluster"] = Cluster
	mux := &http.ServeMux{}
	if ExportWitch {
		witch.Start("")
		mux.Handle("/witch/", witch.Mux)
	}
	if ExportCounselor {
		mux.Handle("/counselor/", counselor.Mux)
	}
	mux.HandleFunc("/", func(respWriter http.ResponseWriter, req *http.Request) {
		respWriter.Write([]byte(req.URL.String() + " not found"))
	})
	if PingUrl != "" {
		pnp.Start(PingUrl, ProcessInfo, mux)
	}
}
