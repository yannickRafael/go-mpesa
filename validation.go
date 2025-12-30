package mpesa

import (
	"strconv"
	"strings"
)

// IsValidMSISDN validates a customer's MSISDN (Phone number).
// Returns the normalized MSISDN (starting with 258) and a boolean indicating validity.
func IsValidMSISDN(msisdn string) (string, bool) {
	// Is it a number? (Basic check, though detailed check happens below via parsing)
	if _, err := strconv.Atoi(msisdn); err != nil {
		return "", false
	}

	// Case 1: Length 12 and starts with 258
	if len(msisdn) == 12 && strings.HasPrefix(msisdn, "258") {
		buffer := msisdn[3:5] // Characters at index 3 and 4
		// Is it an 84 or 85 number?
		if buffer == "84" || buffer == "85" {
			return msisdn, true
		}
	} else if len(msisdn) == 9 {
		// Case 2: Length 9
		buffer := msisdn[0:2]
		// Is it an 84 or 85 number?
		if buffer == "84" || buffer == "85" {
			return "258" + msisdn, true
		}
	}

	return "", false
}

// ValidateAmount checks if the amount is valid (positive number).
func ValidateAmount(amount float64) bool {
	return amount > 0
}
