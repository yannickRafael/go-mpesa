package mpesa

import (
	"testing"
)

func TestEnrichResponseWithDescription(t *testing.T) {
	client, _ := NewClient(Config{})

	tests := []struct {
		name          string
		inputResponse APIResponse
		expectedDesc  string
	}{
		{
			name: "Known Code INS-0",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-0",
				"output_ResponseDesc": "Old Description",
			},
			expectedDesc: "Request processed successfully",
		},
		{
			name: "Known Code INS-10",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-10",
			},
			expectedDesc: "Duplicate Transaction",
		},
		{
			name: "Reversal Error INS-2001",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-2001",
			},
			expectedDesc: "Initiator authentication error.",
		},
		{
			name: "C2B Error INS-2006",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-2006",
			},
			expectedDesc: "Insufficient balance",
		},
		{
			name: "B2C Error INS-996",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-996",
			},
			expectedDesc: "Customer Account Status Not Active",
		},
		{
			name: "Unknown Code",
			inputResponse: APIResponse{
				"output_ResponseCode": "INS-9999",
				"output_ResponseDesc": "Original Desc",
			},
			expectedDesc: "Original Desc",
		},
		{
			name: "No Code",
			inputResponse: APIResponse{
				"output_ResponseDesc": "Just Desc",
			},
			expectedDesc: "Just Desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.enrichResponseWithDescription(tt.inputResponse)

			got := ""
			if val, ok := tt.inputResponse["output_ResponseDesc"].(string); ok {
				got = val
			}

			if got != tt.expectedDesc {
				t.Errorf("enrichResponseWithDescription() = %v, want %v", got, tt.expectedDesc)
			}
		})
	}
}

func TestB2B(t *testing.T) {
	c, _ := NewClient(Config{})

	tests := []struct {
		name    string
		req     B2BRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Request (Validation Only)",
			req: B2BRequest{
				Amount:              100,
				PrimaryPartyCode:    "100000",
				ReceiverPartyCode:   "200000",
				Reference:           "REF123",
				ThirdPartyReference: "TP123",
			},
			wantErr: true, // expect true because network call will fail
		},
		{
			name: "Invalid Amount",
			req: B2BRequest{
				Amount:              0,
				PrimaryPartyCode:    "100000",
				ReceiverPartyCode:   "200000",
				Reference:           "REF123",
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "invalid amount",
		},
		{
			name: "Missing PrimaryPartyCode",
			req: B2BRequest{
				Amount:              100,
				ReceiverPartyCode:   "200000",
				Reference:           "REF123",
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "missing primary_party_code or receiver_party_code",
		},
		{
			name: "Missing ReceiverPartyCode",
			req: B2BRequest{
				Amount:              100,
				PrimaryPartyCode:    "100000",
				Reference:           "REF123",
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "missing primary_party_code or receiver_party_code",
		},
		{
			name: "Missing Reference",
			req: B2BRequest{
				Amount:              100,
				PrimaryPartyCode:    "100000",
				ReceiverPartyCode:   "200000",
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "missing reference or third_party_reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.B2B(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("B2B() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" && err != nil && err.Error() != tt.errMsg {
				t.Errorf("B2B() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}

func TestQueryCustomerName(t *testing.T) {
	c, _ := NewClient(Config{})

	tests := []struct {
		name    string
		req     QueryCustomerNameRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Request (Validation Only)",
			req: QueryCustomerNameRequest{
				CustomerMSISDN:      "258841234567",
				ThirdPartyReference: "TP123",
			},
			wantErr: true, // expect true due to network failure
		},
		{
			name: "Missing CustomerMSISDN",
			req: QueryCustomerNameRequest{
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "missing customer_msisdn or third_party_reference",
		},
		{
			name: "Missing ThirdPartyReference",
			req: QueryCustomerNameRequest{
				CustomerMSISDN: "258841234567",
			},
			wantErr: true,
			errMsg:  "missing customer_msisdn or third_party_reference",
		},
		{
			name: "Invalid MSISDN",
			req: QueryCustomerNameRequest{
				CustomerMSISDN:      "123",
				ThirdPartyReference: "TP123",
			},
			wantErr: true,
			errMsg:  "invalid MSISDN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.QueryCustomerName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryCustomerName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" && err != nil && err.Error() != tt.errMsg {
				t.Errorf("QueryCustomerName() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}
