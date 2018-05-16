package configuration

import (
	"github.com/hecatoncheir/Loguna/configuration"
	"os"
	"testing"
)

func TestGetConfiguration(test *testing.T) {
	defaultValues := configuration.New()
	if defaultValues.Production.Broker.Host != "192.168.99.100" {
		test.Fail()
	}

	os.Setenv("Production-Broker-Host", "localhost")

	notDefaultValues := configuration.New()
	if notDefaultValues.Production.Broker.Host != "localhost" {
		test.Fail()
	}
}
