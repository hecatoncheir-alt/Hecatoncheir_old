package main

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
	"github.com/prometheus/common/log"
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

	channel, err := bro.ListenTopic(config.ApiVersion, "Parser")
	if err != nil {
		log.Error(err)
	}

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, config.ApiVersion)
		}
	}
}

func handlesNeedProductsOfCategoryOfCompanyEvent(parserInstructionsJSON string, bro *broker.Broker, topic string) {
	fmt.Println(parserInstructionsJSON)
	parserInstructionsOfCompany, err := crawler.NewParserInstructionsFromJSON(parserInstructionsJSON)
	if err != nil {
		log.Error(err)
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
