package mpesa

const (
	// DefaultAPIHost is the default hostname for the Sandbox API.
	DefaultAPIHost = "api.sandbox.vm.co.mz"
	// DefaultServiceProviderCode is the default Service Provider Code for Sandbox.
	DefaultServiceProviderCode = "171717"
)

// Config holds the configuration parameters required by the M-Pesa API.
type Config struct {
	// APIHost is the hostname for the API (e.g., api.sandbox.vm.co.mz).
	APIHost string
	// APIKey is used for creating authorize transactions on the API.
	APIKey string
	// PublicKey is the Public Key for the M-Pesa API. Used for generating Authorization bearer tokens.
	PublicKey string
	// Origin is used for identifying the hostname which is sending transaction requests.
	Origin string
	// ServiceProviderCode is provided by Vodacom MZ.
	ServiceProviderCode string
	// InitiatorIdentifier is provided by Vodacom MZ.
	InitiatorIdentifier string
	// SecurityCredential is provided by Vodacom MZ.
	SecurityCredential string
}
