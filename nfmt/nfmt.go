package nfmt

import (
	"io"
	"os"
	"sync"
	"unsafe"
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

func Println(format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintln(os.Stdout, format, kv...)
}

func Fprintln(writer io.Writer, format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintln(writer, format, kv)
}

func fprintln(writer io.Writer, format string, kv ...interface{}) (int, error) {
	buf := bufPool.Get().([]byte)[:0]
	formatter := FormatterOf(format, kv)
	formatted := formatter.Format(buf, kv)
	formatted = append(formatted, '\n')
	n, err := writer.Write(formatted)
	bufPool.Put(formatted)
	return n, err
}

func Printf(format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintf(os.Stdout, format, kv...)
}

func Fprintf(writer io.Writer, format string, kvObj ...interface{}) (int, error) {
	ptr := unsafe.Pointer(&kvObj)
	kv := castEmptyInterfaces(uintptr(ptr))
	return fprintf(writer, format, kv)
}

func fprintf(writer io.Writer, format string, kv ...interface{}) (int, error) {
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
