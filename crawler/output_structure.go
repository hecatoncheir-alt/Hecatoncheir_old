package crawler

import "time"

// StructureOfParserResponse is a structure of input data for parse products from company category pages
type StructureOfParserResponse struct {
	Language string
	Product  Product
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
}

// Price structure
type Price struct {
	Value    string
	City     string
	DateTime time.Time
}
