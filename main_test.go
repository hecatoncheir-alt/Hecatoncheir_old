package main

//func TestIntegrationCanParseCategoryOfCompanyByBrokerEventRequest(test *testing.T) {
//	test.Skip()
//
//	config := configuration.New()
//	if config.ServiceName == "" {
//		config.ServiceName = "Hecatoncheir"
//	}
//
//	bro := broker.New(config.APIVersion, config.ServiceName)
//
//	err := bro.Connect(config.Development.EventBus.Host, config.Development.EventBus.Port)
//	if err != nil {
//		test.Error(err)
//	}
//
//	parserOfCompany := crawler.ParserOfCompanyInstructions{
//		Language: "ru",
//		Company: crawler.Company{
//			ID:   "0x2786",
//			Name: "М.Видео",
//			IRI:  "http://www.mvideo.ru/"},
//		Category: crawler.Category{
//			ID:   "0x2787",
//			Name: "Тестовая категория компании М.Видео"},
//		City: crawler.City{
//			ID:   "0x2788",
//			Name: "Москва"},
//		PageInstruction: crawler.PageInstruction{
//			ID:   "0x2789",
//			Path: "smartfony-i-svyaz/smartfony-205",
//			PageInPaginationSelector: ".c-pagination > .c-pagination__num",
//			PageParamPath:            "/f/page=",
//			CityParamPath:            "?cityId=",
//			//CityParam:                  "CityCZ_975",
//			ItemSelector:               ".c-product-tile",
//			PreviewImageOfItemSelector: ".c-product-tile-picture__link .lazy-load-image-holder img",
//			NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
//			LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
//			PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"},
//	}
//
//	parseData, err := json.Marshal(parserOfCompany)
//	if err != nil {
//		test.Error(err)
//	}
//
//	otherService := broker.New(config.APIVersion, config.ServiceName)
//
//	err = otherService.Connect(config.Development.EventBus.Host, config.Development.EventBus.Port)
//	if err != nil {
//		test.Error(err)
//	}
//
//	event := broker.EventData{Message: "Need products of category of company", Data: string(parseData)}
//	go otherService.Write(event)
//
//	for data := range bro.InputChannel {
//
//		if data.Message == "Need products of category of company" {
//			break
//		}
//
//		if data.Message != "Need products of category of company" {
//			test.Fail()
//		}
//
//		go handlesNeedProductsOfCategoryOfCompanyEvent(data.Data, bro, config.Development.SprootTopic, nil)
//	}
//
//	for data := range otherService.InputChannel {
//
//		if data.Message == "Product of category of company ready" {
//			break
//		}
//
//		if data.Message != "Product of category of company ready" {
//			test.Fatalf("Expected \"Product of category of company ready\" message, but actual: %v", data.Message)
//		}
//
//		dataOfEvent := map[string]interface{}{}
//		err = json.Unmarshal([]byte(data.Data), &dataOfEvent)
//		if err != nil {
//			test.Error(err)
//		}
//
//		if dataOfEvent["Language"] == "ru" {
//			break
//		}
//
//		if dataOfEvent["Language"] != "ru" {
//			test.Fatalf("Expected \"ru\" language, but actual: %v", dataOfEvent["Language"])
//		}
//	}
//}
