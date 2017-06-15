package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	broker "github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/crawler/mvideo"
	"github.com/hecatoncheir/Hecatoncheir/crawler/ulmart"
	socket "github.com/hecatoncheir/Hecatoncheir/socket"

	"golang.org/x/net/websocket"
)

var (
	once sync.Once
)

func SetUpSocketServer() {
	server := socket.NewEngine("v1.0")
	server.PowerUp("localhost", 8181)
	defer server.Server.Close()
}

// {
// 	"Message": "Get items from categories of company",
// 	"Data": {
// 			"Iri": "http://www.mvideo.ru/",
//			"Name": "M.Video",
//			"Categories": ["Телефоны"],
// 			"Pages": [{
// 				"Path": "smartfony-i-svyaz/smartfony-205",
// 				"PageInPaginationSelector": ".pagination-list .pagination-item",
// 				"PageParamPath": "/f/page=",
// 				"CityParamPath": "?cityId=",
// 				"CityParam": "CityCZ_975",
// 				"ItemSelector": ".grid-view .product-tile",
// 				"NameOfItemSelector": ".product-tile-title",
// 				"PriceOfItemSelector": ".product-price-current"
// 			}]
// 	}
// }
func TestSocketCanParseDocumentOfMvideo(test *testing.T) {
	go once.Do(SetUpSocketServer)

	client, err := websocket.Dial("ws://localhost:8181", "", "http://localhost:8181")

	if err != nil {
		test.Fatal()
	}

	smartphonesPage := mvideo.Page{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		CityParam:                "CityCZ_975",
		ItemConfig: mvideo.ItemConfig{
			ItemSelector:        ".grid-view .product-tile",
			NameOfItemSelector:  ".product-tile-title",
			PriceOfItemSelector: ".product-price-current",
		},
	}

	configuration := mvideo.EntityConfig{
		Company: crawler.Company{
			IRI:        "http://www.mvideo.ru/",
			Name:       "M.Video",
			Categories: []string{"Телефоны"},
		},
		Pages: []mvideo.Page{smartphonesPage},
	}

	websocket.JSON.Send(client, socket.MessageEvent{Message: "Get items from categories of company", Data: configuration})

	message := make(chan socket.MessageEvent)

	go func() {
		for {
			socketEvent := socket.MessageEvent{}
			err := websocket.JSON.Receive(client, &socketEvent)
			if err != nil {
				test.Error(err)
			}
			message <- socketEvent
			break
		}
	}()

	for event := range message {
		if event.Message != "Item from categories of company parsed" ||
			event.Data.(map[string]interface{})["Item"] == nil {
			test.Fail()
		}
		break
	}
}

//{
//	"Message": "Get items from categories of company",
//	"Data": {
//			"Iri": "https://www.ulmart.ru/",
//		"Name": "Ulmart",
//		"Categories": ["Телефоны"],
//			"Pages": [{
//      "Path":                          "catalog/communicators",
//      "TotalCountItemsOnPageSelector": "#total-show-count",
//      "MaxItemsOnPageSelector":        "#max-show-count",
//      "PagePath":                      "catalogAdditional/communicators",
//      "PageParamPath":                 "?pageNum=",
//      "CityInCookieKey":               "city",
//      "CityID":                        "18414",
//				"ItemSelector": ".b-product",
//				"NameOfItemSelector": ".b-product__title a",
//				"PriceOfItemSelector": ".b-product__price .b-price__num"
//			}]
//	}
//}
func TestSocketCanParseDocumentOfUlmart(test *testing.T) {
	go once.Do(SetUpSocketServer)

	client, err := websocket.Dial("ws://localhost:8181", "", "http://localhost:8181")

	if err != nil {
		test.Fatal()
	}

	smartphonesPage := ulmart.Page{
		Path: "catalog/communicators",
		TotalCountItemsOnPageSelector: "#total-show-count",
		MaxItemsOnPageSelector:        "#max-show-count",
		PagePath:                      "catalogAdditional/communicators",
		PageParamPath:                 "?pageNum=",
		CityInCookieKey:               "city",
		CityID:                        "18414",
		ItemConfig: ulmart.ItemConfig{
			ItemSelector:        ".b-product",
			NameOfItemSelector:  ".b-product__title a",
			PriceOfItemSelector: ".b-product__price .b-price__num",
			LinkOfItemSelector:  ".b-product__title a",
		},
	}

	configuration := ulmart.EntityConfig{
		Company: crawler.Company{
			IRI:        "https://www.ulmart.ru/",
			Name:       "Ulmart",
			Categories: []string{"Телефоны"},
		},
		Pages: []ulmart.Page{smartphonesPage},
	}

	websocket.JSON.Send(client, socket.MessageEvent{Message: "Get items from categories of company", Data: configuration})

	message := make(chan socket.MessageEvent)

	go func() {
		for {
			socketEvent := socket.MessageEvent{}
			err := websocket.JSON.Receive(client, &socketEvent)
			if err != nil {
				test.Error(err)
			}
			message <- socketEvent
			break
		}
	}()

	for event := range message {
		if event.Message != "Item from categories of company parsed" ||
			event.Data.(map[string]interface{})["Item"] == nil {
			test.Fail()
		}
		break
	}
}

func TestBrokerMessaging(test *testing.T) {
	var err error

	request := `{
	"Message": "Get items from categories of company",
	"Data": {
		"Iri": "https://www.ulmart.ru/",
		"Name": "Ulmart",
		"Categories": ["Телефоны"],
			"Pages": [{
				"Path":                          "catalog/communicators",
				"TotalCountItemsOnPageSelector": "#total-show-count",
				"MaxItemsOnPageSelector":        "#max-show-count",
				"PagePath":                      "catalogAdditional/communicators",
				"PageParamPath":                 "?pageNum=",
				"CityInCookieKey":               "city",
				"CityID":                        "18414",
				"ItemSelector": ".b-product",
				"NameOfItemSelector": ".b-product__title a",
				"PriceOfItemSelector": ".b-product__price .b-price__num",
				"LinkOfItemSelector": ".b-product__title a"
			}]
		}
	}"`

	bro := broker.New()
	err = bro.Connect("192.168.99.100", 4150)
	if err != nil {
		test.Error(err)
	}

	go SubscribeCrawlerHandler(bro, "CrawlingRequest")

	items, err := bro.ListenTopic("ItemsOfCompanies", "test")
	if err != nil {
		test.Error(err)
	}

	err = bro.WriteToTopic("CrawlingRequest", request)
	if err != nil {
		test.Error(err)
	}

	it := 0
	for data := range items {
		item := crawler.Item{}
		json.Unmarshal(data, &item)
		fmt.Println(item)

		if it < 1 {
			it++
			continue
		}

		if item.Name != "" {
			break
		}

		continue
	}
}
