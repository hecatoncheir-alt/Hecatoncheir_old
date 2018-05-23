package mvideo

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/logger"
)

type MockLogWriter struct{}

func (mockLogWriter *MockLogWriter) Write(logData logger.LogData) error {
	return nil
}

func TestCutSymbolFromPriceOfItem(test *testing.T) {
	price := "7 990¤"

	patternForCutPrice, err := regexp.Compile("[ ¤]*")
	// Or use patternForCutPrice.FindAllString with:
	// patternForCutPrice, err := regexp.Compile("[0-9]*) ([0-9]*)") then cut symbols
	// priceOfItem = strings.Replace(priceOfItem, string([]byte{194, 160}), "", 1)

	if err != nil {
		test.Error(err)
	}

	replacedPrice := patternForCutPrice.ReplaceAllString(price, "")

	if replacedPrice != "7990" {
		test.Fatalf("Expected '7990' price value, actual: %v", replacedPrice)
	}

	realPrice, err := strconv.ParseFloat(replacedPrice, 64)
	if err != nil {
		test.Error(err)
	}

	if realPrice != 7990 {
		test.Fatalf("Expected 7990 price value, actual: %v", realPrice)
	}
}

func TestCutHTMLSymbolFromPriceOfItem(test *testing.T) {
	price := "7&nbsp;990.32¤"

	patternForCutPrice, err := regexp.Compile("(&nbsp;)?[¤ ]*")
	if err != nil {
		test.Error(err)
	}

	replacedPrice := patternForCutPrice.ReplaceAllLiteralString(price, "")

	if replacedPrice != "7990.32" {
		test.Fatalf("Expected '7990' price value, actual: %v", replacedPrice)
	}

	realPrice, err := strconv.ParseFloat(replacedPrice, 64)
	if err != nil {
		test.Error(err)
	}

	if realPrice != 7990.32 {
		test.Fatalf("Expected 7990 price value, actual: %v", realPrice)
	}
}

func TestCrawlerCanGetDocumentByConfig(test *testing.T) {

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

		if product.Language != "ru" {
			test.Fatalf("Expected language \"ru\", actual: %v", product.Language)
		}

		if product.PreviewImageLink == "" {
			test.Fatalf("Expected preview link of product, actual: %v", product.PreviewImageLink)
		}

		if product.IRI == "" {
			test.Fatalf("Expected link to product, actual: %v", product.IRI)
		}

		if product.Category.ID != "0x2787" {
			test.Fatalf("Expected ID of category \"0x2787\", actual: %v", product.Category.ID)
		}

		if product.Category.Name != "Тестовая категория компании М.Видео" {
			test.Fatalf(
				"Expected Name of category \"Тестовая категория компании М.Видео\", actual: %v",
				product.Category.Name)
		}

		if product.Company.ID != "0x2786" {
			test.Fatalf("Expected ID of company \"0x2786\", actual: %v", product.Company.ID)
		}

		if product.Company.Name != "М.Видео" {
			test.Fatalf("Expected Name of company \"М.Видео\", actual: %v",
				product.Category.Name)
		}

		if product.Company.IRI != "http://www.mvideo.ru/" {
			test.Fatalf(
				"Expected IRI of company \"http://www.mvideo.ru/\", actual: %v",
				product.Company.IRI)
		}

		if product.City.ID != "0x2788" {
			test.Fatalf("Expected ID of city \"0x2788\", actual: %v", product.City.ID)
		}

		if product.City.Name != "Москва" {
			test.Fatalf("Expected name of city \"Москва\", actual: %v", product.City.Name)
		}

		if product.Price.Value == 0 {
			test.Fatalf("Expected positive price value, actual: %v", product.Price.Value)
		}

		if product.Price.City.ID != "0x2788" {
			test.Fatalf("Expected ID of city of price \"0x2788\", actual: %v", product.Price.City.ID)
		}

		if product.Price.City.Name != "Москва" {
			test.Fatalf("Expected name of city of price \"Москва\", actual: %v", product.Price.City.Name)
		}

		break
	}
}
