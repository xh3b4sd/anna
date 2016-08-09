package main

import (
	"github.com/xh3b4sd/anna/api"
	"github.com/xh3b4sd/anna/server/interface/text"
	"github.com/xh3b4sd/anna/spec"
)

func createTextInterface(newLog spec.Log, newTextInput chan api.TextRequest, newTextOutput chan api.TextResponse) (spec.TextInterface, error) {
	newTextInterfaceConfig := text.DefaultInterfaceConfig()
	newTextInterfaceConfig.Log = newLog
	newTextInterfaceConfig.TextInput = newTextInput
	newTextInterfaceConfig.TextOutput = newTextOutput
	newTextInterface, err := text.NewInterface(newTextInterfaceConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	return newTextInterface, nil
}
