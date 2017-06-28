package srv

import (
	"context"
)

type HandleMethod func(ctx context.Context, request interface{}) (response interface{}, err error)
type StopServer func()
type StartServer func(server *Server) (StopServer, error)

var executors []StartServer

func AddServerExecutor(executor StartServer) {
	executors = append(executors, executor)
}

func BuildServer(kv ...interface{}) *Server {
	return &Server{toMap(kv), []map[string]interface{}{}, []*Server{}}
}

func toMap(kv []interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

type Server struct {
	Properties map[string]interface{}
	Methods    []map[string]interface{}
	SubServers []*Server
}

func (server *Server) SubServer(subServerName string, kv ...interface{}) *Server {
	subServerProps := toMap(kv)
	subServerProps["name"] = subServerName
	subServer := &Server{subServerProps, []map[string]interface{}{}, []*Server{}}
	server.SubServers = append(server.SubServers, subServer)
	return subServer
}

func (server *Server) Method(methodName string, kv ...interface{}) *Server {
	methodProps := toMap(kv)
	methodProps["name"] = methodName
	server.Methods = append(server.Methods, methodProps)
	return server
}

func (server *Server) Start() (StopServer, error) {
	stopFuncs := []StopServer{}
	for _, executor := range executors {
		stopFunc, err := executor(server)
		if err != nil {
			for _, stopFunc := range stopFuncs {
				stopFunc()
			}
			return nil, err
		}
		if stopFunc != nil {
			stopFuncs = append(stopFuncs, stopFunc)
		}
	}
	if len(stopFuncs) == 0 {
		panic("no executor can start this server")
	}
	return func() {
		for _, stopFunc := range stopFuncs {
			stopFunc()
		}
	}, nil
}
