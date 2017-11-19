package witch

import (
	"net/http"
	"github.com/v2pro/plz/countlog"
	"github.com/json-iterator/go"
	"unsafe"
	"reflect"
	"io"
	"strconv"
)

var stateExporterType = reflect.TypeOf((*countlog.StateExporter)(nil)).Elem()
var myJson = jsoniter.Config{
	SortMapKeys:                   true,
	EscapeHTML:                    false,
	MarshalFloatWith6Digits:       true, // will lose precession
	ObjectFieldMustBeSimpleString: true, // do not unescape object field
}.Froze()

func init() {
	jsoniter.RegisterExtension(&stateExporterExtension{})
}

type exporting struct {
	encodedExporters map[uintptr][]byte // key is object address, value is encoded json
}

func exportState(respWriter http.ResponseWriter, req *http.Request) {
	setCurrentGoRoutineIsKoala()
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!plz.exporter.panic", "err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
	}()
	marshalState(countlog.StateExporters(), respWriter)
}

func marshalState(exporters map[string]countlog.StateExporter, writer io.Writer) {
	stream := myJson.BorrowStream(writer)
	defer myJson.ReturnStream(stream)
	exporting := &exporting{
		encodedExporters: map[uintptr][]byte{},
	}
	stream.Attachment = exporting
	stream.WriteObjectStart()
	stream.WriteObjectField("root")
	stream.WriteVal(exporters)
	for addr, encoded := range exporting.encodedExporters {
		stream.WriteMore()
		stream.WriteObjectField(strconv.Itoa(int(addr)))
		stream.Write(encoded)
	}
	stream.WriteObjectEnd()
	stream.Flush()
	if stream.Error != nil {
		countlog.Error("event!failed to export states", "err", stream.Error)
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
			encoder = &jsoniter.OptionalEncoder{encoder}
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
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	stream.WriteObjectField("__object_address__")
	stream.WriteVal(uintptr(ptr))
	stream.WriteObjectEnd()
	exporting := stream.Attachment.(*exporting)
	if _, found := exporting.encodedExporters[uintptr(ptr)]; !found {
		exporting.encodedExporters[uintptr(ptr)] = nil // placeholder
		state := stateExporter.ExportState()
		subStream := myJson.BorrowStream(nil)
		defer myJson.ReturnStream(subStream)
		subStream.Attachment = stream.Attachment
		subStream.WriteVal(state)
		exporting.encodedExporters[uintptr(ptr)] = append([]byte(nil), subStream.Buffer()...)
	}
}

func (encoder *stateExporterEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, encoder)
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func extractInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}
