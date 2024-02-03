package pulse

import (
	"context"
	"encoding/json"
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

	return &Message{
		client: client,
	}, nil

}

func (m Message) Publisher(message interface{}, topic string) error {
	pro, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return err
	}
	defer pro.Close()

	// Convert the type interface to JSON
	byteRes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := pulsar.ProducerMessage{
		Payload: byteRes,
	}
	_, err = pro.Send(context.Background(), &msg)
	if err != nil {
		return err
	}
	return nil
}

/*func (m Message) Subscriber(topic string) ([]byte, error) {
cons, err := m.client.Subscribe(pulsar.ConsumerOptions{
	Topic:            topic,
	Type:             pulsar.Exclusive,
	SubscriptionName: "my-sub-name",
})
if err != nil {
	return nil, err
}
defer cons.Close()

msg, err := cons.Receive(context.Background())
if err != nil {
	log.Fatal(err)
}
/* trunk-ignore(golangci-lint/errcheck) */
//cons.Ack(msg)

//return msg.Payload(), nil
//}

func (m Message) Subscriber(topic string) ([]byte, error) {
	cons, err := m.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		Type:             pulsar.Exclusive,
		SubscriptionName: "my-sub-name",
	})
	if err != nil {
		return nil, err
	}
	defer cons.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		msg, err := cons.Receive(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil, ctx.Err()
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}

		/* trunk-ignore(golangci-lint/errcheck) */
		cons.Ack(msg)
		return msg.Payload(), nil
	}
}
