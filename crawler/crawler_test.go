package crawler

import (
	"testing"
)

func TestUnitCompanyCanBeDecoded(test *testing.T) {
	inputCompany := `{
	      "Language":"en",
	      "Company":{
	         "ID":"0x2786",
	         "Name":"Company test name",
	         "IRI":""
	      },
	      "Category":{
	         "ID":"",
	         "Name":"Test category"
	      },
	      "City":{
	         "ID":"0x2788",
	         "Name":"Test city"
	      },
	      "PageInstruction":{
	         "uid":"0x2789",
	         "path":"smartfony-i-svyaz/smartfony-205",
	         "pageInPaginationSelector":".pagination-list .pagination-item",
	         "previewImageOfSelector":"",
	         "pageParamPath":"/f/page=",
	         "cityParamPath":"?cityId=",
	         "cityParam":"CityCZ_975",
	         "itemSelector":".grid-view .product-tile",
	         "nameOfItemSelector":".product-tile-title",
			 "linkOfItemSelector":".product-tile-title a",
	         "cityInCookieKey":"",
	         "cityIdForCookie":"",
	         "priceOfItemSelector":".product-price-current"
          }
	    }`

	parser, err := NewParserInstructionsFromJSON(inputCompany)
	if err != nil {
		test.Error(err)
	}

	if parser.Language != "en" {
		test.Fail()
	}

	if parser.Company.Name != "Company test name" {
		test.Fail()
	}

	if parser.Category.Name != "Test category" {
		test.Fail()
	}

	if parser.City.Name != "Test city" {
		test.Fail()
	}

	if parser.PageInstruction.ID != "0x2789" {
		test.Fail()
	}
}
