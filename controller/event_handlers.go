package controller

import (
	"errors"
	"gotempl/controller/service"
	"gotempl/model"
	"gotempl/views/crud"
	"gotempl/views/layout"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventHandler struct {
	Service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{Service: service}
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Create a new event with the provided information
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        event  body      model.Event  true  "Event information"
// @Success      201   {object}  model.Event
// @Failure      400   {object}  object
// @Failure      500   {object}  object
// @Router       /event [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event model.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateEvent(&event); err != nil {
		log.Error("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetAllEvents godoc
// @Summary      Get all events
// @Description  Retrieve a list of all events
// @Tags         Event
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.Event
// @Failure      500  {object}  map[string]string
// @Router       /event [get]
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	events, err := h.Service.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent godoc
// @Summary      Get a event by ID
// @Description  Retrieve a event's information using their ID
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      200  {object}  model.Event
// @Failure      400  {object}  object
// @Failure      404  {object}  object
// @Router       /event/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	event, err := h.Service.GetEventByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent godoc
// @Summary      Update a event
// @Description  Update a event's information in the system
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id    path      string     true  "Event ID"
// @Param        event  body      model.Event true  "Updated event information"
// @Success      200   {object}  model.Event
// @Failure      400   {object}  object
// @Failure      500   {object}  object
// @Router       /event/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = uint64(id)
	if err := h.Service.UpdateEvent(&event); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Delete a event
// @Description  Delete a event from the system using their ID.
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  object
// @Router       /event/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.Service.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.Status(http.StatusNoContent)
}

// EventCRUDHandler godoc
// @Summary      This is a non-REST endpoint that returns an HTML page - not JSON data
// @Description  Fetches all events and renders an HTML page with a CRUD form for event management (non-REST endpoint)
// @Tags         Event
// @Produce      html
// @Success      200  {string}  string  "HTML page content"
// @Router       /admin/event/ [get]
// @Notes
func (h *EventHandler) EventCRUDHandler(c *gin.Context) {
	events, err := h.Service.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	//print(events)
	layout.Render(c, 200, crud.EventForm(events))
}
