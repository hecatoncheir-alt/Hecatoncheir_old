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
	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	bro := broker.New()
	err = bro.Connect(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Connect to channel: '%v'.", config.Production.Channel))

	channel, err := bro.ListenTopic(config.ApiVersion, config.Production.Channel)
	if err != nil {
		log.Println(err)
	}

	log.Println(fmt.Sprintf("Listen topic: '%v' on channel: '%v'", config.ApiVersion, config.Production.Channel))

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			log.Println(fmt.Sprintf("Received message: '%v'", data["Message"]))
			go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, config.ApiVersion)
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
