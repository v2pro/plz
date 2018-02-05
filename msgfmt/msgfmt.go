package msgfmt

import (
	"io"
	"os"
	"sync"
	"unsafe"
	"fmt"
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 128)
	},
}

func Sprintf(format string, kvObj ...interface{}) string {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	buf := FormatterOf(format, kv).Format(nil, kv)
	sliceHeader := (*sliceHeader)(unsafe.Pointer(&buf))
	stringHeader := &stringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(stringHeader))
}

func Println(valuesObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&valuesObj)
	values := castEmptyInterfaces(uintptr(ptr))
	return fprintln(os.Stdout, values)
}

func Fprintln(writer io.Writer, valuesObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&valuesObj)
	values := castEmptyInterfaces(uintptr(ptr))
	return fprintln(writer, values)
}

func fprintln(writer io.Writer, values []interface{}) (int, error) {
	switch len(values) {
	case 0:
		return fmt.Println()
	case 1:
		return Fprintf(writer,"{single_value}\n", "single_value", values[0])
	default:
		return Fprintf(writer, "{multiple_values}\n", "multiple_values", values)
	}
}

func Printf(format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintf(os.Stdout, format, kv)
}

func Fprintf(writer io.Writer, format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintf(writer, format, kv)
}

func fprintf(writer io.Writer, format string, kv []interface{}) (int, error) {
	buf := bufPool.Get().([]byte)[:0]
	formatter := FormatterOf(format, kv)
	formatted := formatter.Format(buf, kv)
	n, err := writer.Write(formatted)
	bufPool.Put(formatted)
	return n, err
}

func castEmptyInterfaces(ptr uintptr) []interface{} {
	return *(*[]interface{})(unsafe.Pointer(ptr))
}

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
