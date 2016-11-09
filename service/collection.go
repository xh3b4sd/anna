package service

import (
	"sync"

	"github.com/xh3b4sd/anna/service/fs/mem"
	"github.com/xh3b4sd/anna/service/id"
	"github.com/xh3b4sd/anna/service/log"
	"github.com/xh3b4sd/anna/service/permutation"
	"github.com/xh3b4sd/anna/service/random"
	servicespec "github.com/xh3b4sd/anna/service/spec"
	"github.com/xh3b4sd/anna/service/textinput"
	"github.com/xh3b4sd/anna/service/textoutput"
)

// CollectionConfig represents the configuration used to create a new service
// collection object.
type CollectionConfig struct {
	// Dependencies.
	FSService          servicespec.FS
	IDService          servicespec.ID
	LogService         servicespec.Log
	PermutationService servicespec.Permutation
	RandomService      servicespec.Random
	TextInputService   servicespec.TextInput
	TextOutputService  servicespec.TextOutput
}

// DefaultCollectionConfig provides a default configuration to create a new
// service collection object by best effort.
func DefaultCollectionConfig() CollectionConfig {
	newConfig := CollectionConfig{
		// Dependencies.
		FSService:          mem.MustNew(),
		IDService:          id.MustNew(),
		LogService:         log.MustNew(),
		PermutationService: permutation.MustNew(),
		RandomService:      random.MustNew(),
		TextInputService:   textinput.MustNew(),
		TextOutputService:  textoutput.MustNew(),
	}

	return newConfig
}

// NewCollection creates a new configured service collection object.
func NewCollection(config CollectionConfig) (servicespec.Collection, error) {
	newCollection := &collection{
		CollectionConfig: config,

		ShutdownOnce: sync.Once{},
	}

	if newCollection.FSService == nil {
		return nil, maskAnyf(invalidConfigError, "file system service must not be empty")
	}
	if newCollection.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if newCollection.LogService == nil {
		return nil, maskAnyf(invalidConfigError, "log service must not be empty")
	}
	if newCollection.PermutationService == nil {
		return nil, maskAnyf(invalidConfigError, "permutation service must not be empty")
	}
	if newCollection.RandomService == nil {
		return nil, maskAnyf(invalidConfigError, "random service must not be empty")
	}
	if newCollection.TextInputService == nil {
		return nil, maskAnyf(invalidConfigError, "text input service must not be empty")
	}
	if newCollection.TextOutputService == nil {
		return nil, maskAnyf(invalidConfigError, "text output service must not be empty")
	}

	return newCollection, nil
}

// MustNewCollection creates either a new default configured service collection,
// or panics.
func MustNewCollection() servicespec.Collection {
	newCollection, err := NewCollection(DefaultCollectionConfig())
	if err != nil {
		panic(err)
	}

	return newCollection
}

type collection struct {
	CollectionConfig

	ShutdownOnce sync.Once
}

func (c *collection) FS() servicespec.FS {
	return c.FSService
}

func (c *collection) ID() servicespec.ID {
	return c.IDService
}

func (c *collection) Log() servicespec.Log {
	return c.LogService
}

func (c *collection) Permutation() servicespec.Permutation {
	return c.PermutationService
}

func (c *collection) Random() servicespec.Random {
	return c.RandomService
}

func (c *collection) Shutdown() {
	c.ShutdownOnce.Do(func() {
		var wg sync.WaitGroup

		//wg.Add(1)
		//go func() {
		//	c.TODO().Shutdown()
		//	wg.Done()
		//}()

		wg.Wait()
	})
}

func (c *collection) TextInput() servicespec.TextInput {
	return c.TextInputService
}

func (c *collection) TextOutput() servicespec.TextOutput {
	return c.TextOutputService
}
