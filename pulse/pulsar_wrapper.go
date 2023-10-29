package pulse

import (
	"context"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Message struct {
	client pulsar.Client
}

func NewMessage(pulsarUrl string) (*Message, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: pulsarUrl,
	})
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return &Message{
		client: client,
	}, nil

}

func (m Message) Publisher(message, topic string) error {
	pro, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return err
	}
	defer pro.Close()

	msg := pulsar.ProducerMessage{
		Payload: []byte(message),
	}
	_, err = pro.Send(context.Background(), &msg)
	if err != nil {
		return err
	}
	return nil
}

func (m Message) Subscriber(topic string) string {
	cons, err := m.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		Type:             pulsar.Exclusive,
		SubscriptionName: "my-sub-name",
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	defer cons.Close()

	msg, err := cons.Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cons.Ack(msg)
	return string(msg.Payload())

}
