package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEvent_ValidationRules(t *testing.T) {
	tests := []struct {
		name    string
		event   Event
		isValid bool
	}{
		{
			name: "valid_complete_event",
			event: Event{
				Name:        "Tech Conference 2024",
				Description: "Annual technology conference",
				Date:        time.Now().Add(24 * time.Hour),
				Location:    "Convention Center",
				Capacity:    500,
				Price:       "99.99",
				Category:    "Technology",
			},
			isValid: true,
		},
		{
			name: "empty_name",
			event: Event{
				Name:        "",
				Description: "Description",
				Date:        time.Now().Add(24 * time.Hour),
				Location:    "Location",
				Capacity:    100,
			},
			isValid: false,
		},
		{
			name: "past_date",
			event: Event{
				Name:        "Past Event",
				Description: "Description",
				Date:        time.Now().Add(-24 * time.Hour),
				Location:    "Location",
				Capacity:    100,
			},
			isValid: false,
		},
		{
			name: "zero_capacity",
			event: Event{
				Name:        "Event",
				Description: "Description",
				Date:        time.Now().Add(24 * time.Hour),
				Location:    "Location",
				Capacity:    0,
			},
			isValid: false,
		},
		{
			name: "negative_capacity",
			event: Event{
				Name:        "Event",
				Description: "Description",
				Date:        time.Now().Add(24 * time.Hour),
				Location:    "Location",
				Capacity:    -10,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.event.Name)
				assert.True(t, tt.event.Date.After(time.Now()))
				assert.Greater(t, tt.event.Capacity, 0)
			} else {
				if tt.event.Name == "" {
					assert.Empty(t, tt.event.Name)
				}
				if tt.event.Capacity <= 0 {
					assert.LessOrEqual(t, tt.event.Capacity, 0)
				}
			}
		})
	}
}

func TestEvent_PriceFormatting(t *testing.T) {
	tests := []struct {
		name     string
		price    string
		expected string
	}{
		{"free_event", "0.00", "0.00"},
		{"standard_price", "25.50", "25.50"},
		{"high_price", "999.99", "999.99"},
		{"whole_number", "100", "100"},
		{"decimal_one_place", "9.9", "9.9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{Price: tt.price}
			assert.Equal(t, tt.expected, event.Price)
		})
	}
}

func TestEvent_DateHandling(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		date     time.Time
		isFuture bool
	}{
		{"tomorrow", now.Add(24 * time.Hour), true},
		{"next_week", now.Add(7 * 24 * time.Hour), true},
		{"next_month", now.Add(30 * 24 * time.Hour), true},
		{"yesterday", now.Add(-24 * time.Hour), false},
		{"last_week", now.Add(-7 * 24 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{Date: tt.date}
			if tt.isFuture {
				assert.True(t, event.Date.After(now))
			} else {
				assert.True(t, event.Date.Before(now))
			}
		})
	}
}

func TestEvent_CategoryValidation(t *testing.T) {
	validCategories := []string{
		"Music",
		"Sports",
		"Technology",
		"Art",
		"Food",
		"Business",
		"Education",
		"Entertainment",
	}

	for _, category := range validCategories {
		t.Run("category_"+category, func(t *testing.T) {
			event := Event{Category: category}
			assert.NotEmpty(t, event.Category)
			assert.Contains(t, validCategories, event.Category)
		})
	}
}

func TestEvent_CapacityBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		isValid  bool
	}{
		{"small_event", 10, true},
		{"medium_event", 100, true},
		{"large_event", 1000, true},
		{"huge_event", 10000, true},
		{"one_person", 1, true},
		{"zero_capacity", 0, false},
		{"negative_capacity", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{Capacity: tt.capacity}
			if tt.isValid {
				assert.Greater(t, event.Capacity, 0)
			} else {
				assert.LessOrEqual(t, event.Capacity, 0)
			}
		})
	}
}

func TestEvent_LocationValidation(t *testing.T) {
	locations := []string{
		"Convention Center",
		"Online",
		"City Hall",
		"Stadium",
		"Virtual",
		"123 Main St, City, Country",
	}

	for _, location := range locations {
		t.Run("location_"+location, func(t *testing.T) {
			event := Event{Location: location}
			assert.NotEmpty(t, event.Location)
		})
	}
}

func TestEvent_DescriptionLength(t *testing.T) {
	tests := []struct {
		name        string
		description string
		isValid     bool
	}{
		{"short", "Short description", true},
		{"medium", "This is a medium length description with some details", true},
		{"long", "This is a very long description with lots of details about the event, including information about speakers, agenda, venue, food, parking, and much more that attendees need to know before registering", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{Description: tt.description}
			if tt.isValid {
				assert.NotEmpty(t, event.Description)
			} else {
				assert.Empty(t, event.Description)
			}
		})
	}
}

func TestEvent_MultipleValidEvents(t *testing.T) {
	events := []Event{
		{
			Name:        "Concert",
			Description: "Live music event",
			Date:        time.Now().Add(24 * time.Hour),
			Location:    "Arena",
			Capacity:    5000,
			Price:       "50.00",
			Category:    "Music",
		},
		{
			Name:        "Workshop",
			Description: "Tech workshop",
			Date:        time.Now().Add(48 * time.Hour),
			Location:    "Office",
			Capacity:    20,
			Price:       "0.00",
			Category:    "Technology",
		},
		{
			Name:        "Marathon",
			Description: "City marathon",
			Date:        time.Now().Add(72 * time.Hour),
			Location:    "City Center",
			Capacity:    1000,
			Price:       "25.00",
			Category:    "Sports",
		},
	}

	for i, event := range events {
		t.Run("event_"+string(rune(i+'1')), func(t *testing.T) {
			assert.NotEmpty(t, event.Name)
			assert.NotEmpty(t, event.Description)
			assert.True(t, event.Date.After(time.Now()))
			assert.NotEmpty(t, event.Location)
			assert.Greater(t, event.Capacity, 0)
			assert.NotEmpty(t, event.Price)
			assert.NotEmpty(t, event.Category)
		})
	}
}

func TestEvent_TimezonHandling(t *testing.T) {
	now := time.Now()
	utc := now.UTC()
	local := now.Local()

	event1 := Event{Date: utc}
	event2 := Event{Date: local}

	assert.NotNil(t, event1.Date)
	assert.NotNil(t, event2.Date)
	
	// Both should represent essentially the same moment
	diff := event1.Date.Sub(event2.Date)
	assert.Less(t, diff.Abs(), time.Second)
}

func TestEvent_Equality(t *testing.T) {
	date := time.Now().Add(24 * time.Hour)
	
	event1 := Event{
		Name:        "Event A",
		Description: "Description A",
		Date:        date,
		Location:    "Location A",
		Capacity:    100,
		Price:       "50.00",
		Category:    "Technology",
	}

	event2 := Event{
		Name:        "Event A",
		Description: "Description A",
		Date:        date,
		Location:    "Location A",
		Capacity:    100,
		Price:       "50.00",
		Category:    "Technology",
	}

	assert.Equal(t, event1.Name, event2.Name)
	assert.Equal(t, event1.Description, event2.Description)
	assert.Equal(t, event1.Date, event2.Date)
	assert.Equal(t, event1.Location, event2.Location)
	assert.Equal(t, event1.Capacity, event2.Capacity)
	assert.Equal(t, event1.Price, event2.Price)
	assert.Equal(t, event1.Category, event2.Category)
}
