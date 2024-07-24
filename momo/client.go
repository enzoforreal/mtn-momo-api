package momo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	baseURL = "https://sandbox.momodeveloper.mtn.com"
)

func NewClient(subscriptionKey, apiKey, apiUserID, environment string) *Client {
	return &Client{
		ApiKey:          apiKey,
		ApiUserID:       apiUserID,
		SubscriptionKey: subscriptionKey,
		Environment:     environment,
	}
}

func (c *Client) CreateAPIUser(referenceID, callbackHost string) error {
	url := fmt.Sprintf("%s/v1_0/apiuser", baseURL)
	if callbackHost == "" {
		callbackHost = "string"
	}
	reqBody, err := json.Marshal(map[string]string{"providerCallbackHost": callbackHost})
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	log.Printf("Making request to %s with reference ID %s and callback host %s", url, referenceID, callbackHost)
	log.Printf("Request headers: %v", req.Header)
	log.Printf("Request body: %s", reqBody)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated {
		errMsg := fmt.Sprintf("failed to create API user, status code: %d, response: %s", resp.StatusCode, string(body))
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	log.Println("API user created successfully")
	return nil
}

func (c *Client) CreateAPIKey(referenceID string) (string, error) {
	url := fmt.Sprintf("%s/v1_0/apiuser/%s/apikey", baseURL, referenceID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s to create API key for reference ID %s", url, referenceID)
	log.Printf("Request headers: %v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated {
		errMsg := fmt.Sprintf("failed to create API key, status code: %d, response: %s", resp.StatusCode, string(body))
		log.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	var result struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return "", err
	}

	log.Println("API key created successfully")
	return result.APIKey, nil
}

func (c *Client) GetAuthToken() (*AuthToken, error) {
	url := fmt.Sprintf("%s/collection/token/", baseURL)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.ApiUserID, c.ApiKey)))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Ocp-Apim-Subscription-Key", c.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s to get auth token", url)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get auth token, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var authToken AuthToken
	if err := json.Unmarshal(body, &authToken); err != nil {
		return nil, err
	}

	return &authToken, nil
}

func (c *Client) GetAccountBalance(token string) (*Balance, error) {
	url := fmt.Sprintf("%s/collection/v1_0/account/balance", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Target-Environment", c.Environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	log.Printf("Making request to %s to get account balance", url)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get account balance, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var balance Balance
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&balance); err != nil {
		return nil, err
	}

	return &balance, nil
}
