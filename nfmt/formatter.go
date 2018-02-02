package nfmt

import (
	"sync"
	"fmt"
	"github.com/v2pro/plz/nfmt/njson"
	"reflect"
)

var formatterCache = &sync.Map{}

type Formatter interface {
	Format(space []byte, kv []interface{}) []byte
}

type Formatters []Formatter

func (formatters Formatters) Format(space []byte, kv []interface{}) []byte {
	for _, formatter := range formatters {
		space = formatter.Format(space, kv)
	}
	return space
}

func FormatterOf(format string, sample []interface{}) Formatter {
	formatterObj, found := formatterCache.Load(format)
	if found {
		return formatterObj.(Formatter)
	}
	formatter := compile(format, sample)
	formatterCache.Store(format, formatter)
	return formatter
}

type formatCompiler struct {
	sample     []interface{}
	format     string
	start      int
	lastKey    string
	onByte     func(int, byte)
	formatters []Formatter
}

func compile(format string, sample []interface{}) Formatter {
	compiler := &formatCompiler{
		sample: sample,
		format: format,
	}
	compiler.onByte = compiler.normal
	compiler.compile()
	return Formatters(compiler.formatters)
}

func (compiler *formatCompiler) compile() {
	format := compiler.format
	for i := 0; i < len(format); i++ {
		compiler.onByte(i, format[i])
	}
	if fmt.Sprintf("%v", compiler.onByte) == fmt.Sprintf("%v", compiler.normal) {
		compiler.formatters = append(compiler.formatters,
			fixedFormatter(compiler.format[compiler.start:len(format)]))
	} else {
		compiler.invalidFormat(len(format)-1, "verb not properly ended")
	}
}

func (compiler *formatCompiler) normal(i int, b byte) {
	format := compiler.format
	if format[i] == '%' {
		compiler.formatters = append(compiler.formatters,
			fixedFormatter(format[compiler.start:i]))
		compiler.onByte = compiler.afterPercent
	}
}

func (compiler *formatCompiler) afterPercent(i int, b byte) {
	if b == '(' {
		compiler.start = i + 1
		compiler.onByte = compiler.afterLeftBrace
	} else {
		compiler.invalidFormat(i, "expect left brace")
	}
}

func (compiler *formatCompiler) afterLeftBrace(i int, b byte) {
	if b == ')' {
		compiler.onByte = compiler.afterRightBrace
	}
}

func (compiler *formatCompiler) afterRightBrace(i int, b byte) {
	key := compiler.format[compiler.start:i-1]
	compiler.lastKey = key
	switch b {
	case 's':
		idx := compiler.findLastKey()
		if idx == -1 {
			compiler.invalidFormat(i, compiler.lastKey+" not found in args")
			return
		}
		sampleValue := compiler.sample[idx]
		switch sampleValue.(type) {
		case string:
			compiler.formatters = append(compiler.formatters, strFormatter(idx))
		case []byte:
			compiler.formatters = append(compiler.formatters, bytesFormatter(idx))
		default:
			compiler.formatters = append(compiler.formatters, &jsonFormatter{
				idx:     idx,
				encoder: njson.EncoderOf(reflect.TypeOf(sampleValue)),
			})
		}
		compiler.start = i + 1
		compiler.onByte = compiler.normal
	default:
		compiler.invalidFormat(i, "verb unknown")
	}
}

func (compiler *formatCompiler) findLastKey() int {
	for i := 0; i < len(compiler.sample); i += 2 {
		key := compiler.sample[i].(string)
		if key == compiler.lastKey {
			return i + 1
		}
	}
	return -1
}

func (compiler *formatCompiler) invalidFormat(i int, err string) {
	compiler.onByte = func(i int, b byte) {
	}
	compiler.formatters = []Formatter{invalidFormatter(fmt.Sprintf(
		"%s at %d %s", err, i, compiler.format))}
}
