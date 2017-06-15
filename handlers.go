package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
	"github.com/hecatoncheir/Hecatoncheir/crawler/ulmart"
)

// MessageEvent struct for incoming event
type MessageEvent struct {
	Message string
	Data    interface{}
}

// SubscribeCrawlerHandler handler for crawling items from coduments
func SubscribeCrawlerHandler(broker *broker.Broker, channel string) {
	events, err := broker.ListenTopic("Hecatoncheir", channel)
	if err != nil {
		log.Println(err)
	}

	go func(channel <-chan []byte) {
		for event := range events {
			var details MessageEvent
			json.Unmarshal(event, &details)

			switch details.Message {
			case "Get items from categories of company":
				var company = crawler.Company{}
				dataBytes, _ := json.Marshal(details.Data)
				json.Unmarshal(dataBytes, &company)
				fmt.Println(company)

				if company.IRI == "http://www.mvideo.ru/" {
					hecatonhair := mvideo.NewCrawler()
					var configuration = mvideo.EntityConfig{}
					json.Unmarshal(dataBytes, &configuration)

					go func(config mvideo.EntityConfig) {

						go hecatonhair.RunWithConfiguration(config)

						go func() {
							for item := range hecatonhair.Items {
								data := map[string]interface{}{"Item": item}

								message := MessageEvent{Message: "Item from categories of company parsed", Data: data}
								broker.WriteToTopic("Hecatoncheir", message)
							}
						}()
					}(configuration)
				}

				if company.IRI == "https://www.ulmart.ru/" {
					fmt.Println("AAAA")
					hecatonhair := ulmart.NewCrawler()
					var configuration = ulmart.EntityConfig{}
					json.Unmarshal(dataBytes, &configuration)

					go func(config ulmart.EntityConfig) {

						go hecatonhair.RunWithConfiguration(config)

						go func() {
							for item := range hecatonhair.Items {
								data := map[string]interface{}{"Item": item}

								message := MessageEvent{Message: "Item from categories of company parsed", Data: data}
								broker.WriteToTopic("Hecatoncheir", message)
							}
						}()
					}(configuration)
				}

			}
		}
	}(events)

}
