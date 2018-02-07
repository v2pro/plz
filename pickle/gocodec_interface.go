package pickle

import (
	"unsafe"
	"reflect"
)

var typeSampleEface emptyInterface

func init() {
	var typ interface{} = reflect.TypeOf("")
	typeSampleEface = *(*emptyInterface)(unsafe.Pointer(&typ))
}

type interfaceEncoder struct {
	BaseCodec
	cfg *frozenConfig
}

func (encoder *interfaceEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	eface := (*emptyInterface)(ptr)
	valType := pTypeToType(eface.typ)
	elemEncoder, err := encoderOfType(encoder.cfg, valType)
	if err != nil {
		panic(err)
	}
	origCursor := stream.cursor
	offset := uintptr(len(stream.buf)) - stream.cursor
	stream.cursor = uintptr(len(stream.buf))
	valAsSlice := ptrAsBytes(int(valType.Size()), eface.word)
	stream.buf = append(stream.buf, valAsSlice...)
	elemEncoder.Encode(eface.word, stream)
	pwEface := unsafe.Pointer(&stream.buf[origCursor])
	wEface := (*writableEface)(pwEface)
	wEface.word = offset
}

type interfaceDecoder struct {
	BaseCodec
	cfg *frozenConfig
}

func (decoder *interfaceDecoder) Decode(iter *Iterator) {
	pObj := unsafe.Pointer(&iter.cursor[0])
	eface := (*emptyInterface)(pObj)
	valType := pTypeToType(eface.typ)
	elemDecoder, err := decoderOfType(decoder.cfg, valType)
	if err != nil {
		panic(err)
	}
	elemDecoder.DecodeEmptyInterface(eface, iter)
}

func (decoder *interfaceDecoder) HasPointer() bool {
	return true
}

func pTypeToType(ptr unsafe.Pointer) reflect.Type {
	typeEface := typeSampleEface
	typeEface.word = ptr
	typeObj := *(*interface{})(unsafe.Pointer(&typeEface))
	return typeObj.(reflect.Type)
}
