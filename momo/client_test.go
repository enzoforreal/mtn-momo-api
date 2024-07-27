package momo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Charger le fichier .env.test pour les tests unitaires
	err := godotenv.Load("../test.env")
	if err != nil {
		panic("Error loading test.env file for tests")
	}

	// Exécuter les tests unitaires
	code := m.Run()

	// Quitter avec le code de résultat des tests
	os.Exit(code)
}

func TestGetAuthToken(t *testing.T) {
	client := NewClient()
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

func TestCreateAPIUser(t *testing.T) {
	client := NewClient()
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
	client := NewClient()
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

func TestRequestToPay(t *testing.T) {
	client := NewClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	baseURL = ts.URL

	token := "test-token"
	request := RequestToPay{
		Amount:     "100",
		Currency:   "EUR",
		ExternalId: "123456",
		Payer: Payer{
			PartyIdType: "MSISDN",
			PartyId:     "46733123453",
		},
		PayerMessage: "Payment for invoice 123456",
		PayeeNote:    "Invoice 123456 payment",
	}

	referenceID, err := client.RequestToPay(token, request)
	if err != nil {
		t.Fatal(err)
	}
	if referenceID == "" {
		t.Fatalf("expected reference ID to be non-empty, got %s", referenceID)
	}
}
