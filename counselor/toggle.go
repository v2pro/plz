package counselor

import "fmt"

type rawToggle struct {
	PreprocessorFn   string
	PreprocessorArgs map[string]interface{}
	Variants         []*variant
	DefaultVariant   itemName
}

type variant struct {
	ItemName itemName
	RuleFn   string
	RuleArgs map[string]interface{}
}

type toggle func(target map[string]string) (itemName, error)

var createToggle = func(rawToggle *rawToggle) (toggle, error) {
	var err error
	var rules = make([]rule, len(rawToggle.Variants))
	variants := []string{}
	for i, variant := range rawToggle.Variants {
		variants = append(variants, string(variant.ItemName))
		if variant.RuleFn == "" {
			return nil, fmt.Errorf("missing rule for variant: %s", variant.ItemName)
		}
		rules[i], err = createRule(variant.RuleFn, variant.RuleArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule for variant %s: %s", variant.ItemName, err.Error())
		}
	}
	variants = append(variants, string(rawToggle.DefaultVariant))
	var preprocessor preprocessor
	if rawToggle.PreprocessorFn != "" {
		preprocessor, err = createPreprocessor(rawToggle.PreprocessorFn, rawToggle.PreprocessorArgs, variants)
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
		return rawToggle.DefaultVariant, nil
	}, nil
}
