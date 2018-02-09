package dump

import (
	"context"
	"unsafe"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"encoding/json"
)

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

type sliceEncoder struct {
	elemEncoder jsonfmt.Encoder
	elemSize    uintptr
}

func (encoder *sliceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	slice := (*sliceHeader)(ptr)
	space = append(space, `{"data":{"__ptr__":"`...)
	ptrStr := ptrToStr(uintptr(slice.data))
	space = append(space, ptrStr...)
	space = append(space, `"},"len":`...)
	space = intEncoderInst.Encode(ctx, space, unsafe.Pointer(&slice.len))
	space = append(space, `,"cap":`...)
	space = intEncoderInst.Encode(ctx, space, unsafe.Pointer(&slice.cap))
	space = append(space, `}`...)
	data := encoder.encodeData(ctx, nil, ptr)
	addrMap := ctx.Value(addrMapKey).(map[string]json.RawMessage)
	addrMap[ptrStr] = json.RawMessage(data)
	return space
}

func (encoder *sliceEncoder) encodeData(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	slice := (*sliceHeader)(ptr)
	space = append(space, '[')
	offset := uintptr(slice.data)
	for i := 0; i < slice.len; i++ {
		if i != 0 {
			space = append(space, ',')
		}
		space = encoder.elemEncoder.Encode(ctx, space, unsafe.Pointer(offset))
		offset += encoder.elemSize
	}
	space = append(space, ']')
	return space
}
