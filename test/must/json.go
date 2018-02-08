package must

import (
	"encoding/json"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/test/testify/assert"
	"runtime"
	"reflect"
)

func JsonEqual(expected string, actual interface{}) {
	t := test.CurrentT()
	t.Helper()
	var expectedObj interface{}
	err := json.Unmarshal([]byte(expected), &expectedObj)
	if err != nil {
		t.Fatal("expected json is invalid: " + err.Error())
		return
	}
	var actualJson []byte
	switch actualVal := actual.(type) {
	case string:
		actualJson = []byte(actualVal)
	case []byte:
		actualJson = actualVal
	default:
		actualJson, err = json.Marshal(actual)
		t.Fatal("actual can not marshal to json: " + err.Error())
		return
	}
	var actualObj interface{}
	err = json.Unmarshal(actualJson, &actualObj)
	if err != nil {
		t.Log(string(actualJson))
		t.Fatal("actual json is invalid: " + err.Error())
		return
	}
	maskAnything(expectedObj, actualObj)
	if assert.Equal(t, expectedObj, actualObj) {
		return
	}
	t.Helper()
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("check failed")
		return
	}
	t.Fatal(test.ExtractFailedLines(file, line))
}

func maskAnything(expected interface{}, actual interface{}) {
	switch reflect.TypeOf(expected).Kind() {
	case reflect.Map:
		if reflect.ValueOf(actual).Kind() != reflect.Map {
			return
		}
		expectedVal := reflect.ValueOf(expected)
		actualVal := reflect.ValueOf(actual)
		keys := expectedVal.MapKeys()
		for _, key := range keys {
			elem := expectedVal.MapIndex(key).Interface()
			if elem == "ANYTHING" {
				actualVal.SetMapIndex(key, reflect.ValueOf("ANYTHING"))
				continue
			}
			maskAnything(elem, actualVal.MapIndex(key).Interface())
		}
	}
}
