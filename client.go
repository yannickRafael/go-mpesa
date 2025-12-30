package mpesa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main entry point for interacting with the M-Pesa API.
type Client struct {
	config Config
	client *http.Client
}

// NewClient creates a new M-Pesa Client.
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// C2BRequest represents the data required for a C2B transaction.
type C2BRequest struct {
	Amount              float64
	MSISDN              string
	Reference           string
	ThirdPartyReference string
}

// QueryRequest represents the data required for a status query.
type QueryRequest struct {
	QueryReference      string
	ThirdPartyReference string
}

// ReversalRequest represents the data required for a transaction reversal.
type ReversalRequest struct {
	Amount              float64
	TransactionID       string
	ThirdPartyReference string
}

// B2CRequest represents the data required for a B2C transaction.
type B2CRequest struct {
	Amount              float64
	MSISDN              string
	Reference           string
	ThirdPartyReference string
}

// APIResponse is a generic map for the response JSON.
// In the future, this could be strongly typed based on specific response schemas.
type APIResponse map[string]interface{}

// C2B Initiates a C2B (Customer-to-Business) transaction.
func (c *Client) C2B(req C2BRequest) (APIResponse, error) {
	validMSISDN, ok := IsValidMSISDN(req.MSISDN)
	if !ok {
		return nil, fmt.Errorf("invalid MSISDN")
	}
	if !ValidateAmount(req.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}
	if req.Reference == "" || req.ThirdPartyReference == "" {
		return nil, fmt.Errorf("missing reference or third_party_reference")
	}

	url := fmt.Sprintf("https://%s:18352/ipg/v1x/c2bPayment/singleStage/", c.config.APIHost)

	payload := map[string]string{
		"input_ServiceProviderCode":  c.config.ServiceProviderCode,
		"input_CustomerMSISDN":       validMSISDN,
		"input_Amount":               fmt.Sprintf("%.2f", req.Amount),
		"input_TransactionReference": req.Reference,
		"input_ThirdPartyReference":  req.ThirdPartyReference,
	}

	return c.makeRequest("POST", url, payload)
}

// Query initiates a transaction status query.
func (c *Client) Query(req QueryRequest) (APIResponse, error) {
	if req.QueryReference == "" || req.ThirdPartyReference == "" {
		return nil, fmt.Errorf("missing query_reference or third_party_reference")
	}

	baseURL := fmt.Sprintf("https://%s:18353/ipg/v1x/queryTransactionStatus/", c.config.APIHost)

	// Constructing URL with query params manually to match the node lib pattern,
	// though net/url is usually better.
	url := fmt.Sprintf("%s?input_ServiceProviderCode=%s&input_QueryReference=%s&input_ThirdPartyReference=%s",
		baseURL,
		c.config.ServiceProviderCode,
		req.QueryReference,
		req.ThirdPartyReference,
	)

	return c.makeRequest("GET", url, nil)
}

// Reverse initiates a transaction reversal.
func (c *Client) Reverse(req ReversalRequest) (APIResponse, error) {
	if !ValidateAmount(req.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}
	if req.TransactionID == "" || req.ThirdPartyReference == "" {
		return nil, fmt.Errorf("missing transaction_id or third_party_reference")
	}

	url := fmt.Sprintf("https://%s:18354/ipg/v1x/reversal/", c.config.APIHost)

	payload := map[string]string{
		"input_ReversalAmount":      fmt.Sprintf("%.2f", req.Amount),
		"input_TransactionID":       req.TransactionID,
		"input_ThirdPartyReference": req.ThirdPartyReference,
		"input_ServiceProviderCode": c.config.ServiceProviderCode,
		"input_InitiatorIdentifier": c.config.InitiatorIdentifier,
		"input_SecurityCredential":  c.config.SecurityCredential,
	}

	return c.makeRequest("PUT", url, payload)
}

// B2C Initiates a B2C (Business-to-Customer) transaction.
func (c *Client) B2C(req B2CRequest) (APIResponse, error) {
	validMSISDN, ok := IsValidMSISDN(req.MSISDN)
	if !ok {
		return nil, fmt.Errorf("invalid MSISDN")
	}
	if !ValidateAmount(req.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}

	url := fmt.Sprintf("https://%s:18345/ipg/v1x/b2cPayment/", c.config.APIHost)

	payload := map[string]string{
		"input_ServiceProviderCode":  c.config.ServiceProviderCode,
		"input_CustomerMSISDN":       validMSISDN,
		"input_Amount":               fmt.Sprintf("%.2f", req.Amount),
		"input_TransactionReference": req.Reference,
		"input_ThirdPartyReference":  req.ThirdPartyReference,
	}

	return c.makeRequest("POST", url, payload)
}

func (c *Client) makeRequest(method, url string, body interface{}) (APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	bearerToken, err := GenerateBearerToken(c.config.APIKey, c.config.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate bearer token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.config.Origin)
	req.Header.Set("Authorization", bearerToken)
	// Add User-Agent to avoid WAF (Incapsula) blocking Go-http-client
	// Using a standard Chrome User-Agent to ensure maximal compatibility
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-2xx status codes (optional, but good practice, though Node lib just returns data)
	// The node lib seems to return response.data regardless of status, but rejects on axios error.
	// For now, we will parse JSON and return it.

	var result APIResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		// If response is not JSON, it might be an error string or empty
		return nil, fmt.Errorf("failed to decode response: %v | Body: %s", err, string(respBytes))
	}

	return result, nil
}
