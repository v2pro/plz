package plzio

import (
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
	"reflect"
)

type Unmarshaller interface {
	Unmarshal(ctx *countlog.Context, request interface{}, input interface{}) error
}

type jsoniterUnmarshaller struct {
	protoInterface emptyInterface
}

func NewJsoniterUnmarshaller() Unmarshaller {
	return &jsoniterUnmarshaller{}
}

func (unmarshaller *jsoniterUnmarshaller) Unmarshal(ctx *countlog.Context, request interface{}, input interface{}) error {
	reflect.ValueOf(request).Elem().Field(0).Set(reflect.ValueOf("hello"))
	return jsoniter.Unmarshal(input.([]byte), request)
}
