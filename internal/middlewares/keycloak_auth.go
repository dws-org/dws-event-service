package middlewares

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RealmAccess represents the realm_access claim from Keycloak which includes roles.
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// KeycloakClaims represents the standard JWT claims we care about from Keycloak.
// Additional fields (like realm_access / resource_access) can be added later if needed.
type KeycloakClaims struct {
	jwt.RegisteredClaims
	AZP         string      `json:"azp,omitempty"`          // Authorized party - the client ID that issued the token
	RealmAccess RealmAccess `json:"realm_access,omitempty"` // Realm roles assigned to the user
}

// jwks represents a JSON Web Key Set as returned by Keycloak's /certs endpoint.
type jwks struct {
	Keys []jwk `json:"keys"`
}

// jwk represents a single JSON Web Key (we only care about RSA keys).
type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"` // modulus
	E   string `json:"e"` // exponent
}

// fetchJWKS loads the JWKS from the given URL and returns a map of kid -> *rsa.PublicKey.
// If skipTLSVerify is true, TLS certificate verification is disabled (useful for self-signed certs).
func fetchJWKS(url string, skipTLSVerify bool) (map[string]*rsa.PublicKey, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if skipTLSVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %s", resp.Status)
	}

	var set jwks
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, k := range set.Keys {
		if strings.ToUpper(k.Kty) != "RSA" {
			continue
		}

		pub, err := jwkToPublicKey(k)
		if err != nil {
			// Skip keys we can't parse; continue with others.
			continue
		}
		keys[k.Kid] = pub
	}

	if len(keys) == 0 {
		return nil, errors.New("no RSA keys found in JWKS")
	}

	return keys, nil
}

// jwkToPublicKey converts an RSA JWK into an *rsa.PublicKey.
func jwkToPublicKey(k jwk) (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eb, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	// Convert big-endian bytes to int
	e := 0
	for _, b := range eb {
		e = e<<8 | int(b)
	}
	if e <= 0 {
		return nil, errors.New("invalid exponent")
	}

	n := new(big.Int).SetBytes(nb)
	return &rsa.PublicKey{N: n, E: e}, nil
}

// KeycloakAuthMiddleware validates a Keycloak-issued JWT access token from the Authorization header.
// On success, it stores the user ID (subject) in the request context under UserIDKey.
func KeycloakAuthMiddleware() gin.HandlerFunc {
	cfg := configs.GetEnvConfig()
	kc := cfg.Keycloak

	issuer := strings.TrimSpace(kc.IssuerURL)
	if issuer == "" {
		// Misconfiguration - fail fast rather than accepting unauthenticated traffic.
		panic("keycloak.issuer_url is not configured")
	}

	// Derive the JWKS URL from the issuer, e.g.:
	// https://<keycloak-host>/realms/<realm-name>/protocol/openid-connect/certs
	jwksURL := strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
	audience := strings.TrimSpace(kc.Audience)

	// Simple in-memory cache of JWKS keys by kid.
	var (
		mu       sync.RWMutex
		jwksKeys = make(map[string]*rsa.PublicKey)
	)

	getPublicKeyForKID := func(kid string) (*rsa.PublicKey, error) {
		mu.RLock()
		if key, ok := jwksKeys[kid]; ok {
			mu.RUnlock()
			return key, nil
		}
		mu.RUnlock()

		keys, err := fetchJWKS(jwksURL, kc.SkipTLSVerify)
		if err != nil {
			return nil, err
		}

		mu.Lock()
		for k, v := range keys {
			jwksKeys[k] = v
		}
		key := jwksKeys[kid]
		mu.Unlock()

		if key == nil {
			return nil, fmt.Errorf("no public key found for kid %q", kid)
		}

		return key, nil
	}

	return func(c *gin.Context) {
		log := logger.NewLogrusLogger()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Debug("Keycloak auth: Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Debugf("Keycloak auth: Invalid authorization format, parts: %v", parts)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid authorization format",
			})
			return
		}

		rawToken := strings.TrimSpace(parts[1])
		if rawToken == "" {
			log.Debug("Keycloak auth: Empty bearer token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Empty bearer token",
			})
			return
		}

		token, err := jwt.ParseWithClaims(rawToken, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token is signed with an RSA algorithm (Keycloak default is RS256).
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, _ := token.Header["kid"].(string)
			if kid == "" {
				return nil, errors.New("token header missing kid")
			}

			log.Debugf("Keycloak auth: Token kid: %s", kid)
			return getPublicKeyForKID(kid)
		})
		if err != nil {
			log.Errorf("Keycloak auth: Failed to parse token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token",
				"details": err.Error(),
			})
			return
		}
		if !token.Valid {
			log.Debug("Keycloak auth: Token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(*KeycloakClaims)
		if !ok {
			log.Debug("Keycloak auth: Failed to extract claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token claims",
			})
			return
		}

		log.Debugf("Keycloak auth: Token claims - issuer: %s, subject: %s, audience: %v, azp: %s",
			claims.Issuer, claims.Subject, claims.Audience, claims.AZP)

		// Basic time-based validation.
		now := time.Now()
		if claims.ExpiresAt != nil && !claims.ExpiresAt.After(now) {
			log.Debugf("Keycloak auth: Token expired at %v, now: %v", claims.ExpiresAt, now)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Token is expired",
			})
			return
		}
		if claims.NotBefore != nil && claims.NotBefore.After(now) {
			log.Debugf("Keycloak auth: Token not valid until %v, now: %v", claims.NotBefore, now)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Token is not yet valid",
			})
			return
		}

		if issuer != "" && claims.Issuer != "" && claims.Issuer != issuer {
			log.Debugf("Keycloak auth: Issuer mismatch - expected: %s, got: %s", issuer, claims.Issuer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token issuer",
				"details": fmt.Sprintf("expected: %s, got: %s", issuer, claims.Issuer),
			})
			return
		}

		// Audience validation: In Keycloak, tokens can have "account" as audience
		// when issued for account management, but the actual client ID is in azp (authorized party).
		// We'll check both the audience array and extract azp from raw claims if needed.
		if audience != "" {
			audienceMatch := false
			if len(claims.Audience) > 0 {
				for _, aud := range claims.Audience {
					if aud == audience {
						audienceMatch = true
						break
					}
				}
				log.Debugf("Keycloak auth: Audience check - expected: %s, token audiences: %v, match: %v",
					audience, claims.Audience, audienceMatch)
			}

			// If audience doesn't match, check azp (authorized party)
			// This is common when Keycloak issues tokens with "account" as audience
			// The actual client ID is then in the azp field
			if !audienceMatch && claims.AZP != "" {
				if claims.AZP == audience {
					log.Debugf("Keycloak auth: Audience match via azp: %s", claims.AZP)
					audienceMatch = true
				} else {
					log.Debugf("Keycloak auth: azp found but doesn't match - azp: %s, expected: %s", claims.AZP, audience)
				}
			}

			if !audienceMatch {
				log.Debugf("Keycloak auth: Audience validation failed - expected: %s, got audiences: %v",
					audience, claims.Audience)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid token audience",
					"details": fmt.Sprintf("expected: %s, got: %v", audience, claims.Audience),
				})
				return
			}
		}

		subject := claims.Subject
		if subject == "" {
			log.Debug("Keycloak auth: Token subject (sub) is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Token subject (sub) is missing",
			})
			return
		}

		log.Debugf("Keycloak auth: Token validated successfully for subject: %s", subject)

		// Store user ID and roles in context so handlers can retrieve them.
		ctx := context.WithValue(c.Request.Context(), UserIDKey, subject)
		if len(claims.RealmAccess.Roles) > 0 {
			ctx = context.WithValue(ctx, UserRolesKey, claims.RealmAccess.Roles)
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequireRole ensures that the authenticated user (via KeycloakAuthMiddleware)
// has the given realm role. If not, the request is rejected with 403.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.NewLogrusLogger()

		roles, ok := c.Request.Context().Value(UserRolesKey).([]string)
		if !ok || len(roles) == 0 {
			log.Debug("RequireRole: no roles found in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Forbidden: missing required role",
			})
			return
		}

		for _, r := range roles {
			if r == requiredRole {
				// User has the required role, continue.
				c.Next()
				return
			}
		}

		log.Debugf("RequireRole: user roles %v do not include required role %s", roles, requiredRole)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "Forbidden: missing required role",
		})
	}
}
