package crawler

import "testing"

func TestUnitCitiesCanGetCityByCode(test *testing.T) {
	cities := Cities{
		"Москва":  "18414",
		"Алексин": "1688",
	}

	_, err := cities.SearchCityByCode("1")
	if err != ErrCodeDoesNotExist {
		test.Fail()
	}

	city, err := cities.SearchCityByCode("18414")
	if err != nil {
		test.Fail()
	}

	if city != "Москва" {
		test.Fail()
	}
}

func TestUnitCitiesCanGetCodeByCityName(test *testing.T) {
	cities := Cities{
		"Москва":  "18414",
		"Алексин": "1688",
	}

	_, err := cities.SearchCodeByCityName("Moscow")
	if err != ErrCityDoesNotExist {
		test.Fail()
	}

	city, err := cities.SearchCodeByCityName("Москва")
	if err != nil {
		test.Fail()
	}

	if city != "18414" {
		test.Fail()
	}
}
