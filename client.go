package mpesa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client is the main entry point for interacting with the M-Pesa API.
type Client struct {
	config Config
	client *http.Client
}

// NewClient creates a new M-Pesa Client.
func NewClient(config Config) (*Client, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	return &Client{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}, nil
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

// B2BRequest represents the data required for a B2B transaction.
type B2BRequest struct {
	Amount              float64
	PrimaryPartyCode    string
	ReceiverPartyCode   string
	Reference           string
	ThirdPartyReference string
}

// B2B Initiates a B2B (Business-to-Business) transaction.
func (c *Client) B2B(req B2BRequest) (APIResponse, error) {
	if !ValidateAmount(req.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}
	if req.PrimaryPartyCode == "" || req.ReceiverPartyCode == "" {
		return nil, fmt.Errorf("missing primary_party_code or receiver_party_code")
	}
	if req.Reference == "" || req.ThirdPartyReference == "" {
		return nil, fmt.Errorf("missing reference or third_party_reference")
	}

	url := fmt.Sprintf("https://%s:18349/ipg/v1x/b2bPayment/", c.config.APIHost)

	payload := map[string]string{
		"input_Amount":               fmt.Sprintf("%.2f", req.Amount),
		"input_PrimaryPartyCode":     req.PrimaryPartyCode,
		"input_ReceiverPartyCode":    req.ReceiverPartyCode,
		"input_TransactionReference": req.Reference,
		"input_ThirdPartyReference":  req.ThirdPartyReference,
	}

	return c.makeRequest("POST", url, payload)
}

// QueryCustomerNameRequest represents the data required to query a customer's masked name.
type QueryCustomerNameRequest struct {
	CustomerMSISDN      string
	ThirdPartyReference string
}

// QueryCustomerName retrieves the masked name of a customer.
func (c *Client) QueryCustomerName(req QueryCustomerNameRequest) (APIResponse, error) {
	if req.CustomerMSISDN == "" || req.ThirdPartyReference == "" {
		return nil, fmt.Errorf("missing customer_msisdn or third_party_reference")
	}

	validMSISDN, ok := IsValidMSISDN(req.CustomerMSISDN)
	if !ok {
		return nil, fmt.Errorf("invalid MSISDN")
	}

	baseURL := fmt.Sprintf("https://%s:19323/ipg/v1x/queryCustomerName/", c.config.APIHost)

	url := fmt.Sprintf("%s?input_CustomerMSISDN=%s&input_ThirdPartyReference=%s&input_ServiceProviderCode=%s",
		baseURL,
		validMSISDN,
		req.ThirdPartyReference,
		c.config.ServiceProviderCode,
	)

	return c.makeRequest("GET", url, nil)
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
	// Either remove UA or set a neutral one:
	req.Header.Set("User-Agent", "mpesa-go-client/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || contentType == "" || contentType == "text/html" {
		return nil, fmt.Errorf(
			"unexpected status or content-type: status=%d, content-type=%q, body=%s",
			resp.StatusCode, contentType, string(respBytes),
		)
	}

	var result APIResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v | Body: %s", err, string(respBytes))
	}

	c.enrichResponseWithDescription(result)

	return result, nil
}

// enrichResponseWithDescription adds a readable description to the response if a known code is present.
func (c *Client) enrichResponseWithDescription(response APIResponse) {
	if response == nil {
		return
	}

	// helper to safely get string from map
	getString := func(key string) string {
		if val, ok := response[key]; ok {
			if strVal, ok := val.(string); ok {
				return strVal
			}
		}
		return ""
	}

	code := getString("output_ResponseCode")
	if code == "" {
		return
	}

	if desc, exists := ResponseCodeDescriptions[code]; exists {
		// We overwrite or add the description with our known good description
		response["output_ResponseDesc"] = desc
	}
}
