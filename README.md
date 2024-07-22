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

func main() {
	client := momo.NewClient("Subscription Key", "sandbox")

	router := gin.Default()

	router.POST("/create-api-user", func(c *gin.Context) {
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

	router.Run(":8080")
}




Explanation
Initialize the client: The NewClient function creates a new client with your API key and target environment.
Obtain an authentication token: The GetAuthToken function retrieves an authentication token necessary for API calls.
Display the token: To confirm that the token has been obtained successfully.
Get account balance: The GetAccountBalance function retrieves the balance of your account.
Display the balance: To see the available balance.
Create a payment request: An example payment request is created with the necessary details.
Send the payment request: The RequestToPay function sends the payment request and retrieves the result.
Display the payment status: To see the status of the payment request.
This main.go file can be used as a practical example of using the library, showing how to authenticate, check balance, and request payment. You can customize the details (such as your-api-key) and payment information as needed.

Testing
To run the tests, use the go test command:

go test -v -cover ./...



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