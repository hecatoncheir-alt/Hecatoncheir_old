package mvideo

import (
	"testing"
	"time"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

func TestCrawlerCanGetDocumentByConfig(test *testing.T) {
	smartphonesPage := Page{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		ItemConfig: ItemConfig{
			ItemSelector:               ".grid-view .product-tile",
			NameOfItemSelector:         ".product-tile-title",
			PriceOfItemSelector:        ".product-price-current",
			LinkOfItemSelector:         ".product-tile-title a",
			PreviewImageOfItemSelector: ".product-tile-picture-link img",
		},
	}

	configuration := EntityConfig{
		Company: crawler.Company{
			IRI:        "http://www.mvideo.ru/",
			Name:       "M.Video",
			Categories: []string{"Телефоны"},
		},
		Pages: []Page{smartphonesPage},
	}

	mechanism := NewCrawler()

	go mechanism.RunWithConfiguration(configuration)

	isRightItems := false

	go func() {
		time.Sleep(time.Second * 6)
		close(mechanism.Items)
	}()

	for item := range mechanism.Items {
		if item.Name != "" && item.Price.Value != "" && item.Link != "" {
			isRightItems = true
			break
		}
	}

	if isRightItems == false {
		test.Fail()
	}
}
