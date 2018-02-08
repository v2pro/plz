package dump

import (
	"unsafe"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"reflect"
	"context"
	"encoding/json"
)

var addrMapKey = 2020020002
var dumper = jsonfmt.Config{
	Extensions: []jsonfmt.Extension{&dumpExtension{}},
}.Froze()

var efaceType = reflect.TypeOf(eface{})
var efaceEncoderInst = dumper.EncoderOf(reflect.TypeOf(eface{}))
var addrMapEncoderInst = dumper.EncoderOf(reflect.TypeOf(map[string]json.RawMessage{}))
var ptrEncoderInst = dumper.EncoderOf(reflect.TypeOf(uint64(0)))

type Var struct {
	Object interface{}
}

func (v Var) String() string {
	addrMap := map[string]json.RawMessage{}
	ctx := context.WithValue(context.Background(), addrMapKey, addrMap)
	rootPtr := unsafe.Pointer(&v.Object)
	output := efaceEncoderInst.Encode(ctx, nil, rootPtr)
	rootPtrStr := string(ptrEncoderInst.Encode(nil, nil, jsonfmt.PtrOf(rootPtr)))
	addrMap["__ptr__"] = json.RawMessage(`"` + rootPtrStr + `"`)
	addrMap[rootPtrStr] = json.RawMessage(output)
	output = addrMapEncoderInst.Encode(nil, nil, jsonfmt.PtrOf(addrMap))
	return string(output)
}

type dumpExtension struct {
}

func (extension *dumpExtension) EncoderOf(prefix string, valType reflect.Type) jsonfmt.Encoder {
	if valType == efaceType {
		return &efaceEncoder{}
	}
	return nil
}

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
}

func (encoder *efaceEncoder) Encode(ctx context.Context, space []byte, ptr unsafe.Pointer) []byte {
	space = append(space, `{"type":"`...)
	eface := (*eface)(ptr)
	valType := sampleType
	(*iface)(unsafe.Pointer(&valType)).data = eface.dataType
	space = append(space, valType.String()...)
	space = append(space, `","data":{"__ptr__":"`...)
	space = ptrEncoderInst.Encode(ctx, space, unsafe.Pointer(&eface.data))
	space = append(space, `"}}`...)
	return space
}
