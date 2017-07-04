package srv

import (
	"context"
)

type HandleMethod func(ctx context.Context, request interface{}) (response interface{}, err error)

type Notifier interface {
	Stop()
	IsStopped() bool
	Wait()
}

type StartServer func(server *Server) (Notifier, error)

var Executors []StartServer

func BuildServer(kv ...interface{}) *Server {
	return &Server{toMap(kv), []map[string]interface{}{}, []*Server{}}
}

func NewNotifier() Notifier {
	return &serverNotifier{
		ch: make(chan bool),
	}
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
	// sub server will map to Service concept in thrift
	// will map to url in http
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

func (server *Server) Start() (Notifier, error) {
	signals := []Notifier{}
	for _, executor := range Executors {
		signal := &serverNotifier{
			ch: make(chan bool),
		}
		signal, err := executor(server)
		if err != nil {
			for _, signal := range signals {
				signal.Stop()
			}
			return nil, err
		}
		if signal != nil {
			signals = append(signals, signal)
		}
	}
	if len(signals) == 0 {
		panic("no executor can start this server")
	}
	return &multipleServerNotifiers{signals}, nil
}
