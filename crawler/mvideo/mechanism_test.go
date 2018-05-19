package mvideo

import (
	"testing"

	"fmt"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/logger"
)

type MockLogWriter struct{}

func (mockLogWriter *MockLogWriter) Write(logData logger.LogData) error {
	return nil
}

func TestCrawlerCanGetDocumentByConfig(test *testing.T) {

	parserOfCompany := crawler.ParserOfCompanyInstructions{
		Language: "en",
		Company: crawler.Company{
			ID:   "0x2786",
			Name: "M.Video",
			IRI:  "http://www.mvideo.ru/"},
		Category: crawler.Category{
			ID:   "",
			Name: "Test category of M.Video company"},
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
			PreviewImageOfItemSelector: ".product-tile-picture__image img",
			NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"},
	}

	logWriter := MockLogWriter{}
	mechanism := New(&logWriter)

	productsChannel, err := mechanism.RunWithConfiguration(parserOfCompany)
	if err != nil {
		test.Fatal(err)
	}

	for product := range productsChannel {
		if product.Name == "" {
			test.Fatalf("Expect name of product but get: %v", product.Name)
		}
		fmt.Println(product)
		break
	}

	// isRightItems := false

	// go func() {
	// 	time.Sleep(time.Second * 6)
	// 	close(mechanism.Items)
	// }()

	// for item := range mechanism.Items {
	// 	fmt.Println(item)
	// 	if item.Name != "" && item.Price.Value != "" && item.IRI != "" {
	// 		isRightItems = true
	// 		break
	// 	}
	// }

	// if isRightItems == false {
	// 	test.Fail()
	// }
}
