package dump

import (
	"context"
	"unsafe"
	"encoding/json"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"reflect"
)

var strEncoderInst = jsonfmt.EncoderOf(reflect.TypeOf(""))

type stringHeader struct {
	data unsafe.Pointer
	len  int
}

type stringEncoder struct {
}

func (encoder *stringEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	header := (*stringHeader)(ptr)
	space = append(space, `{"data":{"__ptr__":"`...)
	ptrStr := ptrToStr(uintptr(header.data))
	space = append(space, ptrStr...)
	space = append(space, `"},"len":`...)
	space = intEncoderInst.Encode(ctx, space, unsafe.Pointer(&header.len))
	space = append(space, `}`...)
	elem := strEncoderInst.Encode(ctx, nil, ptr)
	addrMap := ctx.Value(addrMapKey).(map[string]json.RawMessage)
	addrMap[ptrStr] = json.RawMessage(elem)
	return space
}