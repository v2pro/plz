package nfmt

import (
	"errors"
	"fmt"
	"time"
)

var Formats = []Format{&TimeFormat{}}

type Format interface {
	FormatterOf(cfg map[string]interface{}, targetKey string, sample []interface{}) (Formatter, error)
}

func formatterOf(cfg map[string]interface{}, targetKey string, sample []interface{}) (ret Formatter, err error) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			err = fmt.Errorf("formatterOf panic: %v", recovered)
		}
	}()
	for _, format := range Formats {
		formatter, err := format.FormatterOf(cfg, targetKey, sample)
		if err != nil {
			return nil, err
		}
		if formatter != nil {
			return formatter, nil
		}
	}
	return nil, nil
}

type TimeFormat struct {
}

func (format *TimeFormat) FormatterOf(cfg map[string]interface{}, targetKey string, sample []interface{}) (Formatter, error) {
	if cfg["format"] != "time" {
		return nil, nil
	}
	for i := 0; i < len(sample); i+= 2{
		key := sample[i].(string)
		if key == targetKey {
			_, isTime := sample[i+1].(time.Time)
			if !isTime {
				return nil, fmt.Errorf("%s is not time.Time", targetKey)
			}
			layout, isString := cfg["layout"].(string)
			if !isString {
				return nil, errors.New("missing layout in format config")
			}
			return &timeFormatter{
				idx: i+1,
				layout: layout,
			}, nil
		}
	}
	return nil, fmt.Errorf("%s not found in properties", targetKey)
}

type timeFormatter struct {
	idx    int
	layout string
}

func (formatter *timeFormatter) Format(space []byte, kv []interface{}) []byte {
	return kv[formatter.idx].(time.Time).AppendFormat(space, formatter.layout)
}
