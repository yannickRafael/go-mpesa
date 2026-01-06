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
