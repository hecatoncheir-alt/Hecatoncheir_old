package crawler

import "errors"

// Cities codes for company
type Cities map[string]string

// ErrCityDoesNotExist is an error means city can't be found
var ErrCityDoesNotExist = errors.New("city does not exist")

// SearchCodeByCity method of Cities type
func (cities *Cities) SearchCodeByCityName(cityName string) (string, error) {
	for city, code := range *cities {
		if city == cityName {
			return code, nil
		}
	}

	return "", ErrCityDoesNotExist
}

// ErrCodeDoesNotExist is an error means code an't be found
var ErrCodeDoesNotExist = errors.New("code does not exist")

// SearchCityByCode method of Cities type
func (cities *Cities) SearchCityByCode(codeName string) (string, error) {
	for city, code := range *cities {
		if code == codeName {
			return city, nil
		}
	}
	return "", ErrCodeDoesNotExist
}
