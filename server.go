package plz

import "github.com/v2pro/plz/server"

func Server(kv ...interface{}) *server.Server {
	return server.BuildServer(kv...)
}
