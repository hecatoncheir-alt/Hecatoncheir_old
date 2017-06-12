package mvideo

import "github.com/hecatoncheir/Hecatoncheir/crawler"

// ItemConfig structure for parse one item on page
type ItemConfig struct {
	ItemSelector        string
	NameOfItemSelector  string
	PriceOfItemSelector string
}

// Page structure for parse parameters from one page
type Page struct {
	ItemConfig
	CityParamPath            string
	CityParam                string
	Path                     string
	PageInPaginationSelector string
	PageParamPath            string
}

// EntityConfig agregate pages structures
type EntityConfig struct {
	crawler.Company
	Pages []Page
}
