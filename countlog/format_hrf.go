package countlog
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"reflect"
//)
//
//type HumanReadableFormat struct {
//	ContextPropertyNames []string
//	StringLengthCap      int
//}
//
//func (format *HumanReadableFormat) FormatLog(event Event) []byte {
//	msg := make([]byte, 0, 64)
//	ctx := format.describeContext(event)
//	if len(ctx) == 0 {
//		msg = append(msg, fmt.Sprintf(
//			"=== %s ===\n", event.Event)...)
//	} else {
//		msg = append(msg, fmt.Sprintf(
//			"=== [%s] %s ===\n", string(ctx), event.Event)...)
//	}
//	if len(event.Properties)%2 == 1 {
//		msg = append(msg, "wrong number of properties\n"...)
//		return msg
//	}
//	beforePropMsgLen := len(msg)
//	for i := 0; i < len(event.Properties); i += 2 {
//		k, _ := event.Properties[i].(string)
//		switch k {
//		case "", "ctx", "timestamp", "callee":
//			continue
//		}
//		v := event.Properties[i+1]
//		if event.Level <= LevelInfo && (k == "lineNumber" || k == "closedAt") {
//			continue
//		}
//		if k == "err" && v == nil {
//			continue
//		}
//		formattedV := formatV(v)
//		if formattedV == "" {
//			continue
//		}
//		if event.Level < LevelWarn && format.StringLengthCap > 0 {
//			lenCap := len(formattedV)
//			if format.StringLengthCap < lenCap {
//				lenCap = format.StringLengthCap
//				formattedV = formattedV[:lenCap] + "...more, capped"
//			}
//		}
//		msg = append(msg, k...)
//		msg = append(msg, ": "...)
//		msg = append(msg, formattedV...)
//		msg = append(msg, '\n')
//	}
//	noProp := len(msg) == beforePropMsgLen
//	if noProp && event.Level <= LevelTrace {
//		return nil
//	}
//	return msg
//}
//
//func formatV(v interface{}) string {
//	if v == nil {
//		return "<nil>"
//	}
//	switch typedV := v.(type) {
//	case []byte:
//		return string(encodeAnyByteArray(typedV))
//	case string:
//		return typedV
//	case json.RawMessage:
//		return string(typedV)
//	default:
//		err, _ := v.(error)
//		if err != nil {
//			return err.Error()
//		}
//		stringer, _ := v.(fmt.Stringer)
//		if stringer != nil {
//			return stringer.String()
//		}
//		goStringer, _ := v.(fmt.GoStringer)
//		if goStringer != nil {
//			return goStringer.GoString()
//		}
//		switch reflect.TypeOf(v).Kind() {
//		case reflect.Chan, reflect.Struct, reflect.Interface, reflect.Func, reflect.Array, reflect.Slice, reflect.Ptr:
//			return ""
//		}
//		return fmt.Sprintf("%v", typedV)
//	}
//}
//
//func (format *HumanReadableFormat) describeContext(event Event) []byte {
//	var msg []byte
//	ctx, _ := event.Get("ctx").(context.Context)
//	for _, propName := range format.ContextPropertyNames {
//		propValue := event.Get(propName)
//		if propValue == nil && ctx != nil {
//			propValue = ctx.Value(propName)
//		}
//		if propValue != nil {
//			if len(msg) > 0 {
//				msg = append(msg, ',')
//			}
//			msg = append(msg, propName...)
//			msg = append(msg, '=')
//			msg = append(msg, fmt.Sprintf("%v", propValue)...)
//		}
//	}
//	return msg
//}
//
//func isBinary(buf []byte) bool {
//	for _, b := range buf {
//		if b == '\r' || b == '\n' {
//			continue
//		}
//		if b < 32 || b > 127 {
//			return true
//		}
//	}
//	return false
//}
