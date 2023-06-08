package pulse

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	store "github.com/DecodeWorms/messaging-protocol"
	"github.com/apache/pulsar-client-go/pulsar"
)

type EventStore struct {
	client pulsar.Client
	//add logger
}

func Init(c pulsar.Client, opt store.Options) (store.PulsarStore, error) {
	opts := pulsar.ClientOptions{
		URL: opt.Address,
	}
	opts.TLSAllowInsecureConnection = true
	c, err := pulsar.NewClient(opts)
	if err != nil {
		return nil, err
	}
	return &EventStore{
		client: c,
	}, nil
}

func (e *EventStore) Publish(topic string, message []byte) error {
	opts := pulsar.ProducerOptions{
		Topic:                   topic,
		Name:                    fmt.Sprintf("%s", generateRandomName()),
		DisableBlockIfQueueFull: true,
		SendTimeout:             0,
	}

	prod, err := e.client.CreateProducer(opts)
	if err != nil {
		return err
	}
	defer prod.Close()
	_, err = prod.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *EventStore) PublishRaw(topic string, message interface{}) error {
	opts := pulsar.ProducerOptions{
		Topic:                   topic,
		Name:                    fmt.Sprintf("%s", generateRandomName()),
		DisableBlockIfQueueFull: true,
		SendTimeout:             0,
	}
	producer, err := e.client.CreateProducer(opts)
	if err != nil {
		return err
	}
	defer producer.Close()

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = e.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventStore) Subscribe(topic string, subscription string, messageHandler func(message pulsar.Message)) error {
	opts := pulsar.ConsumerOptions{
		Topic: topic,
	}
	consumer, err := e.client.Subscribe(opts)
	if err != nil {
		return err
	}

	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			return err
		}

		messageHandler(msg)
		//consumer.Ack(msg)
	}
}

func (e *EventStore) Run() error {
	//TODO implement me
	panic("implement me")
}

func generateRandomName() string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, 10)
	for i := range bytes {
		bytes[i] = chars[rand.Intn(len(chars))]
	}
	return string(bytes)
}
