package counselor

import "fmt"

type rawToggle struct {
	preprocessorFn   string
	preprocessorArgs map[string]interface{}
	variants         []*variant
	defaultVariant   itemName
}

type variant struct {
	itemName itemName
	ruleFn   string
	ruleArgs map[string]interface{}
}

type toggle func(target map[string]string) (itemName, error)

var createToggle = func(rawToggle *rawToggle) (toggle, error) {
	var err error
	var rules = [len(rawToggle.variants)]rule{}
	variants := []string{}
	for i, variant := range rawToggle.variants {
		variants = append(variants, string(variant.itemName))
		if variant.ruleFn == "" {
			return nil, fmt.Errorf("missing rule for variant: %s", variant.itemName)
		}
		rules[i], err = createRule(variant.ruleFn, variant.ruleArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule for variant %s: %s", variant.itemName, err.Error())
		}
	}
	variants = append(variants, string(rawToggle.defaultVariant))
	var preprocessor preprocessor
	if rawToggle.preprocessorFn != "" {
		preprocessor, err = createPreprocessor(rawToggle.preprocessorFn, rawToggle.preprocessorArgs, variants)
		if err != nil {
			return nil, err
		}
	}
	return func(target map[string]string) (itemName, error) {
		if preprocessor != nil {
			err := preprocessor(target)
			if err != nil {
				return "", fmt.Errorf("failed to preprocess: %s", err.Error())
			}
		}
		for i, rule := range rules {
			verdict, err := rule(target)
			if err != nil {
				return "", fmt.Errorf("failed to run rule for %s: %s", variants[i], err.Error())
			}
			if verdict {
				return itemName(variants[i]), nil
			}
		}
		return rawToggle.defaultVariant, nil
	}, nil
}
