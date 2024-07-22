package momo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	baseURL = "https://sandbox.momodeveloper.mtn.com"
)

type Client struct {
	ApiKey      string
	Environment string
}

func NewClient(apiKey, environment string) *Client {
	return &Client{
		ApiKey:      apiKey,
		Environment: environment,
	}
}

func (c *Client) CreateAPIUser(referenceID, callbackHost string) error {
	url := fmt.Sprintf("%s/v1_0/apiuser", baseURL)
	if callbackHost == "" {
		callbackHost = "string" // Valeur par défaut
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
	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
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

	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
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

func (c *Client) GetAPIUserDetails(referenceID string) (map[string]string, error) {
	url := fmt.Sprintf("%s/provisioning/v1_0/apiuser/%s", baseURL, referenceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s to get API user details for reference ID %s", url, referenceID)
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
		return nil, fmt.Errorf("failed to get API user details, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetAuthToken() (*AuthToken, error) {
	url := fmt.Sprintf("%s/collection/token/", baseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
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
	if err := json.NewDecoder(resp.Body).Decode(&authToken); err != nil {
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
	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)

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
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

func (c *Client) RequestToPay(token string, request RequestToPay) (*RequestToPayResult, error) {
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", baseURL)
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Target-Environment", c.Environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s to request to pay", url)
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

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("failed to request to pay, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var result RequestToPayResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
