package main

import (
	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"

	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

func TestIntegrationCanParseCategoryOfCompanyByBrokerEventRequest(test *testing.T) {
	config := configuration.New()
	bro := broker.New()

	var err error

	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}

	channel, err := bro.ListenTopic(config.Development.HecatoncheirTopic, config.APIVersion)
	if err != nil {
		test.Error(err)
	}

	parserOfCompany := crawler.ParserOfCompanyInstructions{
		Language: "ru",
		Company: crawler.Company{
			ID:   "0x2786",
			Name: "М.Видео",
			IRI:  "http://www.mvideo.ru/"},
		Category: crawler.Category{
			ID:   "0x2787",
			Name: "Тестовая категория компании М.Видео"},
		City: crawler.City{
			ID:   "0x2788",
			Name: "Москва"},
		PageInstruction: crawler.PageInstruction{
			ID:   "0x2789",
			Path: "smartfony-i-svyaz/smartfony-205",
			PageInPaginationSelector: ".c-pagination > .c-pagination__num",
			PageParamPath:            "/f/page=",
			CityParamPath:            "?cityId=",
			//CityParam:                  "CityCZ_975",
			ItemSelector:               ".c-product-tile",
			PreviewImageOfItemSelector: ".c-product-tile-picture__link .lazy-load-image-holder img",
			NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"},
	}

	parseData, err := json.Marshal(parserOfCompany)
	if err != nil {
		test.Error(err)
	}

	go bro.WriteToTopic(config.Development.HecatoncheirTopic, map[string]interface{}{"Message": "Need products of category of company", "Data": string(parseData)})

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			test.Fail()
		}

		go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, config.Development.LogunaTopic, nil)

		break
	}

	channelForGetProducts, err := bro.ListenTopic(config.Development.LogunaTopic, config.APIVersion)
	if err != nil {
		test.Error(err)
	}

	for message := range channelForGetProducts {
		data := map[string]interface{}{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Product of category of company ready" {
			test.Fail()
		}

		if data["Data"].(map[string]interface{})["Language"] != "ru" {
			test.Fail()
		}

		break
	}
}
