package configs

// Wir importieren das testing-Paket von Go
// und testify/assert für einfachere Vergleiche
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TEST 1: Kann ich ein Service-Struct erstellen und die Werte lesen?
// =============================================================================
//
// Jede Test-Funktion:
// - Muss mit "Test" beginnen
// - Bekommt ein "t *testing.T" als Parameter (damit können wir Fehler melden)
//
func TestServiceStruct_CanCreateAndReadValues(t *testing.T) {
	// SCHRITT 1: Erstelle ein Service-Objekt mit Test-Daten
	service := Service{
		Name:        "Mein Test Service",
		Slug:        "mein-test-service",
		Description: "Das ist eine Beschreibung",
		Version:     "1.0.0",
		Tags:        []string{"tag1", "tag2"},
	}

	// SCHRITT 2: Prüfe ob die Werte korrekt gespeichert wurden
	//
	// assert.Equal(t, erwartet, tatsächlich)
	// -> Prüft: Ist "erwartet" gleich "tatsächlich"?
	// -> Wenn nicht, schlägt der Test fehl

	assert.Equal(t, "Mein Test Service", service.Name)
	assert.Equal(t, "mein-test-service", service.Slug)
	assert.Equal(t, "Das ist eine Beschreibung", service.Description)
	assert.Equal(t, "1.0.0", service.Version)

	// Für die Tags prüfen wir die Länge
	assert.Len(t, service.Tags, 2) // Erwartet: 2 Tags
}

// =============================================================================
// TEST 2: Was passiert wenn ich ein leeres Service-Struct erstelle?
// =============================================================================
func TestServiceStruct_EmptyValues(t *testing.T) {
	// Erstelle ein leeres Service (ohne Werte)
	service := Service{}

	// Alle String-Felder sollten leer sein
	assert.Empty(t, service.Name)        // Name sollte "" sein
	assert.Empty(t, service.Slug)        // Slug sollte "" sein
	assert.Empty(t, service.Description) // Description sollte "" sein
	assert.Empty(t, service.Version)     // Version sollte "" sein

	// Tags sollte nil oder leer sein
	assert.Nil(t, service.Tags) // Tags wurde nicht initialisiert
}

// =============================================================================
// TEST 3: Kann ich die Tags richtig auslesen?
// =============================================================================
func TestServiceStruct_TagsWork(t *testing.T) {
	service := Service{
		Name: "Tag Test",
		Tags: []string{"golang", "microservice", "api"},
	}

	// Prüfe die Anzahl der Tags
	assert.Len(t, service.Tags, 3)

	// Prüfe ob bestimmte Tags enthalten sind
	assert.Contains(t, service.Tags, "golang")
	assert.Contains(t, service.Tags, "microservice")
	assert.Contains(t, service.Tags, "api")

	// Prüfe dass ein falscher Tag NICHT enthalten ist
	assert.NotContains(t, service.Tags, "python")
}

// =============================================================================
// TEST 4: Sub-Tests mit t.Run() - Gruppierte Tests
// =============================================================================
//
// Mit t.Run() können wir mehrere kleine Tests in einem gruppieren.
// Das macht die Ausgabe übersichtlicher.
//
func TestServiceStruct_SubTests(t *testing.T) {

	// Sub-Test 1
	t.Run("Name kann gesetzt werden", func(t *testing.T) {
		service := Service{Name: "Test"}
		assert.Equal(t, "Test", service.Name)
	})

	// Sub-Test 2
	t.Run("Version kann gesetzt werden", func(t *testing.T) {
		service := Service{Version: "2.0.0"}
		assert.Equal(t, "2.0.0", service.Version)
	})

	// Sub-Test 3
	t.Run("Leere Tags sind nil", func(t *testing.T) {
		service := Service{}
		assert.Nil(t, service.Tags)
	})

	// Sub-Test 4
	t.Run("Tags können hinzugefügt werden", func(t *testing.T) {
		service := Service{
			Tags: []string{"eins"},
		}
		// Tags erweitern
		service.Tags = append(service.Tags, "zwei")

		assert.Len(t, service.Tags, 2)
		assert.Equal(t, "zwei", service.Tags[1])
	})
}