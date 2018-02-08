package dump

import (
	"unsafe"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"reflect"
	"context"
)

var dumper = jsonfmt.Config{
	Extensions: []jsonfmt.Extension{&dumpExtension{}},
}.Froze()

type Var struct {
	Object interface{}
}

func (v Var) String() string {
	encoder := dumper.EncoderOf(reflect.TypeOf(eface{}))
	output := encoder.Encode(nil, nil, unsafe.Pointer(&v.Object))
	return string(output)
}

type dumpExtension struct {
}

func (extension *dumpExtension) EncoderOf(prefix string, valType reflect.Type) jsonfmt.Encoder {
	if valType == efaceType {
		return &efaceEncoder{ptrEncoder:jsonfmt.EncoderOf(reflect.TypeOf(uint64(0)))}
	}
	return nil
}

var efaceType = reflect.TypeOf(eface{})

type eface struct {
	dataType unsafe.Pointer
	data     unsafe.Pointer
}

type iface struct {
	itab unsafe.Pointer
	data unsafe.Pointer
}

var sampleType = reflect.TypeOf("")

type efaceEncoder struct {
	ptrEncoder jsonfmt.Encoder
}

func (encoder *efaceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, `{"type":"`...)
	eface := (*eface)(ptr)
	valType := sampleType
	(*iface)(unsafe.Pointer(&valType)).data = eface.dataType
	space = append(space, valType.String()...)
	space = append(space, `","data":{"__ptr__":"`...)
	space = encoder.ptrEncoder.Encode(ctx, space, unsafe.Pointer(&eface.data))
	space = append(space, `"}}`...)
	return space
}
