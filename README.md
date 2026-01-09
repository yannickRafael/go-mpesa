# go-mpesa

A Go library for the M-Pesa Mozambique API. 
Ported from [mpesa-mz-nodejs-lib](https://github.com/ivanruby/mpesa-mz-nodejs-lib).

## Installation

```bash
go get github.com/yannickRafael/go-mpesa
```

## Usage

### Configuration

```go
import "github.com/yannickRafael/go-mpesa"

config := mpesa.Config{
    PublicKey:           "<Public Key>", // Note: The key must NOT be in PEM format
    APIHost:             "api.sandbox.vm.co.mz",
    APIKey:              "<API Key>",
    Origin:              "<Origin>",
    ServiceProviderCode: "<Service Provider Code>",
    InitiatorIdentifier: "<Initiator Identifier>",
    SecurityCredential:  "<Security Credential>",
}

client, err := mpesa.NewClient(config)
if err != nil {
    panic(err)
}
```

### Customer to Business (C2B)

```go
response, err := client.C2B(mpesa.C2BRequest{
    Amount:              100.0,
    MSISDN:              "841234567",
    Reference:           "REF123",
    ThirdPartyReference: "3RD123",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)
```

### Query Status

```go
response, err := client.Query(mpesa.QueryRequest{
    QueryReference:      "TRANS_ID_OR_CONV_ID",
    ThirdPartyReference: "3RD123",
})
```

### Reversal

```go
response, err := client.Reverse(mpesa.ReversalRequest{
    Amount:              100.0,
    TransactionID:       "TRANS_ID",
    ThirdPartyReference: "3RD123",
})
```

### Business to Customer (B2C)

```go
response, err := client.B2C(mpesa.B2CRequest{
    Amount:              100.0,
    MSISDN:              "841234567",
    Reference:           "REF123",
    ThirdPartyReference: "3RD123",
})
```

## Features

- RSA Encryption for Authentication (Bearer Token)
- MSISDN Validation (Mozambique format)
- Parameter handling for standard M-Pesa operations

## Example

You can find a complete runnable example in `test/test.go`.

### Prerequisites

Create a `.env` file in your project root with your credentials:

```bash
MPESA_API_KEY=your_api_key_here
MPESA_PUBLIC_KEY=your_public_key_here
MPESA_API_HOST=api.sandbox.vm.co.mz
MPESA_ORIGIN=developer.mpesa.vm.co.mz
MPESA_SERVICE_PROVIDER_CODE=171717
MPESA_INITIATOR_IDENTIFIER=your_initiator_identifier
MPESA_SECURITY_CREDENTIAL=your_security_credential
```

### Run the Example

To run the example code:

```bash
# Install dependencies
go get github.com/joho/godotenv

# Run the test file
go run test/test.go
```

### Example Code

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	mpesa "github.com/yannickRafael/go-mpesa"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Example Configuration using Sandbox details
	// We load from environment variables.
	config := mpesa.Config{
		APIHost:             getEnv("MPESA_API_HOST", mpesa.DefaultAPIHost),
		APIKey:              os.Getenv("MPESA_API_KEY"),
		PublicKey:           os.Getenv("MPESA_PUBLIC_KEY"),
		Origin:              getEnv("MPESA_ORIGIN", "developer.mpesa.vm.co.mz"),
		ServiceProviderCode: getEnv("MPESA_SERVICE_PROVIDER_CODE", mpesa.DefaultServiceProviderCode),
		InitiatorIdentifier: getEnv("MPESA_INITIATOR_IDENTIFIER", "DEFAULT_INITIATOR"),
		SecurityCredential:  getEnv("MPESA_SECURITY_CREDENTIAL", "DEFAULT_SEC_CRED"),
	}

	if config.APIKey == "" || config.PublicKey == "" {
		fmt.Println("Warning: MPESA_API_KEY and MPESA_PUBLIC_KEY env vars are empty. Usage might fail if real keys are needed.")
		// For demo purposes, we might want to put dummy values to see the request form,
		// but typically we should warn.
	}

	client, err := mpesa.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Initiating C2B Transaction...")

	// 1. C2B Transaction Example
	c2bRequest := mpesa.C2BRequest{
		Amount:              1.0,
		MSISDN:              "844236139", // Valid MSISDN
		Reference:           "REF12346",  // Alphanumeric only, no hyphens
		ThirdPartyReference: "3RD1239",
	}

	c2bResponse, err := client.C2B(c2bRequest)
	if err != nil {
		log.Printf("C2B Error: %v\n", err)
	} else {
		fmt.Printf("C2B Response: %+v\n", c2bResponse)
	}

	var conversationID string
	if c2bResponse != nil {
		if val, ok := c2bResponse["output_ConversationID"].(string); ok {
			conversationID = val
		}
	}

	// 2. Query Transaction Example
	fmt.Println("\nQuerying Transaction Status...")

	queryRef := "SOME_CONVERSATION_ID"
	if conversationID != "" {
		queryRef = conversationID
		fmt.Printf("Using ConversationID from C2B: %s\n", queryRef)
	}

	queryRequest := mpesa.QueryRequest{
		QueryReference:      queryRef,
		ThirdPartyReference: "3RD123",
	}

	queryResponse, err := client.Query(queryRequest)
	if err != nil {
		log.Printf("Query Error: %v\n", err)
	} else {
		fmt.Printf("Query Response: %+v\n", queryResponse)
	}
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is present, its value (which may be empty) is returned.
// Otherwise, the fallback value is returned.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
```

