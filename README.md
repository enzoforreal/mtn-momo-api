# MTN MoMo API Client for Go

This library provides a Go client for the MTN Mobile Money API. It allows developers to interact with the MTN MoMo API for operations like getting account balance, requesting payments, and more.

## Installation

To install the library, run:

```bash
go get github.com/enzoforreal/mtn-momo-api
go get github.com/gin-gonic/gin
go get github.com/google/uuid

```

## Usage

Here's an example of how to use the library:

```bash

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/enzoforreal/mtn-momo-api/momo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientConfig struct {
	SubscriptionKey string
	Environment     string
	ApiKey          string
	ApiUserID       string
	AccessToken     string
	TokenType       string
	ExpiresIn       int
}

func NewClient(config ClientConfig) *momo.Client {
	return momo.NewClient(config.SubscriptionKey, config.ApiKey, config.ApiUserID, config.Environment)
}

func main() {
	clientConfig1 := ClientConfig{
		SubscriptionKey: "0285a68a2e9542ae8fb41d6512172362", // Replace with your subscription key
		Environment:     "sandbox",
		ApiKey:          "b5f50a3e93b64ad4bca4793d4531cc29", // Replace with your API key
		ApiUserID:       "46680c23-5cb8-4f6e-8f75-aecaa6a7d415", // Replace with your API user ID
	}

	router := gin.Default()

	router.POST("/create-api-user", func(c *gin.Context) {
		client := NewClient(clientConfig1)
		var req struct {
			ReferenceID  string `json:"reference_id"`
			CallbackHost string `json:"callback_host"`
		}
		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.ReferenceID == "" {
			req.ReferenceID = uuid.New().String()
		}

		log.Printf("Creating API user with reference ID %s and callback host %s", req.ReferenceID, req.CallbackHost)
		err := client.CreateAPIUser(req.ReferenceID, req.CallbackHost)
		if err != nil {
			log.Printf("Error creating API user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("API user created successfully")
		c.JSON(http.StatusCreated, gin.H{"message": "API user created successfully", "reference_id": req.ReferenceID})
	})

	router.POST("/create-api-key", func(c *gin.Context) {
		client := NewClient(clientConfig1)
		var req struct {
			ReferenceID string `json:"reference_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.ReferenceID == "" {
			log.Printf("Reference ID is missing")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reference ID is required"})
			return
		}

		log.Printf("Creating API key for reference ID %s", req.ReferenceID)
		apiKey, err := client.CreateAPIKey(req.ReferenceID)
		if err != nil {
			log.Printf("Error creating API key: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("API key created successfully")
		c.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
	})

	router.POST("/get-auth-token", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
			return
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			log.Println("Failed to decode authorization header:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header"})
			return
		}

		authParts := strings.SplitN(string(decodedAuth), ":", 2)
		if len(authParts) != 2 {
			log.Println("Invalid authorization format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
			return
		}

		apiUserID := authParts[0]
		apiKey := authParts[1]
		log.Printf("Received API User ID: %s, API Key: %s\n", apiUserID, apiKey)

		client := NewClient(clientConfig1)
		authToken, err := client.GetAuthToken()
		if err != nil {
			log.Printf("Error getting auth token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Token retrieved successfully")
		c.JSON(http.StatusOK, gin.H{"token": authToken.AccessToken, "expires_in": authToken.ExpiresIn})
	})

	router.Run(":8080")
}

```

## Endpoint documentation 

1. Create API User


```bash

curl -X POST http://localhost:8080/create-api-user \
    -H "X-Reference-Id: your-X-Reference-id" \
    -H "Content-Type: application/json" \
    -H "Ocp-Apim-Subscription-Key: Your-Subscription-Key" \
    -d '{
          "callback_host": "https://your_callback_host.ngrok-free.app"
        }'

```

```bash

curl -X POST http://localhost:8080/create-api-key \
    -H "Content-Type: application/json" \
    -d '{
          "reference_id": "your-reference-id"
        }'

```

```bash

curl -X POST "http://localhost:8080/get-auth-token" \
    -H "Authorization: Basic your-base64-encoded-auth" \
    -H "Content-Type: application/json" \
    -H "Ocp-Apim-Subscription-Key: Your-Subscription-Key"
```


```bash
   curl -X GET http://localhost:8080/get-account-balance \
    -H "Content-Type: application/json" \
    -H "Authorization: your-access-token" \
	-H "X-Target-Environment: X_TARGET_ENVIRONMENT" \ (replace with sandbox or mtncotedivoire or mtncongo ...)
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: Your-Subscription-Key" \
```

```bash
curl -X POST http://localhost:8080/request-to-pay \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer your-access-token" \
    -d '{
          "amount": "100",
          "currency": "EUR",
          "external_id": "7890",
          "payer": {
              "party_id_type": "MSISDN",
              "party_id": "1234567890"
          },
          "payer_message": "Payment for services",
          "payee_note": "Thank you for your service"
        }'
```

## Explanations

	Initialize the client: The NewClient function creates a new client with your API key and target environment.
	Get an authentication token: The GetAuthToken function retrieves an authentication token that is required for API calls.
	Show token: To confirm that the token has been successfully obtained.
	Get the balance from your account: The GetAccountBalance feature retrieves the balance from your account.
	Show balance: To see the available balance.
	Create a payment request: A sample payment request is created with the necessary details.
	Send payment request: The RequestToPay function sends the payment request and retrieves the result.
	Show payment status: To see the status of the payment request.

This main.go file can be used as a practical example of library usage, showing how to authenticate, check the balance, and request payment. You can customize the details (such as your-api-key) and payment information to suit your needs.

## Testing

To run the tests, use the go test command:

```bash

go test -v -cover ./...

```

## Useful Links

Official documentation of the MTN MoMo API

License

This project is licensed under the MIT License.

## Project Structure


```go

mtn-momo-api/
├── LICENSE
├── README.md
├── example
│   └── main.go
├── go.mod
├── go.sum
├── momo
│   ├── client.go
│   ├── client_test.go
│   ├── errors.go
│   └── models.go
└── scripts
    ├── create-api-key.sh
    ├── create-api-user.sh
    ├── test-auth-token.sh
    └── test.sh


```

## Summary

	Installation: Clear instructions on how to install the library.
	Usage: A concrete example of using the library, complete with explanations for each step.
	Endpoint Documentation: curl commands to use during your tests.
	Explanations: Detailed descriptions of the functionality.
	Testing: Command to run unit tests.
	License: Project license type.
	Project Structure: Overview of the structure of the project files and directories.
