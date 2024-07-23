package momo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthToken(t *testing.T) {
	client := NewClient("subscriptionKey ", "apiKey", "apiUserID  ", "sandbox")
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
	client := NewClient("subscriptionKey ", "apiKey", "apiUserID  ", "sandbox")
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
	client := NewClient("subscriptionKey ", "apiKey", "apiUserID  ", "sandbox")
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
