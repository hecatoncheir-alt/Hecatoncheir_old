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
			PageInPaginationSelector:   ".pagination-list .pagination-item",
			PreviewImageOfItemSelector: ".product-tile-picture-link img",
			PageParamPath:              "/f/page=",
			CityParamPath:              "?cityId=",
			//CityParam:                  "CityCZ_975",
			ItemSelector:        ".grid-view .product-tile",
			NameOfItemSelector:  ".product-tile-title",
			LinkOfItemSelector:  ".product-tile-title a",
			PriceOfItemSelector: ".product-price-current"},
	}
	fmt.Println(parserOfCompany)

	logWriter := MockLogWriter{}
	mechanism := New(&logWriter)
	fmt.Println(mechanism.LogWriter)
	//
	//go mechanism.RunWithConfiguration(parserOfCompany)
	//
	//isRightItems := false
	//
	//go func() {
	//	time.Sleep(time.Second * 6)
	//	close(mechanism.Items)
	//}()
	//
	//for item := range mechanism.Items {
	//	if item.Name != "" && item.Price.Value != "" && item.IRI != "" {
	//		isRightItems = true
	//		break
	//	}
	//}
	//
	//if isRightItems == false {
	//	test.Fail()
	//}
}
