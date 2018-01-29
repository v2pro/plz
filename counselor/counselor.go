package counselor

import (
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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
	for i := 0; i < len(targetKV); i += 2 {
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

// SetObjectVariant is only intended to be used in test
func SetObjectVariant(namespace string, objectName string, variantName string, dataFormat string, content []byte) {
	dir := path.Join(SourceDir, namespace, objectName)
	os.MkdirAll(dir, 0777)
	err := ioutil.WriteFile(path.Join(dir, variantName), append([]byte("data_format:"+dataFormat+"\n"), content...), 0666)
	if err != nil {
		panic(err)
	}
}

// SetObjectToggle is only intended to be used in test
func SetObjectToggle(namespace string, objectName string, content string) {
	dir := path.Join(SourceDir, namespace, objectName)
	os.MkdirAll(dir, 0777)
	err := ioutil.WriteFile(path.Join(dir, toggleItemName), []byte("data_format:toggle\n"+content), 0666)
	if err != nil {
		panic(err)
	}
}

// SetObject is only intended to be used in test
func SetObject(namespace string, objectName string, dataFormat string, content []byte) {
	SetObjectToggle(namespace, objectName, `{"DefaultVariant":"default"}`)
	SetObjectVariant(namespace, objectName, "default", dataFormat, content)
}
