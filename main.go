package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
	"github.com/hecatoncheir/Hecatoncheir/logger"
)

func main() {
	config := configuration.New()

	bro := broker.New()
	err := bro.Connect(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	loguna := logger.New(config.APIVersion, config.Production.LogunaTopic, bro)
	eventMessage := fmt.Sprintf("Configuration of hecatoncheir: '%v'.", config)
	log.Println(eventMessage)

	err = loguna.Write(logger.LogData{Message: eventMessage, Level: "info"})
	if err != nil {
		log.Println(err)
	}

	eventMessage = fmt.Sprintf("Connect to channel: '%v'.", config.Production.HecatoncheirTopic)
	log.Println(eventMessage)

	event := logger.LogData{
		Message: eventMessage,
		Level:   "info"}

	err = loguna.Write(event)
	if err != nil {
		log.Println(err)
	}

	channel, err := bro.ListenTopic(config.APIVersion, config.Production.HecatoncheirTopic)
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

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			eventMessage = fmt.Sprintf("Received message: '%v'", data["Message"])
			log.Println(eventMessage)

			event = logger.LogData{
				Message: eventMessage,
				Level:   "info"}

			err = loguna.Write(event)
			if err != nil {
				log.Println(err)
			}

			go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, config.APIVersion, loguna)
		}
	}
}

func handlesNeedProductsOfCategoryOfCompanyEvent(parserInstructionsJSON string, bro *broker.Broker, topic string, loguna *logger.LogWriter) {
	parserInstructionsOfCompany, err := crawler.NewParserInstructionsFromJSON(parserInstructionsJSON)
	if err != nil {
		log.Println(err)
		loguna.Write(logger.LogData{Message: err.Error(), Level: "warning"})
	}

	if parserInstructionsOfCompany.Company.IRI == "http://www.mvideo.ru/" {
		crawlerOfCompany := mvideo.New(loguna)

		channelWithProducts, err := crawlerOfCompany.RunWithConfiguration(parserInstructionsOfCompany)
		if err != nil {
			log.Println(err)
			loguna.Write(logger.LogData{Message: err.Error(), Level: "warning"})
		}

		for product := range channelWithProducts {
			go bro.WriteToTopic(topic,
				map[string]interface{}{"Message": "Product of category of company ready", "Data": product})
		}
	}

}
