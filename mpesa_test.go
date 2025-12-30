package mpesa

import "testing"

func TestIsValidMSISDN(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"841234567", "258841234567", true},
		{"851234567", "258851234567", true},
		{"258841234567", "258841234567", true},
		{"258851234567", "258851234567", true},
		{"123456789", "", false},
		{"821234567", "", false}, // Only 84/85 handling was in basic node lib, though 82/83/86/87 exist in MZ, we clone the lib behavior
		{"invalid", "", false},
	}

	for _, tt := range tests {
		res, valid := IsValidMSISDN(tt.input)
		if valid != tt.valid {
			t.Errorf("IsValidMSISDN(%s) validity = %v; want %v", tt.input, valid, tt.valid)
		}
		if res != tt.expected {
			t.Errorf("IsValidMSISDN(%s) result = %s; want %s", tt.input, res, tt.expected)
		}
	}
}

func TestValidateAmount(t *testing.T) {
	if !ValidateAmount(10.0) {
		t.Error("10.0 should be valid")
	}
	if ValidateAmount(0) {
		t.Error("0 should be invalid")
	}
	if ValidateAmount(-5.0) {
		t.Error("-5.0 should be invalid")
	}
}
