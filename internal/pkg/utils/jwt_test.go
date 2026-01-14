package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/stretchr/testify/assert"
)

// Hilfsfunktion: setzt eine Test-Config (damit GenerateToken/JwtVerify ein Secret haben)
func setTestEnvConfig(secret string) func() {
	orig := configs.EnvConfig
	configs.EnvConfig = &configs.Config{
		JWT: configs.JWT{
			Secret: secret,
		},
	}

	// Cleanup-Funktion zurückgeben
	return func() {
		configs.EnvConfig = orig
	}
}

func TestGenerateToken_And_JwtVerify_HappyPath(t *testing.T) {
	cleanup := setTestEnvConfig("test-secret-123")
	defer cleanup()

	claims := &Claims{
		Uid:      42,
		Username: "lea",
	}

	token := GenerateToken(claims)
	assert.NotEmpty(t, token)

	verified, err := JwtVerify(token)
	assert.NoError(t, err)
	assert.NotNil(t, verified)

	assert.Equal(t, uint(42), verified.Uid)
	assert.Equal(t, "lea", verified.Username)

	// ExpiresAt sollte gesetzt sein (und in der Zukunft liegen)
	assert.NotNil(t, verified.ExpiresAt)
	assert.True(t, verified.ExpiresAt.Time.After(time.Now()))
}

func TestJwtVerify_InvalidTokenString(t *testing.T) {
	cleanup := setTestEnvConfig("test-secret-123")
	defer cleanup()

	_, err := JwtVerify("this-is-not-a-jwt")
	assert.Error(t, err)

	// Dein Code gibt "token invalid" zurück
	assert.Contains(t, err.Error(), "token invalid")
}

func TestJwtVerify_WrongSecret(t *testing.T) {
	// Token mit Secret A erstellen
	cleanupA := setTestEnvConfig("secret-A")
	claims := &Claims{Uid: 1, Username: "user"}
	token := GenerateToken(claims)
	cleanupA()

	// Jetzt Verify mit Secret B -> muss failen
	cleanupB := setTestEnvConfig("secret-B")
	defer cleanupB()

	_, err := JwtVerify(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token invalid")
}

func TestJwtVerify_ExpiredToken(t *testing.T) {
	cleanup := setTestEnvConfig("test-secret-123")
	defer cleanup()

	// Wichtig: GenerateToken überschreibt ExpiresAt immer auf now+30m,
	// daher bauen wir hier einen Token manuell mit ExpiresAt in der Vergangenheit.
	expiredClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
		Uid:      99,
		Username: "expired-user",
	}

	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).
		SignedString([]byte(configs.EnvConfig.JWT.Secret))
	assert.NoError(t, err)
	assert.NotEmpty(t, expiredToken)

	_, verr := JwtVerify(expiredToken)
	assert.Error(t, verr)

	// Je nach Verhalten der jwt/v5 Validierung kann es "token invalid" ODER "token expired" sein.
	// Dein Code hat beide Pfade möglich.
	msg := verr.Error()
	assert.True(t,
		msg == "token invalid" || msg == "token expired",
		"expected 'token invalid' or 'token expired', got: "+msg,
	)
}
