package witch

import "net/http"
import (
	"github.com/rakyll/statik/fs"
	"github.com/v2pro/plz/countlog"
	"io/ioutil"
	_ "github.com/v2pro/plz/witch/statik"
	"bytes"
)

//go:generate $GOPATH/bin/statik -src $PWD/webroot

var viewerHtml []byte

func initViewerHtml() error {
	if viewerHtml != nil {
		return nil
	}
	statikFS, err := fs.New()
	if err != nil {
		countlog.Error("event!witch.failed to load witch viewer web resource", "err", err)
		return err
	}
	var files = []string{"ide.html", "log-viewer.html", "filters.html"}
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

func StartViewer(addr string) {
	err := initViewerHtml()
	if err != nil {
		countlog.Error("event!witch.failed to init viewer html", "err", err)
		return
	}
	countlog.Info("event!witch.viewer started", "addr", addr)
	countlog.LogWriters = append(countlog.LogWriters, TheEventQueue)
	http.HandleFunc("/", homepage)
	http.HandleFunc("/more-events", moreEvents)
	http.ListenAndServe(addr, nil)
}

func homepage(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Write(viewerHtml)
}
