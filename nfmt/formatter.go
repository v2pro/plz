package nfmt

import (
	"sync"
	"fmt"
	"github.com/v2pro/plz/nfmt/njson"
	"reflect"
	"encoding/json"
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
	levels     int
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
	if reflect.ValueOf(compiler.onByte).Pointer() == reflect.ValueOf(compiler.endState).Pointer() {
		return
	}
	if reflect.ValueOf(compiler.onByte).Pointer() == reflect.ValueOf(compiler.normal).Pointer() {
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
	case '{':
		compiler.start = i
		compiler.levels = 1
		compiler.onByte = compiler.afterLeftCurlyBrace
	default:
		compiler.invalidFormat(i, "verb unknown")
	}
}

func (compiler *formatCompiler) afterLeftCurlyBrace(i int, b byte) {
	if b == '"' {
		compiler.onByte = compiler.afterDoubleQuoteInCurlyBrace
	} else if b == '}' {
		compiler.levels -= 1
		if compiler.levels == 0 {
			var cfg map[string]interface{}
			err := json.Unmarshal([]byte(compiler.format[compiler.start:i+1]), &cfg)
			if err != nil {
				compiler.invalidFormat(i, "custom format should be specified in valid json: "+err.Error())
				return
			}
			formatter, err := formatterOf(cfg, compiler.lastKey, compiler.sample)
			if err != nil {
				compiler.invalidFormat(i, "custom format is invalid: "+err.Error())
				return
			}
			if formatter == nil {
				compiler.invalidFormat(i, "unknown custom format")
				return
			}
			compiler.formatters = append(compiler.formatters, formatter)
			compiler.start = i + 1
			compiler.onByte = compiler.normal
		}
	} else if b == '{' {
		compiler.levels += 1
	}
}

func (compiler *formatCompiler) afterDoubleQuoteInCurlyBrace(i int, b byte) {
	if b == '"' {
		compiler.onByte = compiler.afterLeftCurlyBrace
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
	compiler.onByte = compiler.endState
	compiler.formatters = []Formatter{invalidFormatter(fmt.Sprintf(
		"%s at %d %s", err, i, compiler.format))}
}

func (compiler *formatCompiler) endState(i int, b byte) {
}
