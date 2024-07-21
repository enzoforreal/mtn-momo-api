package momo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthToken(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AuthToken{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		})
	}))
	defer ts.Close()

	baseURL = ts.URL
	token, err := client.GetAuthToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "test-access-token" {
		t.Fatalf("expected access token to be 'test-access-token', got %s", token.AccessToken)
	}
}

func TestGetAccountBalance(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Balance{
			AvailableBalance: "1000",
			Currency:         "USD",
		})
	}))
	defer ts.Close()

	baseURL = ts.URL
	balance, err := client.GetAccountBalance("test-access-token")
	if err != nil {
		t.Fatal(err)
	}
	if balance.AvailableBalance != "1000" {
		t.Fatalf("expected balance to be '1000', got %s", balance.AvailableBalance)
	}
}

func TestRequestToPay(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(RequestToPayResult{
			Amount:                 "100",
			Currency:               "USD",
			FinancialTransactionId: "123456",
			ExternalId:             "7890",
			Payer: Payer{
				PartyIdType: "MSISDN",
				PartyId:     "1234567890",
			},
			Status: "SUCCESSFUL",
		})
	}))
	defer ts.Close()

	baseURL = ts.URL
	request := RequestToPay{
		Amount:     "100",
		Currency:   "USD",
		ExternalId: "7890",
		Payer: Payer{
			PartyIdType: "MSISDN",
			PartyId:     "1234567890",
		},
		PayerMessage: "Payment",
		PayeeNote:    "Note",
	}
	result, err := client.RequestToPay("test-access-token", request)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "SUCCESSFUL" {
		t.Fatalf("expected status to be 'SUCCESSFUL', got %s", result.Status)
	}
}

func TestGetAccountBalanceError(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorReason{
			Code:    "INTERNAL_ERROR",
			Message: "An internal error occurred.",
		})
	}))
	defer ts.Close()

	baseURL = ts.URL
	_, err := client.GetAccountBalance("test-access-token")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestCreateAPIUser(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	baseURL = ts.URL
	err := client.CreateAPIUser("test-reference-id", "test-callback-host")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAPIKey(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"apiKey": "test-api-key"})
	}))
	defer ts.Close()

	baseURL = ts.URL
	apiKey, err := client.CreateAPIKey("test-reference-id")
	if err != nil {
		t.Fatal(err)
	}
	if apiKey != "test-api-key" {
		t.Fatalf("expected API key to be 'test-api-key', got %s", apiKey)
	}
}

func TestGetAPIUserDetails(t *testing.T) {
	client := NewClient("test-api-key", "sandbox")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"providerCallbackHost": "test-callback-host",
			"targetEnvironment":    "sandbox",
		})
	}))
	defer ts.Close()

	baseURL = ts.URL
	details, err := client.GetAPIUserDetails("test-reference-id")
	if err != nil {
		t.Fatal(err)
	}
	if details["providerCallbackHost"] != "test-callback-host" {
		t.Fatalf("expected callback host to be 'test-callback-host', got %s", details["providerCallbackHost"])
	}
	if details["targetEnvironment"] != "sandbox" {
		t.Fatalf("expected target environment to be 'sandbox', got %s", details["targetEnvironment"])
	}
}
