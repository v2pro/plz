package service

import (
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
)

type Response struct {
	Object interface{}
	Error  error
}

type Marshaller interface {
	Marshal(ctx *countlog.Context, output interface{}, obj interface{}) error
}

type jsoniterResponseMarshaller struct {
}


func NewJsoniterResponseMarshaller() Marshaller {
	return &jsoniterResponseMarshaller{}
}

func (marshaller *jsoniterResponseMarshaller) Marshal(ctx *countlog.Context, output interface{}, obj interface{}) error {
	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)
	stream.WriteObjectStart()
	stream.WriteObjectField("errno")
	resp := obj.(Response)
	if resp.Error != nil {
		errno, _ := resp.Error.(errNo)
		if errno == nil {
			stream.WriteInt(1)
		} else {
			stream.WriteInt(errno.ErrorNumber())
		}
		stream.WriteMore()
		stream.WriteObjectField("errmsg")
		stream.WriteString(resp.Error.Error())
	} else {
		stream.WriteInt(0)
	}
	stream.WriteMore()
	stream.WriteObjectField("data")
	stream.WriteVal(resp.Object)
	stream.WriteObjectEnd()
	if stream.Error != nil {
		return stream.Error
	}
	ptrBuf := output.(*[]byte)
	*ptrBuf = append(([]byte)(nil), stream.Buffer()...)
	return nil
}

type errNo interface {
	ErrorNumber() int
}

type jsoniterMarshaller struct {
}

func NewJsoniterMarshaller() Marshaller {
	return &jsoniterMarshaller{}
}

func (marshaller *jsoniterMarshaller) Marshal(ctx *countlog.Context, output interface{}, obj interface{}) error {
	buf, err := jsoniter.Marshal(obj)
	if err != nil {
		return err
	}
	ptrBuf := output.(*[]byte)
	*ptrBuf = buf
	return nil
}
