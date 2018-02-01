package pnp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func Start(pingUrl string, processInfo map[string]interface{}, handler http.Handler) {
	go maintainTunnelInBackground(pingUrl, processInfo, handler)
}

func maintainTunnelInBackground(pingUrl string, processInfo map[string]interface{}, handler http.Handler) {
	defer func() {
		countlog.LogPanic(recover())
	}()
	listener, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		countlog.Fatal("event!pnp.failed to start tunnel server", "err", err)
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port
	go serveTunnel(listener, handler)
	processInfo["tunnel_addr"] = fmt.Sprintf("127.0.0.1:%d", port)
	countlog.Info("event!pnp.tunnel started", "add", processInfo["tunnel_addr"])
	req, err := json.Marshal(map[string]interface{}{
		"ProcessId":   os.Getpid(),
		"ProcessInfo": processInfo,
	})
	if err != nil {
		countlog.Fatal("event!pnp.failed to marshal ping request",
			"err", err, "processInfo", processInfo)
		return
	}
	for {
		resp, err := http.Post(pingUrl, "application/json", bytes.NewBuffer(req))
		if err != nil {
			countlog.Warn("event!pnp.failed to ping", "err", err, "pingUrl", pingUrl)
			time.Sleep(time.Second * 5)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			countlog.Warn("event!pnp.failed to read ping response", "err", err, "pingUrl", pingUrl)
			time.Sleep(time.Second * 5)
			continue
		}
		var pingResp pingResponse
		err = json.Unmarshal(body, &pingResp)
		if err != nil {
			countlog.Warn("event!pnp.failed to unmarshal ping response",
				"err", err, "pingUrl", pingUrl, "body", body)
			time.Sleep(time.Second * 5)
			continue
		}
		if pingResp.Errno != 0 {
			countlog.Warn("event!pnp.ping response has non zero error number",
				"err", err, "pingUrl", pingUrl, "body", body)
			time.Sleep(time.Second * 5)
			continue
		}
		// ping again immediately
	}
}

type pingResponse struct {
	Errno int `json:"errno"`
}

func serveTunnel(listener net.Listener, handler http.Handler) {
	defer func() {
		countlog.LogPanic(recover())
	}()
	http.Serve(listener, handler)
	countlog.Info("event!pnp.tunnel quit")
}
