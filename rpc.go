package plz

import (
	"github.com/v2pro/plz/logging"
	"github.com/v2pro/plz/srv"
	"github.com/v2pro/plz/clt"
)

func LoggerOf(loggerKv ...interface{}) logging.Logger {
	return logging.LoggerOf(loggerKv...)
}

func ClientOf(serviceName string, methodName string, kv ...interface{}) clt.Client {
	return clt.ClientOf(serviceName, methodName, kv...)
}

func BuildServer(kv ...interface{}) *srv.Server {
	return srv.BuildServer(kv...)
}
