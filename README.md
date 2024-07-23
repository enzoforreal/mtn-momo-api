# MTN MoMo API Client for Go

This library provides a Go client for the MTN Mobile Money API. It allows developers to interact with the MTN MoMo API for operations like getting account balance, requesting payments, and more.

## Installation

To install the library, run:

```bash
go get github.com/enzoforreal/mtn-momo-api


Usage
Here's an example of how to use the library:

package main

import (
	"log"
	"net/http"

	"github.com/enzoforreal/mtn-momo-api/momo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientConfig struct {
	SubscriptionKey string
	Environment     string
}

func NewClient(config ClientConfig) *momo.Client {
	return momo.NewClient(config.SubscriptionKey, config.Environment)
}

func main() {
	clientConfig1 := ClientConfig{
		SubscriptionKey: "0285a68a2e9542ae8fb41d6512172362", // Remplacez par votre clé API
		Environment:     "sandbox",
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
		client := NewClient(clientConfig1)
		token, err := client.GetAuthToken()
		if err != nil {
			log.Printf("Error obtaining auth token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Auth token obtained successfully")
		c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken, "expires_in": token.ExpiresIn})
	})

	router.Run(":8080")
}


```

##  Endpoint documentation

1. Create API User

```bash

curl -X POST http://localhost:8080/create-api-user \
    -H "X-Reference-Id: your-X-Reference-id" \
    -H "Content-Type: application/json" \
    -H "Ocp-Apim-Subscription-Key: Your-Subscription-Key" \
    -d '{
          "callback_host": "https://your_callback_host.ngrok-free.app"
        }'


 
2. Create API key 



curl -X POST http://localhost:8080/create-api-key \
    -H "Content-Type: application/json" \
    -d '{
          "reference_id": "your-reference-id"
        }'



3. Get an authentication token

curl -X POST "http://localhost:8080/get-auth-token" \
    -H "Authorization: ZThjNWEzNjQtNzAxOC00YmNmLWI2NWQtZGViYTBmYjk4MTdhOjIwZjViYjA3ODBkYTQwMTg5MDc0YjBkMTVhNzhkYTAw" \  (replace Basic authentication header containing API user ID and API key. Should be sent in as B64 encoded)
    -H "Content-Type: application/json" \
    -H "Ocp-Apim-Subscription-Key: 0285a68a2e9542ae8fb41d6512172362" (replace with your Subscription-Key)



4. Get API User Details

curl -X GET http://localhost:8080/get-api-user-details/votre-reference-id \
    -H "Content-Type: application/json"

5. Get the balance from the account

curl -X GET http://localhost:8080/get-account-balance \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer votre-access-token"


6. Send a payment request (Request to Pay)


curl -X POST http://localhost:8080/request-to-pay \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer votre-access-token" \
    -d '{
          "amount": "100",
          "currency": "EUR",
          "external_id": "7890",
          "payer": {
              "party_id_type": "MSISDN",
              "party_id": "1234567890"
        },
		
		{
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
Send payment request: RequestToPay function sends the payment request and retrieves the result.
Show payment status: To see the status of the payment request.

This main.go file can be used as a practical example of library usage, showing how to authenticate, check the balance and request payment. You can customize the details (such as your-api-key) and payment information to suit your needs.

Testing
To run the tests, use the go test command:

go test -v -cover ./...



## Liens utiles

[Official documentation of the MTN MoMo API](https://momodeveloper.mtn.com/)


License
This project is licensed under the MIT License.

Project Structure

mtn-momo-api/
├── go.mod
├── go.sum
├── README.md
├── main.go
├── momo/
│   ├── client.go
│   ├── models.go
│   ├── errors.go
│   ├── client_test.go
└── .gitignore




- **Installation**: Clear instructions on how to install the library.
- **Usage**: A concrete example of using the library, complete with explanations in French for each step.
- **Testing**: Command to run unit tests.
- **License**: Project license type.
- **Project Structure**: Overview of the structure of the project files and directories.