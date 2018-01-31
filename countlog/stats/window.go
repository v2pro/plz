package stats

import (
	"github.com/v2pro/plz/countlog/core"
	"unsafe"
	"strings"
	"reflect"
)

type Point struct {
	Event     string
	Timestamp int64
	Dimension map[string]string
	Value     float64
}

type Collector interface {
	Collect(event string, timestamp int64, dimension map[string]string, value float64)
}

type Window struct {
	MapMonoid
}

func newWindow() *Window {
	return &Window{
		MapMonoid: MapMonoid{},
	}
}

func (window *Window) Export(collector Collector) {
}

func (window *Window) Reset() {
}

type propIdx int

type dimensionExtractor interface{
	Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid
}

func newDimensionExtractor(site *core.EventSite) dimensionExtractor {
	var dimensionElems []string
	for i := 0; i < len(site.Sample); i += 2 {
		key := site.Sample[i].(string)
		if key == "dim" {
			dimensionElems = strings.Split(site.Sample[i+1].(string), ",")
		}
	}
	indices := make([]propIdx, 0, len(dimensionElems))
	for i := 0; i < len(site.Sample); i += 2 {
		key := site.Sample[i].(string)
		for _, dimension := range dimensionElems {
			if key == dimension {
				indices = append(indices, propIdx(i+1))
			}
		}
	}
	arrayType := reflect.ArrayOf(len(dimensionElems), reflect.TypeOf(""))
	arrayObj := reflect.New(arrayType).Elem().Interface()
	sampleInterface := *(*emptyInterface)(unsafe.Pointer(&arrayObj))
	if len(indices) == 0 {
		return &dimensionExtractor0{}
	}
	if len(indices) <= 2 {
		return &dimensionExtractor2{
			sampleInterface: sampleInterface,
			indices: indices,
		}
	}
	if len(indices) <= 4 {
		return &dimensionExtractor4{
			sampleInterface: sampleInterface,
			indices: indices,
		}
	}
	if len(indices) <= 8 {
		return &dimensionExtractor8{
			sampleInterface: sampleInterface,
			indices: indices,
		}
	}
	return &dimensionExtractorAny{
		sampleInterface: sampleInterface,
		indices: indices,
	}
}

type dimensionExtractor0 struct {
}

func (extractor *dimensionExtractor0) Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	elem := monoid[0]
	if elem == nil {
		elem = createElem()
		monoid[0] = elem
	}
	return elem
}

type dimensionExtractor2 struct {
	sampleInterface emptyInterface
	indices []propIdx
}
func (extractor *dimensionExtractor2) Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	dimensionArr := [2]string{}
	dimension := dimensionArr[:len(extractor.indices)]
	return extractDimension(extractor.sampleInterface, dimension,
		extractor.indices, event, monoid, createElem)
}

type dimensionExtractor4 struct {
	sampleInterface emptyInterface
	indices []propIdx
}
func (extractor *dimensionExtractor4) Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	dimensionArr := [4]string{}
	dimension := dimensionArr[:len(extractor.indices)]
	return extractDimension(extractor.sampleInterface, dimension,
		extractor.indices, event, monoid, createElem)
}

type dimensionExtractor8 struct {
	sampleInterface emptyInterface
	indices []propIdx
}
func (extractor *dimensionExtractor8) Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	dimensionArr := [8]string{}
	dimension := dimensionArr[:len(extractor.indices)]
	return extractDimension(extractor.sampleInterface, dimension,
		extractor.indices, event, monoid, createElem)
}

type dimensionExtractorAny struct {
	sampleInterface emptyInterface
	indices []propIdx
}
func (extractor *dimensionExtractorAny) Extract(event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	dimension := make([]string, len(extractor.indices))
	return extractDimension(extractor.sampleInterface, dimension,
		extractor.indices, event, monoid, createElem)
}

func extractDimension(
	sampleInterface emptyInterface, dimension []string, indices []propIdx,
	event *core.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
	for i, idx := range indices {
		dimension[i] = event.Properties[idx].(string)
	}
	dimensionInterface := sampleInterface
	dimensionInterface.word = unsafe.Pointer(&dimension[0])
	dimensionObj := castEmptyInterface(uintptr(unsafe.Pointer(&dimensionInterface)))
	elem := monoid[dimensionObj]
	if elem == nil {
		elem = createElem()
		monoid[dimensionObj] = elem
	}
	return elem
}

func castEmptyInterface(ptr uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(ptr))
}


// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}