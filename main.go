package main

import (
	"fmt"
	"log"

	"github.com/enzoforreal/mtn-momo-api/momo"
)

func main() {
	client := momo.NewClient("your-api-key", "sandbox")

	// Create API User
	referenceID := "your-reference-id"
	err := client.CreateAPIUser(referenceID, "your-callback-host")
	if err != nil {
		log.Fatalf("Error creating API user: %v", err)
	}
	fmt.Println("API user created successfully")

	// Create API Key
	apiKey, err := client.CreateAPIKey(referenceID)
	if err != nil {
		log.Fatalf("Error creating API key: %v", err)
	}
	fmt.Printf("API key created successfully: %s\n", apiKey)

	// Get API User Details
	userDetails, err := client.GetAPIUserDetails(referenceID)
	if err != nil {
		log.Fatalf("Error getting API user details: %v", err)
	}
	fmt.Printf("API user details: %v\n", userDetails)

	// Authenticate and get access token
	token, err := client.GetAuthToken()
	if err != nil {
		log.Fatalf("Error obtaining auth token: %v", err)
	}
	fmt.Printf("Authentication token: %s\n", token.AccessToken)

	// Get account balance
	balance, err := client.GetAccountBalance(token.AccessToken)
	if err != nil {
		log.Fatalf("Error getting account balance: %v", err)
	}
	fmt.Printf("Available balance: %s %s\n", balance.AvailableBalance, balance.Currency)

	// Create a payment request
	request := momo.RequestToPay{
		Amount:     "100",
		Currency:   "USD",
		ExternalId: "7890",
		Payer: momo.Payer{
			PartyIdType: "MSISDN",
			PartyId:     "1234567890",
		},
		PayerMessage: "Payment for services",
		PayeeNote:    "Thank you for your service",
	}

	// Send the payment request
	result, err := client.RequestToPay(token.AccessToken, request)
	if err != nil {
		log.Fatalf("Error requesting payment: %v", err)
	}
	fmt.Printf("Payment status: %s\n", result.Status)
}
