package counselor

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"sync"
)

const paramVariants = "variants"
const paramCreatePreprocessor = "create_preprocessor"

type preprocessor func(target map[string]string) error
type preprocessorFactory func(args map[string]interface{}) (preprocessor, error)

var preprocessorFactories = map[string]preprocessorFactory{
	"md5":                md5Preprocessor,
	"chain_preprocessor": chainPreprocessor,
}
var preprocessorFactoriesMutex = &sync.Mutex{}

func createPreprocessor(fn string, args map[string]interface{}, variants []string) (preprocessor, error) {
	args[paramVariants] = variants
	args[paramCreatePreprocessor] = _createPreprocessor
	return _createPreprocessor(fn, args)
}

func _createPreprocessor(fn string, args map[string]interface{}) (
	func(target map[string]string) error, error) {
	factory, err := getPreprocessorFactory(fn)
	if err != nil {
		return nil, err
	}
	preprocessor, err := factory(args)
	if err != nil {
		countlog.Error("event!counselor.failed to create preprocessor", "fn", fn, "err", err)
		return nil, err
	}
	return preprocessor, nil
}

func getPreprocessorFactory(fn string) (preprocessorFactory, error) {
	preprocessorFactoriesMutex.Lock()
	factory := preprocessorFactories[fn]
	preprocessorFactoriesMutex.Unlock()
	if factory != nil {
		return factory, nil
	}
	sym, err := loadFn(fn)
	if err != nil {
		countlog.Error("event!counselor.failed to load fn", "fn", fn, "err", err)
		return nil, err
	}
	untypedFactory, _ := sym.(func(args map[string]interface{}) (func(target map[string]string) error, error))
	if untypedFactory == nil {
		countlog.Error("event!counselor.fn is not valid preprocessor factory", "fn", fn, "err", err)
		return nil, err
	}
	factory = func(args map[string]interface{}) (preprocessor, error) {
		return untypedFactory(args)
	}
	preprocessorFactoriesMutex.Lock()
	preprocessorFactories[fn] = factory
	preprocessorFactoriesMutex.Unlock()
	return factory, nil
}

func md5Preprocessor(args map[string]interface{}) (preprocessor, error) {
	targetKey, _ := args["target_key"].(string)
	if targetKey == "" {
		return nil, errors.New("missing target_key")
	}
	return func(target map[string]string) error {
		value, found := target[targetKey]
		if !found {
			return fmt.Errorf("key %s not found in target", targetKey)
		}
		bytes := md5.Sum([]byte(value))
		value = fmt.Sprintf("%x", bytes)
		target[targetKey] = value
		return nil
	}, nil
}

func chainPreprocessor(args map[string]interface{}) (preprocessor, error) {
	preprocessorDefs, _ := args["preprocessors"].([]map[string]interface{})
	if preprocessorDefs == nil {
		return nil, errors.New("missing preprocessors")
	}
	createPreprocessor, _ := args[paramCreatePreprocessor].(func(fn string, args map[string]interface{}) (
		func(target map[string]string) error, error))
	if createPreprocessor == nil {
		return nil, errors.New("missing " + paramCreatePreprocessor)
	}
	preprocessors := []func(target map[string]string) error{}
	for _, preprocessorDef := range preprocessorDefs {
		fn, _ := preprocessorDef["fn"].(string)
		if fn == "" {
			return nil, errors.New("missing fn in preprocessor")
		}
		args, _ := preprocessorDef["args"].(map[string]interface{})
		if args == nil {
			return nil, errors.New("missing args in preprocessor")
		}
		preprocessor, err := createPreprocessor(fn, args)
		if err != nil {
			return nil, fmt.Errorf("create fn %s failed: %s", fn, err.Error())
		}
		preprocessors = append(preprocessors, preprocessor)
	}
	return func(target map[string]string) error {
		for i, preprocessor := range preprocessors {
			err := preprocessor(target)
			if err != nil {
				return fmt.Errorf("%s failed to preprocess: %s", preprocessorDefs[i]["fn"], err.Error())
			}
		}
		return nil
	}, nil
}
