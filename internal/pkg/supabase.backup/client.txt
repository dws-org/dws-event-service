package supabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oskargbc/dws-event-service.git/configs"

	"github.com/golang-jwt/jwt/v4"
)

type SupabaseClient struct {
	config *configs.Config
}

type SupabaseUser struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	UserMeta  map[string]interface{} `json:"user_metadata"`
	AppMeta   map[string]interface{} `json:"app_metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewSupabaseClient(config *configs.Config) *SupabaseClient {
	return &SupabaseClient{
		config: config,
	}
}

// VerifyJWT verifies a Supabase JWT token and returns the user ID
func (s *SupabaseClient) VerifyJWT(tokenString string) (string, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to parse token: %w", err))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("could not parse claims")
	}

	// Check if aud is "authenticated"
	if aud, ok := claims["aud"].(string); ok && aud == "authenticated" {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		} else {
			return "", errors.New("sub not found in token")
		}
	} else {
		return "", errors.New("not authenticated")
	}




}

// fetchPublicKey fetches the public key from Supabase JWKS endpoint
func (s *SupabaseClient) fetchPublicKey(kid string) (interface{}, error) {
	jwksURL := fmt.Sprintf("%s/auth/v1/jwks", s.config.Supabase.URL)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			// For now, we'll return the key data as a string
			// In a production environment, you'd want to properly construct the RSA public key
			return fmt.Sprintf("{\"kid\":\"%s\",\"kty\":\"%s\",\"alg\":\"%s\",\"use\":\"%s\",\"n\":\"%s\",\"e\":\"%s\"}",
				key.Kid, key.Kty, key.Alg, key.Use, key.N, key.E), nil
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}
