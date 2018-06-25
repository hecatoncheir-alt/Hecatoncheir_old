package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
	"github.com/hecatoncheir/Logger"
)

func main() {
	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Hecatoncheir"
	}

	bro := broker.New(config.APIVersion, config.ServiceName)
	err := bro.Connect(config.Production.EventBus.Host,
		config.Production.EventBus.Port)
	if err != nil {
		log.Fatal(err)
	}

	loguna := logger.New(config.APIVersion,
		config.ServiceName, config.Production.LogunaTopic, bro)

	eventMessage :=
		fmt.Sprintf("Configuration of hecatoncheir: '%v'.", config)

	//log.Println(eventMessage)

	err = loguna.Write(logger.LogData{Message: eventMessage, Level: "info"})
	if err != nil {
		log.Println(err)
	}

	eventMessage = fmt.Sprintf("Connect to channel: '%v'.",
		config.Production.HecatoncheirTopic)
	log.Println(eventMessage)

	event := logger.LogData{
		Message: eventMessage,
		Level:   "info"}

	err = loguna.Write(event)
	if err != nil {
		log.Println(err)
	}

	eventMessage = fmt.Sprintf("Listen topic: '%v' on channel: '%v'", config.APIVersion, config.Production.HecatoncheirTopic)
	log.Println(eventMessage)

	event = logger.LogData{
		Message: eventMessage,
		Level:   "info"}

	err = loguna.Write(event)
	if err != nil {
		log.Println(err)
	}

	for data := range bro.InputChannel {

		if data.Message != "Need products of category of company" {
			eventMessage = fmt.Sprintf("Received message: '%v'", data.Message)
			log.Println(eventMessage)
			fmt.Println(eventMessage)

			event = logger.LogData{
				Message: eventMessage,
				Level:   "info"}

			err = loguna.Write(event)
			if err != nil {
				log.Println(err)
			}

			go handlesNeedProductsOfCategoryOfCompanyEvent(data.Data, bro, config.APIVersion, loguna)
		}
	}
}

func handlesNeedProductsOfCategoryOfCompanyEvent(parserInstructionsJSON string, bro *broker.Broker, topicForProductsPush string, loguna logger.Writer) {

	parserInstructionsOfCompany, err := crawler.NewParserInstructionsFromJSON(parserInstructionsJSON)
	if err != nil {
		log.Println(err)
		err := loguna.Write(logger.LogData{Message: err.Error(), Level: "warning"})
		if err != nil {
			log.Println(err)
		}
	}

	if parserInstructionsOfCompany.Company.IRI == "http://www.mvideo.ru/" {
		crawlerOfCompany := mvideo.New(loguna)

		channelWithProducts, err := crawlerOfCompany.RunWithConfiguration(parserInstructionsOfCompany)
		if err != nil {
			log.Println(err)
			err := loguna.Write(logger.LogData{Message: err.Error(), Level: "warning"})
			if err != nil {
				log.Println(err)
			}
		}

		for product := range channelWithProducts {
			data, err := json.Marshal(product)
			if err != nil {
				log.Println(err)
				err := loguna.Write(logger.LogData{Message: err.Error(), Level: "warning"})
				if err != nil {
					log.Println(err)
				}
			}

			event := broker.EventData{Message: "Product of category of company ready", Data: string(data)}
			fmt.Printf("Write: %v to %v", event, topicForProductsPush)
			go bro.Write(event)
		}
	}

}
