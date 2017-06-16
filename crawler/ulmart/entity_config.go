package ulmart

import "github.com/hecatoncheir/Hecatoncheir/crawler"

//ItemConfig is structure of parameters for one item on page
type ItemConfig struct {
	ItemSelector               string
	NameOfItemSelector         string
	PriceOfItemSelector        string
	LinkOfItemSelector         string
	PreviewImageOfItemSelector string
}

// Page is structure of parameters for parse one web page
type Page struct {
	ItemConfig
	CityInCookieKey               string
	CityID                        string
	Path                          string
	TotalCountItemsOnPageSelector string
	MaxItemsOnPageSelector        string
	PagePath                      string
	PageParamPath                 string
}

// EntityConfig for make crawler work.
type EntityConfig struct {
	crawler.Company
	Pages []Page
}
