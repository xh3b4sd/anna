package event

var (
	ExtendHeadType Type = "extend-head"
)

// Config represents the configuration used to create a new queue object.
type ExtendHeadConfig struct {
	// Connection represents the new connection being tracked during the current
	// event. This connection consist out of two peers. The first peer is
	// Destination. The second peer is Source.
	Connection string

	// ConnectionPath represents the stored connection path matching the new
	// connection according to the event being tracked. In this case,
	// ExtendHeadType.
	ConnectionPath string

	// Destination represents the destination of the network payload currently
	// being processed.
	Destination string

	// Source represents one source of the network payload currently being
	// processed.
	Source string
}

// DefaultEventQueueConfig provides a default configuration to create a new
// extend head object by best effort.
func DefaultExtendHeadConfig() ExtendHeadConfig {
	newConfig := ExtendHeadConfig{
		Connection:     "",
		ConnectionPath: "",
		Destination:    "",
		Source:         "",
	}

	return newConfig
}

// New creates a new configured extend head object.
func New(config ExtendHeadConfig) (Event, error) {
	newEvent := &extendHead{
		ExtendHeadConfig: config,

		Type: ExtendHeadType,
	}

	if newEvent.Connection == "" {
		return nil, maskAnyf(invalidConfigError, "connection must not be empty")
	}
	if newEvent.ConnectionPath == "" {
		return nil, maskAnyf(invalidConfigError, "connection path must not be empty")
	}
	if newEvent.Destination == "" {
		return nil, maskAnyf(invalidConfigError, "destination must not be empty")
	}
	if newEvent.Source == "" {
		return nil, maskAnyf(invalidConfigError, "source must not be empty")
	}

	return newEvent, nil
}

type extendHead struct {
	ExtendHeadConfig

	Type Type
}

func (eh *extendHead) GetConnection() string {
	return eh.Connection
}

func (eh *extendHead) GetConnectionPath() string {
	return eh.ConnectionPath
}

func (eh *extendHead) GetDestination() string {
	return eh.Destination
}

func (eh *extendHead) GetSource() string {
	return eh.Source
}

func (eh *extendHead) GetType() Type {
	return eh.Type
}