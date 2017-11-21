package counselor

import (
	"github.com/v2pro/plz/countlog"
	"encoding/json"
	"fmt"
	"net/http"
)

const toggleItemName = "toggle"

var Mux = &http.ServeMux{}

func init() {
	RegisterParserByFunc("toggle", func(data []byte) (interface{}, error) {
		var rawToggle rawToggle
		err := json.Unmarshal(data, &rawToggle)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal toggle: %s", err.Error())
		}
		return createToggle(&rawToggle)
	})
	RegisterParserByFunc("json", func(data []byte) (interface{}, error) {
		var obj interface{}
		err := json.Unmarshal(data, &obj)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %s", err.Error())
		}
		return obj, nil
	})
}

func GetObject(namespace string, objectName string, targetKV ...string) interface{} {
	objectId := objectId{namespace, objectName}
	toggle, _ := getItem(objectId, toggleItemName).(toggle)
	if toggle == nil {
		countlog.Error("event!counselor.object not found",
			"namespace", namespace, "objectName", objectName)
		return nil
	}
	target := map[string]string{}
	for i := 0; i < len(targetKV); i+=2 {
		target[targetKV[i]] = targetKV[i+1]
	}
	variant, err := toggle(target)
	if err != nil {
		countlog.Error("event!counselor.failed to execute toggle", "err", err)
		return nil
	}
	return getItem(objectId, variant)
}

func ShouldUse(namespace string, objectName string, targetKV ...string) bool {
	verdict, _ := GetObject(namespace, objectName, targetKV...).(bool)
	return verdict
}
