// Package network implements spec.Network to provide a neural network based on
// dynamic and self improving CLG execution. The network provides input and
// output channels. When input is received it is injected into the neural
// communication. The following neural activity calculates output which is
// streamed through the output channel back to the requestor.
package network

import (
	"reflect"
	"sync"
	"time"

	"github.com/xh3b4sd/anna/api"
	"github.com/xh3b4sd/anna/context"
	"github.com/xh3b4sd/anna/factory/id"
	"github.com/xh3b4sd/anna/factory/permutation"
	"github.com/xh3b4sd/anna/log"
	"github.com/xh3b4sd/anna/spec"
	"github.com/xh3b4sd/anna/storage/memory"
)

const (
	// ObjectType represents the object type of the network object. This is used
	// e.g. to register itself to the logger.
	ObjectType spec.ObjectType = "network"
)

// Config represents the configuration used to create a new network object.
type Config struct {
	// Dependencies.
	IDFactory          spec.IDFactory
	Log                spec.Log
	PermutationFactory spec.PermutationFactory
	Storage            spec.Storage
	TextInput          chan spec.TextRequest
	TextOutput         chan spec.TextResponse

	// Settings.

	// Delay causes each CLG execution to be delayed. This value represents a
	// default value. The actually used value is optimized based on experience
	// and learning.
	// TODO implement
	Delay time.Duration
}

// DefaultConfig provides a default configuration to create a new network
// object by best effort.
func DefaultConfig() Config {
	newPermutationFactory, err := permutation.NewFactory(permutation.DefaultFactoryConfig())
	if err != nil {
		panic(err)
	}

	newStorage, err := memory.NewStorage(memory.DefaultStorageConfig())
	if err != nil {
		panic(err)
	}

	newConfig := Config{
		// Dependencies.
		IDFactory:          id.MustNewFactory(),
		Log:                log.NewLog(log.DefaultConfig()),
		PermutationFactory: newPermutationFactory,
		Storage:            newStorage,
		TextInput:          make(chan spec.TextRequest, 1000),
		TextOutput:         make(chan spec.TextResponse, 1000),

		// Settings.
		Delay: 0,
	}

	return newConfig
}

// New creates a new configured network object.
func New(config Config) (spec.Network, error) {
	newNetwork := &network{
		Config: config,

		BootOnce:     sync.Once{},
		CLGIDs:       map[string]spec.ObjectID{},
		CLGs:         newCLGs(),
		ID:           id.MustNew(),
		Mutex:        sync.Mutex{},
		ShutdownOnce: sync.Once{},
		Type:         ObjectType,
	}

	if newNetwork.IDFactory == nil {
		return nil, maskAnyf(invalidConfigError, "ID factory must not be empty")
	}
	if newNetwork.Log == nil {
		return nil, maskAnyf(invalidConfigError, "logger must not be empty")
	}
	if newNetwork.PermutationFactory == nil {
		return nil, maskAnyf(invalidConfigError, "permutation factory must not be empty")
	}
	if newNetwork.Storage == nil {
		return nil, maskAnyf(invalidConfigError, "storage must not be empty")
	}
	if newNetwork.TextInput == nil {
		return nil, maskAnyf(invalidConfigError, "text input channel must not be empty")
	}
	if newNetwork.TextOutput == nil {
		return nil, maskAnyf(invalidConfigError, "text output channel must not be empty")
	}

	newNetwork.Log.Register(newNetwork.GetType())

	return newNetwork, nil
}

type network struct {
	Config

	BootOnce sync.Once

	// CLGIDs provides a mapping of CLG names pointing to their corresponding CLG
	// ID.
	CLGIDs map[string]spec.ObjectID

	CLGs         map[spec.ObjectID]spec.CLG // TODO there is probably no reason to index the CLGs like this
	ID           spec.ObjectID
	Mutex        sync.Mutex
	ShutdownOnce sync.Once
	Type         spec.ObjectType
}

func (n *network) Activate(CLG spec.CLG, payload spec.NetworkPayload, queue []spec.NetworkPayload) (spec.NetworkPayload, []spec.NetworkPayload, error) {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Activate")

	queue = append(queue, payload)

	// Prepare the permutation list to find out which combination of payloads
	// satisfies the requested CLG's interface.
	newConfig := permutation.DefaultListConfig()
	newConfig.MaxGrowth = len(CLG.GetInputTypes())
	newConfig.Values = queueToValues(queue)
	newPermutationList, err := permutation.NewList(newConfig)
	if err != nil {
		return nil, nil, maskAny(err)
	}

	for {
		err := n.PermutationFactory.MapTo(newPermutationList)
		if err != nil {
			return nil, nil, maskAny(err)
		}

		// Check if the given payload satisfies the requested CLG's interface.
		members := newPermutationList.GetMembers()
		types, err := membersToTypes(members)
		if err != nil {
			return nil, nil, maskAny(err)
		}
		if reflect.DeepEqual(types, CLG.GetInputTypes()) {
			newPayload, err := membersToPayload(members)
			if err != nil {
				return nil, nil, maskAny(err)
			}
			newQueue, err := filterMembersFromQueue(members, queue)
			if err != nil {
				return nil, nil, maskAny(err)
			}

			// In case the current queue exeeds the interface of the requested CLG, it is
			// trimmed to cause a more strict behaviour of the neural network.
			if len(newPermutationList.GetValues()) > len(CLG.GetInputTypes()) {
				newQueue = newQueue[1:]
			}

			return newPayload, newQueue, nil
		}

		err = n.PermutationFactory.PermuteBy(newPermutationList, 1)
		if permutation.IsMaxGrowthReached(err) {
			// We cannot permute the given list anymore. So far the requested CLG's
			// interface could not be satisfied.
			return nil, nil, maskAnyf(invalidInterfaceError, "types must match")
		}
	}
}

func (n *network) Boot() {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Boot")

	n.BootOnce.Do(func() {
		n.CLGs = n.configureCLGs(n.CLGs)
		n.CLGIDs = n.mapCLGIDs(n.CLGs)

		go n.Listen()
		go n.listenCLGs()
	})
}

func (n *network) Calculate(CLG spec.CLG, payload spec.NetworkPayload) (spec.NetworkPayload, error) {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Calculate")

	calculatedPayload, err := CLG.Calculate(payload)
	if err != nil {
		return nil, maskAny(err)
	}

	return calculatedPayload, nil
}

// TODO
func (n *network) Forward(CLG spec.CLG, payload spec.NetworkPayload) error {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Forward")

	// Check if the given spec.Context provides a CLG tree ID.
	ctx, err := payload.GetContext()
	if err != nil {
		return maskAny(err)
	}
	clgTreeID := ctx.GetCLGTreeID()

	if clgTreeID == "" {
		// create new
	} else {
		// lookup existing
	}

	return nil
}

func (n *network) Listen() {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Listen")

	// Listen on TextInput from the outside to receive text requests.
	CLG, err := n.clgByName("input")
	if err != nil {
		n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
	}

	clgID := CLG.GetID()
	networkID := n.GetID()
	clgChannel := CLG.GetInputChannel()

	for {
		textRequest := <-n.TextInput

		// This should only be used for testing to bypass the neural network
		// and directly respond with the received input.
		if textRequest.GetEcho() {
			newTextResponseConfig := api.DefaultTextResponseConfig()
			newTextResponseConfig.Output = textRequest.GetInput()
			newTextResponse, err := api.NewTextResponse(newTextResponseConfig)
			if err != nil {
				n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
			}
			n.TextOutput <- newTextResponse
			continue
		}

		ctxConfig := context.DefaultConfig()
		ctxConfig.SessionID = textRequest.GetSessionID()
		ctx, err := context.New(ctxConfig)
		if err != nil {
			n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
			continue
		}

		payloadConfig := api.DefaultNetworkPayloadConfig()
		payloadConfig.Args = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(textRequest.GetInput())}
		payloadConfig.Destination = clgID
		payloadConfig.Sources = []spec.ObjectID{networkID}
		newPayload, err := api.NewNetworkPayload(payloadConfig)
		if err != nil {
			n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
			continue
		}

		clgChannel <- newPayload
	}
}

func (n *network) Shutdown() {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Shutdown")

	n.ShutdownOnce.Do(func() {
		// TODO shutdown CLG listeners
		// TODO stop listening on text requests
	})
}
