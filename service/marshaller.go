package service

import (
	"github.com/v2pro/plz/countlog"
)

type Response struct {
	Object interface{}
	Error  error
}

type Marshaller interface {
	Marshal(ctx *countlog.Context, output interface{}, obj interface{}) error
}
