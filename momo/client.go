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
	url := fmt.Sprintf("%s/provisioning/v1_0/apiuser", baseURL)
	reqBody, err := json.Marshal(map[string]string{"providerCallbackHost": callbackHost})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s with reference ID %s and callback host %s", url, referenceID, callbackHost)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create API user, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) CreateAPIKey(referenceID string) (string, error) {
	url := fmt.Sprintf("%s/provisioning/v1_0/apiuser/%s/apikey", baseURL, referenceID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Making request to %s to create API key for reference ID %s", url, referenceID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create API key, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var result struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

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

	body, _ := io.ReadAll(resp.Body)
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

	body, _ := io.ReadAll(resp.Body)
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

	body, _ := io.ReadAll(resp.Body)
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

	body, _ := io.ReadAll(resp.Body)
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
