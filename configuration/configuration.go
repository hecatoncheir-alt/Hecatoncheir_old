package configuration

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	ApiVersion string

	Proxy string

	Production struct {
		ParserChannel string
		Broker        struct {
			Host string
			Port int
		}
	}

	Development struct {
		ParserChannel string
		Broker        struct {
			Host string
			Port int
		}
	}
}

func GetConfiguration() (Configuration, error) {
	configuration := Configuration{}

	apiVersion := os.Getenv("Api-Version")
	if apiVersion == "" {
		configuration.ApiVersion = "v1"
	} else {
		configuration.ApiVersion = apiVersion
	}

	productionParserChannel := os.Getenv("Production-Parser-Channel")
	if productionParserChannel == "" {
		configuration.Production.ParserChannel = "Parser"
	} else {
		configuration.Production.ParserChannel = productionParserChannel
	}

	developmentParserChannel := os.Getenv("Development-Parser-Channel")
	if developmentParserChannel == "" {
		configuration.Development.ParserChannel = "Parser"
	} else {
		configuration.Development.ParserChannel = developmentParserChannel
	}

	productionBrokerHostFromEnvironment := os.Getenv("Production-Broker-Host")
	if productionBrokerHostFromEnvironment == "" {
		configuration.Production.Broker.Host = "192.168.99.100"
	} else {
		configuration.Production.Broker.Host = productionBrokerHostFromEnvironment
	}

	productionBrokerPortFromEnvironment := os.Getenv("Production-Broker-Port")
	if productionBrokerPortFromEnvironment == "" {
		configuration.Production.Broker.Port = 4150
	} else {
		port, err := strconv.Atoi(productionBrokerPortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Production.Broker.Port = port
	}

	developmentBrokerHostFromEnvironment := os.Getenv("Development-Broker-Host")
	if developmentBrokerHostFromEnvironment == "" {
		configuration.Development.Broker.Host = "192.168.99.100"
	} else {
		configuration.Development.Broker.Host = developmentBrokerHostFromEnvironment
	}

	developmentBrokerPortFromEnvironment := os.Getenv("Development-Broker-Port")
	if developmentBrokerPortFromEnvironment == "" {
		configuration.Development.Broker.Port = 4150
	} else {
		port, err := strconv.Atoi(developmentBrokerPortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Development.Broker.Port = port
	}

	return configuration, nil
}