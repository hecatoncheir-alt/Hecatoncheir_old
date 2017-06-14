package broker

import "testing"

func TestBrokerCanConnectToNSQ(test *testing.T) {
	broker := New()
	items, err := broker.ListenTopic("Hecatoncheir", "Parsed items")
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		if item == "" {
			test.Fail()
		}
	}
}
