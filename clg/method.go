package clg

import (
	"reflect"
	"strings"
)

func (i *clgIndex) CallCLGByName(args ...interface{}) ([]interface{}, error) {
	methodName, err := ArgToString(args, 0)
	if err != nil {
		return nil, maskAny(err)
	}

	inputValues := ArgsToValues(args[1:])
	methodValue := reflect.ValueOf(i).MethodByName(methodName)
	if !methodValue.IsValid() {
		return nil, maskAnyf(methodNotFoundError, methodName)
	}

	outputValues := methodValue.Call(inputValues)
	results, err := ValuesToArgs(outputValues)
	if err != nil {
		return nil, maskAny(err)
	}

	return results, nil
}

func (i *clgIndex) GetCLGNames(args ...interface{}) ([]interface{}, error) {
	if len(args) > 1 {
		return nil, maskAnyf(tooManyArgumentsError, "expected 1 got %d", len(args))
	}
	var pattern string
	if len(args) == 1 {
		var err error
		pattern, err = ArgToString(args, 0)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var allCLGNames []string

	t := reflect.TypeOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if pattern != "" && !strings.Contains(methodName, pattern) {
			continue
		}
		allCLGNames = append(allCLGNames, methodName)
	}

	return []interface{}{allCLGNames}, nil
}
