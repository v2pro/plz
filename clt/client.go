package clt

import (
	"context"
	"fmt"
)

type Client interface {
	Call(ctx context.Context, req interface{}, resp interface{}) error
}

var Providers = []func(serviceName string, methodName string, kv ...interface{}) Client{
}

func ClientOf(serviceName string, methodName string, kv ...interface{}) Client {
	for _, provider := range Providers {
		client := provider(serviceName, methodName, kv...)
		if client != nil {
			return client
		}
	}
	panic(fmt.Sprintf("no client for %s %s", serviceName, methodName))
}
