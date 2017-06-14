package ulmart

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hecatoncheir/Hecatoncheir/crawler"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestCookies(test *testing.T) {
	iri := "https://www.ulmart.ru/catalog/communicators?sort=5&viewType=2&rec=true"

	cookie, _ := cookiejar.New(nil)
	// 18414 - Москва
	// 1688 - Алексин
	city := &http.Cookie{Name: "city", Value: "18414"}
	allCookies := []*http.Cookie{}
	allCookies = append(allCookies, city)

	pageURL, _ := url.Parse(iri)
	cookie.SetCookies(pageURL, allCookies)
	client := &http.Client{
		Jar: cookie,
	}

	request, _ := http.NewRequest("GET", iri, nil)
	response, _ := client.Do(request)

	document, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		test.Log(err)
	}

	cityName := document.Find("#load-cities").Text()
	cityName = strings.TrimSpace(cityName)

	if cityName != "Москва" {
		test.Fail()
	}
}

func TestCrawlerCanGetDocumentByConfig(test *testing.T) {
	smartphonesPage := Page{
		Path: "catalog/communicators",
		TotalCountItemsOnPageSelector: "#total-show-count",
		MaxItemsOnPageSelector:        "#max-show-count",
		PagePath:                      "catalogAdditional/communicators",
		PageParamPath:                 "?pageNum=",
		CityInCookieKey:               "city",
		CityID:                        "18414",
		ItemConfig: ItemConfig{
			ItemSelector:        ".b-product",
			NameOfItemSelector:  ".b-product__title a",
			PriceOfItemSelector: ".b-product__price .b-price__num",
			LinkOfItemSelector:  ".b-product__title a",
		},
	}

	configuration := EntityConfig{
		Company: crawler.Company{
			IRI:        "https://www.ulmart.ru/",
			Name:       "Ulmart",
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
			assert.Equal(test, item.Price.City, "Москва", "Not right city")
			isRightItems = true
			break
		}
	}

	if isRightItems == false {
		test.Fail()
	}
}
