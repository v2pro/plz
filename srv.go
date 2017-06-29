package plz

import (
	"github.com/v2pro/plz/srv"
)

func Server(kv ...interface{}) *srv.Server {
	return srv.BuildServer(kv...)
}
