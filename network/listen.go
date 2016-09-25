package network

import (
	"reflect"

	"github.com/xh3b4sd/anna/api"
	"github.com/xh3b4sd/anna/clg/output"
	"github.com/xh3b4sd/anna/context"
	"github.com/xh3b4sd/anna/spec"
)

// receiver

func (n *network) listenCLGs() {
	// Make all CLGs listening in their specific input channel.
	for ID, CLG := range n.CLGs {
		go func(ID spec.ObjectID, CLG spec.CLG) {
			var queue []spec.NetworkPayload
			queueBuffer := len(CLG.GetInputTypes()) + 1
			inputChannel := CLG.GetInputChannel()

			for {
				select {
				case <-n.Closer:
					break
				case payload := <-inputChannel:
					// In case the current queue exeeds a certain amount of payloads, it
					// is unlikely that the queue is going to be helpful when growing any
					// further. Thus we cut the queue at some point beyond the interface
					// capabilities of the requested CLG.
					queue = append(queue, payload)
					if len(queue) > queueBuffer {
						queue = queue[1:]
					}

					go func(payload spec.NetworkPayload) {
						// Activate if the CLG's interface is satisfied by the given
						// network payload.
						newPayload, newQueue, err := n.Activate(CLG, queue)
						if IsInvalidInterface(err) {
							// The interface of the requested CLG was not fulfilled. We
							// continue listening for the next network payload without doing
							// any work.
							return
						} else if err != nil {
							n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
						}
						queue = newQueue

						// Calculate based on the CLG's implemented business logic.
						calculatedPayload, err := n.Calculate(CLG, newPayload)
						if output.IsExpectationNotMet(err) {
							n.Log.WithTags(spec.Tags{C: nil, L: "W", O: n, V: 7}, "%#v", maskAny(err))

							err = n.forwardInputCLG(calculatedPayload)
							if err != nil {
								n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
							}

							return
						} else if err != nil {
							n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
						}

						// Forward to other CLG's, if necessary.
						err = n.Forward(CLG, calculatedPayload)
						if err != nil {
							n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
						}
					}(payload)
				}
			}
		}(ID, CLG)
	}
}

func (n *network) listenInputCLG() {
	n.Log.WithTags(spec.Tags{C: nil, L: "D", O: n, V: 13}, "call Listen")

	// Listen on TextInput from the outside to receive text requests.
	CLG, err := n.clgByName("input")
	if err != nil {
		n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
	}

	networkID := n.GetID()
	clgChannel := CLG.GetInputChannel()

	for {
		select {
		case <-n.Closer:
			break
		case textRequest := <-n.TextInput:

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

			// Prepare the context and a unique behaviour ID for the input CLG.
			ctxConfig := context.DefaultConfig()
			ctxConfig.Expectation = textRequest.GetExpectation()
			ctxConfig.SessionID = textRequest.GetSessionID()
			ctx, err := context.New(ctxConfig)
			if err != nil {
				n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
				continue
			}
			// TODO write a new CLG tree ID and add it to context
			behaviorID, err := n.Factory().ID().New()
			if err != nil {
				n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
				continue
			}

			// We transform the received input to a network payload to have a
			// conventional data structure within the neural network. Note the
			// following details.
			//
			//     The list of arguments always contains a context as first argument.
			//
			//     Destination is always the behavior ID of the input CLG, since this
			//     one is the connecting building block to other CLGs within the
			//     neural network. This behavior ID is always a new one, because it
			//     will eventually be part of a completely new CLG tree within the
			//     connection space.
			//
			//     Sources is here only the individual network ID to have at least
			//     any reference of origin.
			//
			payloadConfig := api.DefaultNetworkPayloadConfig()
			payloadConfig.Args = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(textRequest.GetInput())}
			payloadConfig.Destination = behaviorID
			payloadConfig.Sources = []spec.ObjectID{networkID}
			newPayload, err := api.NewNetworkPayload(payloadConfig)
			if err != nil {
				n.Log.WithTags(spec.Tags{C: nil, L: "E", O: n, V: 4}, "%#v", maskAny(err))
				continue
			}

			// Send the new network payload to the input CLG.
			clgChannel <- newPayload
		}
	}
}