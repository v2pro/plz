package plzio

import (
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
)

type Unmarshaller interface {
	Unmarshal(ctx *countlog.Context, obj interface{}, input interface{}) error
}

type jsoniterUnmarshaller struct {
	protoInterface emptyInterface
}

func NewJsoniterUnmarshaller() Unmarshaller {
	return &jsoniterUnmarshaller{}
}

func (unmarshaller *jsoniterUnmarshaller) Unmarshal(ctx *countlog.Context, obj interface{}, input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), obj)
}
