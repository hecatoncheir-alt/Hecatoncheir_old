package broker

import (
	"encoding/json"
	"testing"
)

func TestBrokerCanConnectToNSQ(test *testing.T) {
	broker := New()

	err := broker.Connect("192.168.99.100", 4150)
	if err != nil {
		test.Error(err)
	}

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	broker.Producer.Publish("Hecatoncheir", message)

	items, err := broker.ListenTopic("Hecatoncheir", "ParsedItems")
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		data := map[string]string{}
		json.Unmarshal(item, &data)
		if data["test key"] == "test value" {
			break
		}
	}
}
