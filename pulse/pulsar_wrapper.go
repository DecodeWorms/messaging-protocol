package pulse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	messagingprotocol "github.com/DecodeWorms/messaging-protocol"
	"github.com/apache/pulsar-client-go/pulsar"
)

type Message struct {
	client pulsar.Client
	/* trunk-ignore(golangci-lint/unused) */
	opt         messagingprotocol.Options
	serviceName string
}

func NewMessage(opt messagingprotocol.Options) (*Message, error) {
	addr := opt.Address
	if addr == "" {
		return nil, fmt.Errorf("empty pulsar url")
	}
	serviceName := opt.ServiceName
	if serviceName == "" {
		return nil, fmt.Errorf("empty service name")
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                        addr,
		TLSAllowInsecureConnection: true,
	})
	if err != nil {
		return nil, err
	}

	return &Message{
		client:      client,
		serviceName: serviceName,
	}, nil

}

func (m Message) Publisher(message interface{}, topic string) error {
	pro, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Name:  generateRandomName(),
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
		Payload:   byteRes,
		EventTime: time.Now(),
	}
	_, err = pro.Send(context.Background(), &msg)
	if err != nil {
		return err
	}
	return nil
}

func (m Message) Subscriber(topic string) ([]byte, error) {
	cons, err := m.client.Subscribe(pulsar.ConsumerOptions{
		Topic: topic,
		Type:  pulsar.Exclusive,
		// Remove the hard coded part
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

func generateRandomName() string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, 10)
	for i := range bytes {
		bytes[i] = chars[rand.Intn(len(chars))]
	}
	return string(bytes)
}
