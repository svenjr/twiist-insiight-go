package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type InitiateAuthRequest struct {
	AuthFlow       string            `json:"AuthFlow"`
	ClientId       string            `json:"ClientId"`
	AuthParameters map[string]string `json:"AuthParameters"`
}

type InitiateAuthResponse struct {
	AuthenticationResult struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
		TokenType    string `json:"TokenType"`
		ExpiresIn    int    `json:"ExpiresIn"`
	} `json:"AuthenticationResult"`
}

type RefreshTokenRequest struct {
	AuthFlow       string            `json:"AuthFlow"`
	ClientId       string            `json:"ClientId"`
	AuthParameters map[string]string `json:"AuthParameters"`
}

type RefreshTokenResponse struct {
	AuthenticationResult struct {
		AccessToken string `json:"AccessToken"`
		TokenType   string `json:"TokenType"`
		ExpiresIn   int    `json:"ExpiresIn"`
	} `json:"AuthenticationResult"`
}

type PWDOverview struct {
	PWDId       string `json:"pwdId"`
	PWDNickname string `json:"pwdNickname"`
	Status      struct {
		Date    string `json:"date"`
		Summary struct {
			GlucoseDate      string `json:"glucoseDate"`
			GlucoseUnit      string `json:"glucoseUnit"`
			CGMRateArrow     string `json:"cgmRateArrow"`
			GlucoseQuantity  string `json:"glucoseQuantity"`
			PumpBatteryLevel string `json:"pumpBatteryLevel"`
		} `json:"summary"`
	} `json:"status"`
}

type InsiightClient struct {
	accessToken  string
	refreshToken string
	clientId     string
	userPoolId   string
	httpClient   *http.Client
}

func NewInsiightClient() *InsiightClient {
	return &InsiightClient{
		clientId:   "65ev2vbkr2mle7uu4cqkn7ohgl",
		userPoolId: "us-east-1_fnkWvSdfv",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *InsiightClient) Login(username, password string) error {
	authReq := InitiateAuthRequest{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: c.clientId,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	body, err := json.Marshal(authReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://cognito-idp.us-east-1.amazonaws.com/%s", c.userPoolId), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: %d %s", resp.StatusCode, string(respBody))
	}

	var authResp InitiateAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return err
	}

	c.accessToken = authResp.AuthenticationResult.AccessToken
	c.refreshToken = authResp.AuthenticationResult.RefreshToken

	fmt.Printf("Login successful! Token expires in %d seconds\n", authResp.AuthenticationResult.ExpiresIn)
	return nil
}

func (c *InsiightClient) RefreshToken() error {
	refreshReq := RefreshTokenRequest{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		ClientId: c.clientId,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": c.refreshToken,
		},
	}

	body, err := json.Marshal(refreshReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://cognito-idp.us-east-1.amazonaws.com/%s", c.userPoolId), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("token refresh failed: %d %s", resp.StatusCode, string(respBody))
	}

	var refreshResp RefreshTokenResponse
	if err := json.Unmarshal(respBody, &refreshResp); err != nil {
		return err
	}

	c.accessToken = refreshResp.AuthenticationResult.AccessToken

	fmt.Printf("Token refreshed! New token expires in %d seconds\n", refreshResp.AuthenticationResult.ExpiresIn)
	return nil
}

func (c *InsiightClient) GetPWDOverviews() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://follower-service.mytwiistportal.com/pwd/overviews", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API call failed: %d %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <username> <password>")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]

	client := NewInsiightClient()

	fmt.Println("Logging in...")
	if err := client.Login(username, password); err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Fetching PWD overviews...")
	rawData, err := client.GetPWDOverviews()
	if err != nil {
		fmt.Printf("Failed to get overviews: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Raw API Response:")
	fmt.Println("==================")
	
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawData, "", "  "); err != nil {
		fmt.Printf("Raw JSON (unparsed): %s\n", string(rawData))
	} else {
		fmt.Println(prettyJSON.String())
	}
}
