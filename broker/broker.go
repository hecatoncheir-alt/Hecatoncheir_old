package broker

import (
	"fmt"
	"log"

	nsq "github.com/bitly/go-nsq"
)

// New constructor for Broker
func New() *Broker {
	broker := Broker{}
	broker.сonfiguration = nsq.NewConfig()
	return &broker
}

// Broker is a object of message stream
type Broker struct {
	IP            string
	Port          int
	сonfiguration *nsq.Config
	Producer      *nsq.Producer
}

// connectToMessageBroker method for connect to message broker
func (broker *Broker) connectToMessageBroker(host string, port int) (*nsq.Producer, error) {
	if host != "" && string(port) != "" {
		broker.IP = host
		broker.Port = port
	}

	hostAddr := fmt.Sprintf("%v:%v", broker.IP, string(broker.Port))
	producer, err := nsq.NewProducer(hostAddr, broker.сonfiguration)

	if err != nil {
		// log.Panic("Could not connect")
		log.Print("Could not connect to message broker")
	}

	return producer, err

}

// Connect to message broker for publish events
func (broker *Broker) Connect(host string, port int) error {
	producer, err := broker.connectToMessageBroker(host, port)
	broker.Producer = producer
	return err
}

// func (broker Broker) WriteToTopic(topic string, message interface{}) {

// }

// AddHandlerToTopic method for adding handler for listing events
func (broker *Broker) AddHandlerToTopic(topic string, channel string, handler func(event string)) {}

// ListenTopic get events in channel of topic
func (broker *Broker) ListenTopic(topic string, channel string) (<-chan string, error) {
	consumer, err := nsq.NewConsumer(topic, channel, broker.сonfiguration)
	if err != nil {
		return nil, err
	}

	events := make(chan string)

	handler := nsq.HandlerFunc(func(message *nsq.Message) error {
		fmt.Println(message)
		return nil
	})

	consumer.AddHandler(handler)

	hostAddr := fmt.Sprintf("%v:%v", broker.IP, string(broker.Port))
	consumer.ConnectToNSQD(hostAddr)

	return events, nil
}