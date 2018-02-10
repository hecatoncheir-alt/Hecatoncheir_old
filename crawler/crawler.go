package crawler

import (
	"encoding/json"
	"time"
)

// ParserOfCompany is a type for parse company products
type ParserOfCompanyInstructions struct {
	Language        string
	Company         Company
	Category        Category
	City            City
	PageInstruction PageInstruction
}

// NewParserInstructionsFromJSON is a constructor for ParserOfCompany from json string
func NewParserInstructionsFromJSON(data string) (ParserOfCompanyInstructions, error) {

	parser := ParserOfCompanyInstructions{}
	err := json.Unmarshal([]byte(data), &parser)
	if err != nil {
		return parser, err
	}

	return parser, nil
}

// Company structure for parse
type Company struct {
	ID   string
	Name string
	IRI  string
}

// Category of company structure for parse
type Category struct {
	ID   string
	Name string
}

// City of products for parse
type City struct {
	ID   string
	Name string
}

// PageInstruction is a structure of settings for parse of one page crawling
type PageInstruction struct {
	ID                         string `json:"uid, omitempty"`
	Path                       string `json:"path, omitempty"`
	PageInPaginationSelector   string `json:"pageInPaginationSelector, omitempty"`
	PreviewImageOfItemSelector string `json:"previewImageOfSelector, omitempty"`
	PageParamPath              string `json:"pageParamPath, omitempty"`
	CityParamPath              string `json:"cityParamPath, omitempty"`
	ItemSelector               string `json:"itemSelector, omitempty"`
	NameOfItemSelector         string `json:"nameOfItemSelector, omitempty"`
	LinkOfItemSelector         string `json:"linkOfItemSelector, omitempty"`
	CityInCookieKey            string `json:"cityInCookieKey, omitempty"`
	CityIDForCookie            string `json:"cityIdForCookie, omitempty"`
	PriceOfItemSelector        string `json:"priceOfItemSelector, omitempty"`
}

// Product is a structure of one product from one page
type Product struct {
	Name             string
	IRI              string
	PreviewImageLink string
	Language         string
	Price            Price
	Company          Company
	Category         Category
	City             City
}

// Price structure
type Price struct {
	Value    string
	City     City
	DateTime time.Time
}
