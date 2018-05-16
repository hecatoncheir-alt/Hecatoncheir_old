package broker

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/configuration"
)

func TestBrokerCanConnectToNSQ(test *testing.T) {
	bro := New()

	config, _ := configuration.GetConfiguration()

	err := bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println("Need started NSQ")
		log.Println(err)
	}

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	bro.Producer.Publish(config.Development.HecatoncheirTopic, message)

	items, err := bro.ListenTopic(config.Development.HecatoncheirTopic, config.APIVersion)
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
