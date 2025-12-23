package configs

// Keycloak holds configuration for validating Keycloak-issued JWT access tokens.
// These values should match your Keycloak realm and client configuration.
type Keycloak struct {
	// IssuerURL is typically: https://<keycloak-host>/realms/<realm-name>
	IssuerURL string `mapstructure:"issuer_url"`

	// Audience is usually the client-id of the application that this service trusts.
	Audience string `mapstructure:"audience"`

	// PublicKeyPEM is the RSA public key for the realm in PEM format.
	// You can copy the realm public key from the Keycloak admin console and
	// convert it to PEM, then paste it into config.yaml as a multi-line string.
	PublicKeyPEM string `mapstructure:"public_key_pem"`

	// SkipTLSVerify disables TLS certificate verification when fetching JWKS.
	// Only use this for development/testing with self-signed certificates.
	// Defaults to false.
	SkipTLSVerify bool `mapstructure:"skip_tls_verify"`
}
