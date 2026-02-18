package service

import (
	"fmt"
	"testing"
)

func TestDemoPrettyPrintNewsItems(t *testing.T) {
	mockNews := []NewsItem{
		{
			ID:          "INT-101",
			Description: "Promozione News: 10% di sconto su abbonamenti.",
			Channel:     "PUSH",
		},
		{
			ID:          "INT-102",
			Description: "Benvenuto nel nuovo sistema di notifiche Julia!",
			Channel:     "EMAIL",
		},
	}

	fmt.Println("\nDemo: Pretty Printing Mock News Items")
	for _, item := range mockNews {
		fmt.Println(item)
	}
}
