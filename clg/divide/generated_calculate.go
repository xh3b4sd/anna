package divide

// This file is generated by the CLG generator. Don't edit it manually. The CLG
// generator is invoked by go generate. For more information about the usage of
// the CLG generator check https://github.com/xh3b4sd/clggen or have a look at
// the clg package. There is the go generate statement placed to invoke clggen.

import (
	"reflect"

	"github.com/xh3b4sd/anna/factory/id"
	"github.com/xh3b4sd/anna/log"
	"github.com/xh3b4sd/anna/spec"
	"github.com/xh3b4sd/anna/storage/memory"
)

const (
	// ObjectType represents the object type of the CLG object. This is used e.g.
	// to register itself to the logger.
	ObjectType spec.ObjectType = "divide-clg"
)

// Config represents the configuration used to create a new CLG object.
type Config struct {
	// Dependencies.
	IDFactory spec.IDFactory
	Log       spec.Log
	Storage   spec.Storage

	// Settings.
	InputChannel chan spec.NetworkPayload
}

// DefaultConfig provides a default configuration to create a new CLG object by
// best effort.
func DefaultConfig() Config {
	newConfig := Config{
		// Dependencies.
		IDFactory: id.MustNewFactory(),
		Log:       log.New(log.DefaultConfig()),
		Storage:   memory.MustNew(),

		// Settings.
		InputChannel: make(chan spec.NetworkPayload, 1000),
	}

	return newConfig
}

// New creates a new configured CLG object.
func New(config Config) (spec.CLG, error) {
	newCLG := &clg{
		Config: config,
		ID:     id.MustNew(),
		Type:   ObjectType,
	}

	// Dependencies.
	if newCLG.Log == nil {
		return nil, maskAnyf(invalidConfigError, "logger must not be empty")
	}
	if newCLG.Storage == nil {
		return nil, maskAnyf(invalidConfigError, "storage must not be empty")
	}

	// Settings.
	if newCLG.InputChannel == nil {
		return nil, maskAnyf(invalidConfigError, "input channel must not be empty")
	}

	newCLG.Log.Register(newCLG.GetType())

	return newCLG, nil
}

// MustNew creates either a new default configured CLG object, or panics.
func MustNew() spec.CLG {
	newCLG, err := New(DefaultConfig())
	if err != nil {
		panic(err)
	}

	return newCLG
}

type clg struct {
	Config

	ID   spec.ObjectID
	Type spec.ObjectType
}

func (c *clg) Calculate(payload spec.NetworkPayload) (spec.NetworkPayload, error) {
	outputs := reflect.ValueOf(c.calculate).Call(payload.GetArgs())

	modifiedPayload, err := injectValues(payload, outputs)
	if err != nil {
		return nil, maskAny(err)
	}

	return modifiedPayload, nil
}

func (c *clg) GetName() string {
	return "divide"
}

func (c *clg) GetInputChannel() chan spec.NetworkPayload {
	return c.InputChannel
}

func (c *clg) GetInputTypes() []reflect.Type {
	t := reflect.TypeOf(c.calculate)

	var inputType []reflect.Type

	for i := 0; i < t.NumIn(); i++ {
		inputType = append(inputType, t.In(i))
	}

	return inputType
}

func (c *clg) SetIDFactory(idFactory spec.IDFactory) {
	c.IDFactory = idFactory
}

func (c *clg) SetLog(log spec.Log) {
	c.Log = log
}

func (c *clg) SetStorage(storage spec.Storage) {
	c.Storage = storage
}
