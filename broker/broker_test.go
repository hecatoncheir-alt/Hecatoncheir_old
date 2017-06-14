package broker

import (
	"testing"
)

func TestBrokerCanConnectToNSQ(test *testing.T) {
	broker := New()

	err := broker.Connect("192.168.99.100", 4150)
	if err != nil {
		test.Error(err)
	}

	broker.Producer.Publish("Hecatoncheir", []byte("Message"))

	items, err := broker.ListenTopic("Hecatoncheir", "ParsedItems")
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		if item != "" {
			break
		}
	}
}
