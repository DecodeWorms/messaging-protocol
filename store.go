package store

import "github.com/apache/pulsar-client-go/pulsar"

type PulsarStore interface {
	Publish(topic string, message []byte) error
	PublishRaw(topic string, message interface{}) error
	Subscribe(topic string, subscription string, messageHandler func(message pulsar.Message)) error
	Run() error
}
