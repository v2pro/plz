package witch

import (
	"net/http"
	"github.com/v2pro/plz/countlog"
	"github.com/json-iterator/go"
	"unsafe"
	"reflect"
)

var stateExporterType = reflect.TypeOf((*countlog.StateExporter)(nil)).Elem()

func init() {
	jsoniter.RegisterExtension(&stateExporterExtension{})
}

func exportState(respWriter http.ResponseWriter, req *http.Request) {
	encoder := jsoniter.NewEncoder(respWriter)
	err := encoder.Encode(countlog.StateExporters)
	if err != nil {
		countlog.Error("event!failed to export states", "err", err)
	}
}

type stateExporterExtension struct {
	jsoniter.DummyExtension
}

func (extension *stateExporterExtension) CreateEncoder(typ reflect.Type) jsoniter.ValEncoder {
	if typ == stateExporterType {
		// will use nonEmptyInterfaceCodec
		return nil
	}
	if reflect.PtrTo(typ).Implements(stateExporterType) {
		templateInterface := reflect.New(typ).Interface()
		var encoder jsoniter.ValEncoder = &stateExporterEncoder{
			templateInterface: extractInterface(templateInterface),
		}
		return encoder
	}
	if typ.Implements(stateExporterType) {
		templateInterface := reflect.New(typ).Elem().Interface()
		var encoder jsoniter.ValEncoder = &stateExporterEncoder{
			templateInterface: extractInterface(templateInterface),
		}
		if typ.Kind() == reflect.Ptr {
			encoder = &optionalEncoder{encoder}
		}
		return encoder
	}
	return nil
}

type stateExporterEncoder struct {
	templateInterface emptyInterface
}

func (encoder *stateExporterEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (encoder *stateExporterEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	templateInterface := encoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	stateExporter, ok := (*realInterface).(countlog.StateExporter)
	if !ok {
		stream.WriteVal(nil)
		return
	}
	state := stateExporter.ExportState()
	if state != nil {
		state["__object_address__"] = uintptr(ptr)
	}
	stream.WriteVal(state)
}

func (encoder *stateExporterEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, encoder)
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

type optionalEncoder struct {
	valueEncoder jsoniter.ValEncoder
}

func extractInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}

func (encoder *optionalEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		encoder.valueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (encoder *optionalEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, encoder)
}

func (encoder *optionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}
