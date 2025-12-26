package events

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dws-org/dws-event-service/internal/services"
	"github.com/dws-org/dws-event-service/prisma/db"
	"github.com/shopspring/decimal"
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

// GetEvents godoc
// @Summary      List events
// @Description  Returns a list of all events
// @Tags         events
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /events [get]
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

// GetEventByID godoc
// @Summary      Get event by ID
// @Description  Returns a single event by its ID
// @Tags         events
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /events/{id} [get]
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

// CreateEventRequest represents the JSON payload for creating an event
// @Description  Event creation payload
type CreateEventRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description" binding:"required"`
	StartDate   time.Time       `json:"startDate" binding:"required"`
	StartTime   time.Time       `json:"startTime" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	EndDate     time.Time       `json:"endDate" binding:"required"`
	Location    string          `json:"location" binding:"required"`
	Capacity    int             `json:"capacity" binding:"required"`
	ImageURL    string          `json:"imageUrl" binding:"required"`
	Category    string          `json:"category" binding:"required"`
	OrganizerID string          `json:"organizerId" binding:"required"`
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Creates a new event
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event  body      CreateEventRequest  true  "Event to create"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /events [post]
func (ec *Controller) CreateEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	event, err := ec.dbService.GetClient().Event.CreateOne(
		db.Event.Name.Set(req.Name),
		db.Event.Description.Set(req.Description),
		db.Event.StartDate.Set(req.StartDate),
		db.Event.StartTime.Set(req.StartTime),
		db.Event.Price.Set(req.Price),
		db.Event.EndDate.Set(req.EndDate),
		db.Event.Location.Set(req.Location),
		db.Event.Capacity.Set(req.Capacity),
		db.Event.ImageURL.Set(req.ImageURL),
		db.Event.Category.Set(req.Category),
		db.Event.Organizer.Link(db.Organizer.ID.Equals(req.OrganizerID)),
	).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create event",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}
