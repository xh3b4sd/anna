package spec

import (
	"reflect"
)

// CLG represents the CLGs interacting with each other within the neural
// network. Each CLG is registered in the Network. From there signal are
// dispatched in a dynamic fashion until some useful calculation took place.
type CLG interface {
	// Calculate provides the CLG's actual business logic.
	Calculate(payload NetworkPayload) (NetworkPayload, error)

	// GetName returns the CLG's human readable name.
	GetName() string

	// GetInputChannel returns the CLG's input channel, which acts as
	// communication channel to reach the CLG inside of the neural network.
	GetInputChannel() chan NetworkPayload

	// GetInputTypes returns the CLG's underlying input types. These reflect the
	// real interface hidden behind the Calculate API.
	GetInputTypes() []reflect.Type

	Object

	// SetFeatureStorage configures the CLG's feature storage. This is done for
	// all CLGs, regardless if a CLG is making use of the logger or not.
	SetFeatureStorage(storage Storage)

	// SetGeneralStorage configures the CLG's feature storage. This is done for
	// all CLGs, regardless if a CLG is making use of the logger or not.
	SetGeneralStorage(storage Storage)

	// SetIDFactory configures the CLG's id factory. This is done for all CLGs,
	// regardless if a CLG is making use of the logger or not.
	SetIDFactory(idFactory IDFactory)

	// SetLog configures the CLG's logger. This is done for all CLGs, regardless
	// if a CLG is making use of the logger or not.
	SetLog(log Log)
}
