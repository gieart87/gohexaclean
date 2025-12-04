package broker

import "errors"

// Message broker infrastructure errors
var (
	ErrBrokerConnection    = errors.New("broker connection failed")
	ErrBrokerPublish       = errors.New("failed to publish message")
	ErrBrokerSubscribe     = errors.New("failed to subscribe to topic")
	ErrBrokerTimeout       = errors.New("broker operation timeout")
	ErrBrokerChannelClosed = errors.New("broker channel closed")
	ErrBrokerAck           = errors.New("failed to acknowledge message")
	ErrBrokerNack          = errors.New("failed to negatively acknowledge message")
)
