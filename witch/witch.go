package witch

import "net/http"
import (
	"bytes"
	"github.com/rakyll/statik/fs"
	"github.com/v2pro/plz/countlog"
	_ "github.com/v2pro/plz/witch/statik"
	"io/ioutil"
)

var files = []string{
	"ide.html",
	"log-viewer.html", "filters.html", "data-sources.html", "columns.html",
	"state-viewer.html", "snapshots.html"}

//go:generate $GOPATH/bin/statik -src $PWD/webroot

var viewerHtml []byte

var Mux = &http.ServeMux{}

func init() {
	Mux.HandleFunc("/witch/more-events", moreEvents)
	Mux.HandleFunc("/witch/export-state", exportState)
	Mux.HandleFunc("/witch/", homepage)
}

func initViewerHtml() error {
	if viewerHtml != nil {
		return nil
	}
	statikFS, err := fs.New()
	if err != nil {
		countlog.Error("event!witch.failed to load witch viewer web resource", "err", err)
		return err
	}
	indexHtmlFile, err := statikFS.Open("/index.html")
	if err != nil {
		countlog.Error("event!witch.failed to open index.html", "err", err)
		return err
	}
	indexHtml, err := ioutil.ReadAll(indexHtmlFile)
	if err != nil {
		countlog.Error("event!witch.failed to read index.html", "err", err)
		return err
	}
	components := []byte{}
	for _, file := range files {
		f, err := statikFS.Open("/" + file)
		if err != nil {
			countlog.Error("event!witch.failed to open file", "err", err, "file", file)
			return err
		}
		fileHtml, err := ioutil.ReadAll(f)
		if err != nil {
			countlog.Error("event!witch.failed to read file html", "err", err, "file", file)
			return err
		}
		components = append(components, fileHtml...)
	}
	viewerHtml = bytes.Replace(indexHtml, []byte("{{ COMPONENTS }}"), components, -1)
	return nil
}

func Start(addr string) {
	err := initViewerHtml()
	if err != nil {
		countlog.Error("event!witch.failed to init viewer html", "err", err)
		return
	}
	countlog.Info("event!witch.viewer started", "addr", addr)
	countlog.LogWriters = append(countlog.LogWriters, theEventQueue)
	if addr != "" {
		go func() {
			setCurrentGoRoutineIsKoala()
			http.ListenAndServe(addr, Mux)
		}()
	}
}

func homepage(respWriter http.ResponseWriter, req *http.Request) {
	setCurrentGoRoutineIsKoala()
	respWriter.Write(viewerHtml)
}
