package counselor

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/v2pro/plz/countlog"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

const paramCreateRule = "create_rule"
const keyDivideBucketsBy = "divide_buckets_by"

type rule func(target map[string]string) (bool, error)
type ruleFactory func(args map[string]interface{}) (rule, error)

var ruleFactories = map[string]ruleFactory{
	"thousand_buckets": thousandBucketsRule,
	"eq":               eqRule,
	"and":              andRule,
	"or":               orRule,
}
var ruleFactoriesMutex = &sync.Mutex{}

func createRule(fn string, args map[string]interface{}) (rule, error) {
	args[paramCreateRule] = _createRule
	return _createRule(fn, args)
}

func _createRule(fn string, args map[string]interface{}) (
	func(target map[string]string) (bool, error), error) {
	factory, err := getRuleFactory(fn)
	if err != nil {
		return nil, err
	}
	rule, err := factory(args)
	if err != nil {
		countlog.Error("event!counselor.failed to create rule", "fn", fn, "err", err)
		return nil, err
	}
	return rule, nil
}

func getRuleFactory(fn string) (ruleFactory, error) {
	ruleFactoriesMutex.Lock()
	factory := ruleFactories[fn]
	ruleFactoriesMutex.Unlock()
	if factory != nil {
		return factory, nil
	}
	sym, err := loadFn(fn)
	if err != nil {
		countlog.Error("event!counselor.failed to load fn", "fn", fn, "err", err)
		return nil, err
	}
	untypedFactory, _ := sym.(func(args map[string]interface{}) (func(target map[string]string) (bool, error), error))
	if untypedFactory == nil {
		countlog.Error("event!counselor.fn is not valid rule factory", "fn", fn, "err", err)
		return nil, err
	}
	factory = func(args map[string]interface{}) (rule, error) {
		return untypedFactory(args)
	}
	ruleFactoriesMutex.Lock()
	ruleFactories[fn] = factory
	ruleFactoriesMutex.Unlock()
	return factory, nil
}

func thousandBucketsRule(args map[string]interface{}) (rule, error) {
	salt, _ := args["salt"].(string)
	if salt == "" {
		salt = "SALT"
	}
	segmentStartFloat, ok := args["segment_start"].(float64)
	if !ok {
		return nil, errors.New("missing segment start")
	}
	segmentStart := int(segmentStartFloat)
	segmentEndFloat, ok := args["segment_end"].(float64)
	if !ok {
		return nil, errors.New("missing segment end")
	}
	segmentEnd := int(segmentEndFloat)
	return func(target map[string]string) (bool, error) {
		value, found := target[keyDivideBucketsBy]
		if !found {
			value = strconv.Itoa(int(rand.Int31n(math.MaxInt32)))
		}
		index := saltHash(value, salt, 1000)
		return index >= segmentStart && index < segmentEnd, nil
	}, nil
}

func saltHash(uuid string, salt string, bucketSize int) int {
	h := sha1.New()
	h.Write([]byte(uuid + salt))
	bytes := h.Sum(nil)
	size := len(bytes)
	value := binary.BigEndian.Uint32(bytes[size-4 : size])
	// avoid negative number mod
	mod := int64(value) % int64(bucketSize)
	return int(mod)
}

func eqRule(args map[string]interface{}) (rule, error) {
	targetKey, _ := args["target_key"].(string)
	if targetKey == "" {
		return nil, errors.New("missing key")
	}
	operand, _ := args["operand"].([]string)
	if operand == nil {
		return nil, errors.New("missing operand")
	}
	m := map[string]struct{}{}
	for _, elem := range operand {
		m[elem] = struct{}{}
	}
	return func(target map[string]string) (bool, error) {
		value, found := target[targetKey]
		if !found {
			return false, fmt.Errorf("key %s not found in target", targetKey)
		}
		_, found = m[value]
		return found, nil
	}, nil
}

func andRule(args map[string]interface{}) (rule, error) {
	ruleDefs, _ := args["rules"].([]map[string]interface{})
	if ruleDefs == nil {
		return nil, errors.New("missing rules")
	}
	createRule, _ := args[paramCreateRule].(func(fn string, args map[string]interface{}) (
		func(target map[string]string) (bool, error), error))
	if createRule == nil {
		return nil, errors.New("missing " + paramCreateRule)
	}
	rules := []func(target map[string]string) (bool, error){}
	for _, ruleDef := range ruleDefs {
		fn, _ := ruleDef["fn"].(string)
		if fn == "" {
			return nil, errors.New("missing fn in rule")
		}
		args, _ := ruleDef["args"].(map[string]interface{})
		if args == nil {
			return nil, errors.New("missing args in rule")
		}
		rule, err := createRule(fn, args)
		if err != nil {
			return nil, fmt.Errorf("create fn %s failed: %s", fn, err.Error())
		}
		rules = append(rules, rule)
	}
	return func(target map[string]string) (bool, error) {
		if len(rules) == 0 {
			return false, nil
		}
		finalVerdict := true
		for i, rule := range rules {
			verdict, err := rule(target)
			if err != nil {
				return false, fmt.Errorf("%s failed to execute rule: %s", ruleDefs[i]["fn"], err.Error())
			}
			finalVerdict = finalVerdict && verdict
		}
		return finalVerdict, nil
	}, nil
}

func orRule(args map[string]interface{}) (rule, error) {
	ruleDefs, _ := args["rules"].([]map[string]interface{})
	if ruleDefs == nil {
		return nil, errors.New("missing rules")
	}
	createRule, _ := args[paramCreateRule].(func(fn string, args map[string]interface{}) (
		func(target map[string]string) (bool, error), error))
	if createRule == nil {
		return nil, errors.New("missing " + paramCreateRule)
	}
	rules := []func(target map[string]string) (bool, error){}
	for _, ruleDef := range ruleDefs {
		fn, _ := ruleDef["fn"].(string)
		if fn == "" {
			return nil, errors.New("missing fn in rule")
		}
		args, _ := ruleDef["args"].(map[string]interface{})
		if args == nil {
			return nil, errors.New("missing args in rule")
		}
		rule, err := createRule(fn, args)
		if err != nil {
			return nil, fmt.Errorf("create fn %s failed: %s", fn, err.Error())
		}
		rules = append(rules, rule)
	}
	return func(target map[string]string) (bool, error) {
		if len(rules) == 0 {
			return true, nil
		}
		finalVerdict := true
		for i, rule := range rules {
			verdict, err := rule(target)
			if err != nil {
				return false, fmt.Errorf("%s failed to execute rule: %s", ruleDefs[i]["fn"], err.Error())
			}
			finalVerdict = finalVerdict || verdict
		}
		return finalVerdict, nil
	}, nil
}
