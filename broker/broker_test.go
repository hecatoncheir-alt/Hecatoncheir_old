package broker

import (
	"encoding/json"
	"log"
	"sync"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

var broker *Broker
var once sync.Once

func SetUp() {
	broker = New()

	err := broker.Connect("192.168.99.100", 4150)
	if err != nil {
		log.Println(err)
	}
}

func TestBrokerCanConnectToNSQ(test *testing.T) {
	once.Do(SetUp)

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
	once.Do(SetUp)

	item := crawler.Item{Name: "test item"}

	items, err := broker.ListenTopic("Hecatoncheir", "items")
	if err != nil {
		test.Error(err)
	}

	err = broker.WriteToTopic("Hecatoncheir", item)
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
