package broker

import (
	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

func TestBrokerCanConnectToNSQ(test *testing.T) {
	broker := New()

	err := broker.Connect("192.168.99.100", 4150)
	if err != nil {
		test.Error(err)
	}

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	broker.Producer.Publish("Hecatoncheir", message)

	items, err := broker.ListenTopic("Hecatoncheir", "test")
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

func TestBrokerCanSendMessageToNSQ(test *testing.T) {
	var err error
	broker := New()

	err = broker.Connect("192.168.99.100", 4150)
	if err != nil {
		test.Error(err)
	}

	item := crawler.Item{Name: "test item"}

	err = broker.WriteToTopic("Hecatoncheir", item)
	if err != nil {
		test.Error(err)
	}

	items, err := broker.ListenTopic("Hecatoncheir", "ParsedItem")
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		data := crawler.Item{}
		json.Unmarshal(item, &data)
		if data.Name == "test item" {
			break
		}
	}
}
