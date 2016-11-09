package splitfeatures

// This file is generated by the CLG generator. Don't edit it manually. The CLG
// generator is invoked by go generate. For more information about the usage of
// the CLG generator check https://github.com/xh3b4sd/clggen or have a look at
// the clg package. There is the go generate statement placed to invoke clggen.

import (
	"github.com/xh3b4sd/anna/service"
	servicespec "github.com/xh3b4sd/anna/service/spec"
	"github.com/xh3b4sd/anna/storage"
	storagespec "github.com/xh3b4sd/anna/storage/spec"
)

// Config represents the configuration used to create a new CLG object.
type Config struct {
	// Dependencies.
	ServiceCollection servicespec.Collection
	StorageCollection storagespec.Collection
}

// DefaultConfig provides a default configuration to create a new CLG object by
// best effort.
func DefaultConfig() Config {
	newConfig := Config{
		// Dependencies.
		ServiceCollection: service.MustNewCollection(),
		StorageCollection: storage.MustNewCollection(),
	}

	return newConfig
}

// New creates a new configured CLG object.
func New(config Config) (servicespec.CLG, error) {
	newService := &clg{
		Config: config,
	}

	// Dependencies.
	if newService.ServiceCollection == nil {
		return nil, maskAnyf(invalidConfigError, "factory collection must not be empty")
	}
	if newService.StorageCollection == nil {
		return nil, maskAnyf(invalidConfigError, "storage collection must not be empty")
	}

	id, err := newService.Service().ID().New()
	if err != nil {
		return nil, maskAny(err)
	}
	newService.Metadata["id"] = id
	newService.Metadata["kind"] = "splitfeatures"
	newService.Metadata["name"] = "clg"
	newService.Metadata["type"] = "service"

	return newService, nil
}

// MustNew creates either a new default configured CLG object, or panics.
func MustNew() servicespec.CLG {
	newService, err := New(DefaultConfig())
	if err != nil {
		panic(err)
	}

	return newService
}

type clg struct {
	Config

	Metadata map[string]string
}

func (c *clg) Service() servicespec.Collection {
	return c.ServiceCollection
}

func (c *clg) GetCalculate() interface{} {
	return c.calculate
}

func (c *clg) SetServiceCollection(serviceCollection servicespec.Collection) {
	c.ServiceCollection = serviceCollection
}

func (c *clg) SetStorageCollection(storageCollection storagespec.Collection) {
	c.StorageCollection = storageCollection
}

func (c *clg) Storage() storagespec.Collection {
	return c.StorageCollection
}
