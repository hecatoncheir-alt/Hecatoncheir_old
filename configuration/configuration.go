package configuration

import (
	"log"
	"os"
	"strconv"
)

// Configuration is structure of config data from environment with default value
type Configuration struct {
	APIVersion string

	Proxy string

	Production struct {
		Channel string
		Broker  struct {
			Host string
			Port int
		}
	}

	Development struct {
		Channel string
		Broker  struct {
			Host string
			Port int
		}
	}
}

// GetConfiguration function check environment variables and return structure with value
func GetConfiguration() (Configuration, error) {
	configuration := Configuration{}

	apiVersion := os.Getenv("Api-Version")
	if apiVersion == "" {
		configuration.APIVersion = "v1"
	} else {
		configuration.APIVersion = apiVersion
	}

	productionParserChannel := os.Getenv("Production-Channel")
	if productionParserChannel == "" {
		configuration.Production.Channel = "Parser"
	} else {
		configuration.Production.Channel = productionParserChannel
	}

	developmentParserChannel := os.Getenv("Development-Channel")
	if developmentParserChannel == "" {
		configuration.Development.Channel = "Parser"
	} else {
		configuration.Development.Channel = developmentParserChannel
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
