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
    "fmt"
    "log"
    "github.com/enzoforreal/mtn-momo-api/momo"
)

func main() {
    // Create a new client with your API key and target environment
    client := momo.NewClient("your-api-key", "sandbox")

    // Obtain an authentication token
    token, err := client.GetAuthToken()
    if err != nil {
        log.Fatalf("Error obtaining auth token: %v", err)
    }

    // Display the authentication token for confirmation
    fmt.Printf("Authentication token: %s\n", token.AccessToken)

    // Get the balance of your account
    balance, err := client.GetAccountBalance(token.AccessToken)
    if err != nil {
        log.Fatalf("Error getting account balance: %v", err)
    }

    // Display the available balance
    fmt.Printf("Available balance: %s %s\n", balance.AvailableBalance, balance.Currency)

    // Create a payment request
    request := momo.RequestToPay{
        Amount:    "100",
        Currency:  "USD",
        ExternalId: "7890",
        Payer: momo.Payer{
            PartyIdType: "MSISDN",
            PartyId: "1234567890",
        },
        PayerMessage: "Payment for services",
        PayeeNote:    "Thank you for your service",
    }

    // Send the payment request
    result, err := client.RequestToPay(token.AccessToken, request)
    if err != nil {
        log.Fatalf("Error requesting payment: %v", err)
    }

    // Display the payment status
    fmt.Printf("Payment status: %s\n", result.Status)
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

go test ./...


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