package events

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oskargbc/dws-event-service.git/internal/services"
	"github.com/oskargbc/dws-event-service.git/prisma/db"
)

// Controller handles event-related HTTP requests
type Controller struct {
	dbService *services.DatabaseService
}

// NewController creates a new events controller
func NewController() *Controller {
	return &Controller{
		dbService: services.GetDatabaseSeviceInstance(),
	}
}

// GetEvents returns a list of all events
// GET /api/v1/events
func (ec *Controller) GetEvents(c *gin.Context) {
	ctx := c.Request.Context()

	events, err := ec.dbService.GetClient().Event.FindMany().Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch events",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventByID returns a single event by ID
// GET /api/v1/events/:id
func (ec *Controller) GetEventByID(c *gin.Context) {
	ctx := c.Request.Context()
	eventID := c.Param("id")

	event, err := ec.dbService.GetClient().Event.FindUnique(
		db.Event.ID.Equals(eventID),
	).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Event not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}
