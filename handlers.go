package main

//
//import (
//	"encoding/json"
//	"log"
//
//	"github.com/hecatoncheir/Hecatoncheir/broker"
//	"github.com/hecatoncheir/Hecatoncheir/crawler"
//	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
//)
//
//// MessageEvent struct for incoming event
//type MessageEvent struct {
//	Message string
//	Data    interface{}
//}
//
//// SubscribeCrawlerHandler handler for crawling items from coduments
//func SubscribeCrawlerHandler(broker *broker.Broker, inputTopic string, outputTopic string) {
//	events, err := broker.ListenTopic(inputTopic, "Crawler")
//	if err != nil {
//		log.Println(err)
//	}
//
//	go func(channel <-chan []byte) {
//		for event := range events {
//			details := MessageEvent{}
//			json.Unmarshal(event, &details)
//
//			switch details.Message {
//			case "Get items from categories of company":
//				var company = crawler.Company{}
//				dataBytes, _ := json.Marshal(details.Data)
//				json.Unmarshal(dataBytes, &company)
//
//				if company.IRI == "http://www.mvideo.ru/" {
//					hecatonhair := mvideo.NewCrawler()
//					var configuration = mvideo.EntityConfig{}
//					json.Unmarshal(dataBytes, &configuration)
//
//					go func(config mvideo.EntityConfig) {
//						go hecatonhair.RunWithConfiguration(config)
//
//						go func() {
//							for item := range hecatonhair.Items {
//								data := map[string]interface{}{"Item": item}
//
//								message := MessageEvent{Message: "Product of company category parsed", Data: data}
//								broker.WriteToTopic(outputTopic, message)
//							}
//						}()
//					}(configuration)
//				}
//
//				if company.IRI == "https://www.ulmart.ru/" {
//					hecatonhair := ulmart.NewCrawler()
//					var configuration = ulmart.EntityConfig{}
//					json.Unmarshal(dataBytes, &configuration)
//
//					go func(config ulmart.EntityConfig) {
//						go hecatonhair.RunWithConfiguration(config)
//
//						go func() {
//							for item := range hecatonhair.Items {
//								data := map[string]interface{}{"Item": item}
//
//								message := MessageEvent{Message: "Product of company category parsed", Data: data}
//								broker.WriteToTopic(outputTopic, message)
//							}
//						}()
//					}(configuration)
//				}
//
//			}
//		}
//	}(events)
//
//}
