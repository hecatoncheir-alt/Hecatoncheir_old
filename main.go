package main

import (
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

	// TODO listening broker channel
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
