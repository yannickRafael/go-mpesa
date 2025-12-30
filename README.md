# go-mpesa

A Go library for the M-Pesa Mozambique API. 
Ported from [mpesa-mz-nodejs-lib](https://github.com/ivanruby/mpesa-mz-nodejs-lib).

## Installation

```bash
go get github.com/coffeebit/go-mpesa
```

## Usage

### Configuration

```go
import "github.com/coffeebit/go-mpesa"

config := mpesa.Config{
    PublicKey:           "<Public Key>",
    APIHost:             "api.sandbox.vm.co.mz",
    APIKey:              "<API Key>",
    Origin:              "<Origin>",
    ServiceProviderCode: "<Service Provider Code>",
    InitiatorIdentifier: "<Initiator Identifier>",
    SecurityCredential:  "<Security Credential>",
}

client := mpesa.NewClient(config)
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
