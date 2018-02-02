package crawler

// StructureOfRequestToParse is a structure of input data for parse products from company category pages
type StructureOfRequestToParse struct {
	Language string
	Price    Price
	Company  Company
	Category Category
	Page     PageInstruction
}

// PageInstruction is a structure of settings for parse of one page crawling
type PageInstruction struct {
	ID                         string `json:"uid, omitempty"`
	Path                       string `json:"path, omitempty"`
	PageInPaginationSelector   string `json:"pageInPaginationSelector, omitempty"`
	PreviewImageOfItemSelector string `json:"previewImageOfSelector, omitempty"`
	PageParamPath              string `json:"pageParamPath, omitempty"`
	CityParamPath              string `json:"cityParamPath, omitempty"`
	CityParam                  string `json:"cityParam, omitempty"`
	ItemSelector               string `json:"itemSelector, omitempty"`
	NameOfItemSelector         string `json:"nameOfItemSelector, omitempty"`
	CityInCookieKey            string `json:"cityInCookieKey, omitempty"`
	CityIDForCookie            string `json:"cityIdForCookie, omitempty"`
	PriceOfItemSelector        string `json:"priceOfItemSelector, omitempty"`
}
