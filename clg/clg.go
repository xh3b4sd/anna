// Package clg implementes fundamental actions used to create strategies that
// allow to discover new behavior for problem solving.
package clg

import (
	"sync"

	"github.com/xh3b4sd/anna/id"
	"github.com/xh3b4sd/anna/log"
	"github.com/xh3b4sd/anna/spec"
)

const (
	// ObjectTypeCLGIndex represents the object type of the CLG index object.
	// This is used e.g. to register itself to the logger.
	ObjectTypeCLGIndex spec.ObjectType = "clg-index"
)

// Config represents the configuration used to create a new CLG index object.
type Config struct {
	// Dependencies.
	Log spec.Log
}

// DefaultConfig provides a default configuration to create a new CLG index
// object by best effort.
func DefaultConfig() Config {
	newConfig := Config{
		// Dependencies.
		Log: log.NewLog(log.DefaultConfig()),
	}

	return newConfig
}

// NewCLGIndex creates a new configured CLG index object.
func NewCLGIndex(config Config) (spec.CLGIndex, error) {
	newCLGIndex := &clgIndex{
		Config: config,
		ID:     id.NewObjectID(id.Hex128),
		Mutex:  sync.Mutex{},
		Type:   ObjectTypeCLGIndex,
	}

	newCLGIndex.Log.Register(newCLGIndex.GetType())

	return newCLGIndex, nil
}

type clgIndex struct {
	Config

	ID    spec.ObjectID
	Mutex sync.Mutex
	Type  spec.ObjectType
}
