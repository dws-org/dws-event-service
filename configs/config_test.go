package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TEST: GetEnvConfig() Funktion
// =============================================================================
func TestGetEnvConfig(t *testing.T) {

	t.Run("gibt Config zurück wenn bereits geladen", func(t *testing.T) {
		// Zuerst eine Test-Config setzen
		testConfig := &Config{
			Service: Service{
				Name:    "Test Service",
				Version: "1.0.0",
			},
		}
		
		// Config setzen (simuliert dass sie schon geladen wurde)
		EnvConfig = testConfig

		// Jetzt GetEnvConfig aufrufen
		result := GetEnvConfig()

		// Prüfen ob wir die gleiche Config zurückbekommen
		assert.NotNil(t, result)
		assert.Equal(t, "Test Service", result.Service.Name)
		assert.Equal(t, "1.0.0", result.Service.Version)
		
		// Aufräumen
		EnvConfig = nil
	})
}

// =============================================================================
// TEST: Config Struct komplett
// =============================================================================
func TestConfigStruct(t *testing.T) {

	t.Run("alle Felder können gesetzt werden", func(t *testing.T) {
		config := &Config{
			Service: Service{
				Name:        "Event Service",
				Slug:        "event-service",
				Description: "Manages events",
				Version:     "2.0.0",
				Tags:        []string{"events", "api"},
			},
			Server: Server{
				Host:    "localhost",
				Port:    ":8080",
				GinMode: "debug",
			},
			JWT: JWT{
				Secret: "super-secret-key",
			},
		}

		// Service prüfen
		assert.Equal(t, "Event Service", config.Service.Name)
		assert.Equal(t, "event-service", config.Service.Slug)
		
		// Server prüfen
		assert.Equal(t, "localhost", config.Server.Host)
		assert.Equal(t, ":8080", config.Server.Port)
		
		// JWT prüfen
		assert.Equal(t, "super-secret-key", config.JWT.Secret)
	})
}

// =============================================================================
// TEST: EnvConfig Variable direkt setzen und lesen
// =============================================================================
func TestEnvConfigVariable(t *testing.T) {

	t.Run("EnvConfig ist anfangs nil", func(t *testing.T) {
		// Sicherstellen dass wir sauber starten
		EnvConfig = nil
		
		assert.Nil(t, EnvConfig)
	})

	t.Run("EnvConfig kann gesetzt werden", func(t *testing.T) {
		EnvConfig = &Config{
			Service: Service{Name: "Test"},
		}

		assert.NotNil(t, EnvConfig)
		assert.Equal(t, "Test", EnvConfig.Service.Name)
		
		// Aufräumen
		EnvConfig = nil
	})
}