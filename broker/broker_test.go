package broker

import (
	"encoding/json"
	"log"
	"sync"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/configuration"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

var broker *Broker
var once sync.Once

func SetUp() {
	broker = New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println(err)
	}

	err = broker.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println("Need started NSQ")
		log.Println(err)
	}
}

func TestBrokerCanConnectToNSQ(test *testing.T) {
	once.Do(SetUp)

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	broker.Producer.Publish("test", message)

	items, err := broker.ListenTopic("test", "testing")
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

	item := crawler.Product{Name: "test item"}

	items, err := broker.ListenTopic("test", "Parser")
	if err != nil {
		test.Error(err)
	}

	err = broker.WriteToTopic("test", item)
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		data := crawler.Product{}
		json.Unmarshal(item, &data)
		if data.Name == "test item" {
			break
		}
	}
}
