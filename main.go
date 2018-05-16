package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
)

func main() {
	config := configuration.New()

	bro := broker.New()
	err := bro.Connect(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Connect to channel: '%v'.", config.Production.HecatoncheirTopic))

	channel, err := bro.ListenTopic(config.APIVersion, config.Production.HecatoncheirTopic)
	if err != nil {
		log.Println(err)
	}

	log.Println(fmt.Sprintf("Listen topic: '%v' on channel: '%v'", config.APIVersion, config.Production.HecatoncheirTopic))

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			log.Println(fmt.Sprintf("Received message: '%v'", data["Message"]))
			go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, config.APIVersion)
		}
	}
}

func handlesNeedProductsOfCategoryOfCompanyEvent(parserInstructionsJSON string, bro *broker.Broker, topic string) {
	fmt.Println(parserInstructionsJSON)
	parserInstructionsOfCompany, err := crawler.NewParserInstructionsFromJSON(parserInstructionsJSON)
	if err != nil {
		log.Println(err)
	}

	if parserInstructionsOfCompany.Company.IRI == "http://www.mvideo.ru/" {
		crawlerOfCompany := mvideo.NewCrawler()

		go crawlerOfCompany.RunWithConfiguration(parserInstructionsOfCompany)

		for product := range crawlerOfCompany.Items {
			go bro.WriteToTopic(topic,
				map[string]interface{}{"Message": "Product of category of company ready", "Data": product})
		}
	}

}
