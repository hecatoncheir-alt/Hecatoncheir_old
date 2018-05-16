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
		HecatoncheirTopic string
		LogunaTopic       string
		Broker            struct {
			Host string
			Port int
		}
	}

	Development struct {
		HecatoncheirTopic string
		LogunaTopic       string
		Broker            struct {
			Host string
			Port int
		}
	}
}

// GetConfiguration function check environment variables and return structure with value
func New() *Configuration {
	configuration := Configuration{}

	apiVersion := os.Getenv("Api-Version")
	if apiVersion == "" {
		configuration.APIVersion = "1.0.0"
	} else {
		configuration.APIVersion = apiVersion
	}

	productionLogunaTopic := os.Getenv("Production-Loguna-Topic")
	if productionLogunaTopic == "" {
		configuration.Production.LogunaTopic = "Loguna"
	} else {
		configuration.Production.LogunaTopic = productionLogunaTopic
	}

	developmentLogunaTopic := os.Getenv("Development-Loguna-Topic")
	if developmentLogunaTopic == "" {
		configuration.Development.LogunaTopic = "DevLoguna"
	} else {
		configuration.Development.LogunaTopic = developmentLogunaTopic
	}

	productionParserHecatoncheirTopic := os.Getenv("Production-Hecatoncheir-Topic")
	if productionParserHecatoncheirTopic == "" {
		configuration.Production.HecatoncheirTopic = "Hecatoncheir"
	} else {
		configuration.Production.HecatoncheirTopic = productionParserHecatoncheirTopic
	}

	developmentParserHecatoncheirTopic := os.Getenv("Development-Hecatoncheir-Topic")
	if developmentParserHecatoncheirTopic == "" {
		configuration.Development.HecatoncheirTopic = "DevHecatoncheir"
	} else {
		configuration.Development.HecatoncheirTopic = developmentParserHecatoncheirTopic
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

	return &configuration
}
