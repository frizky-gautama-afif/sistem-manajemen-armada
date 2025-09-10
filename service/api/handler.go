package api

import (
	"net/http"
	"sistem-manajemen-armada/service/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	loc, err := h.DB.GetLastLocation(vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last location"})
		return
	}
	if loc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (h *Handler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp"})
		return
	}
	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp"})
		return
	}

	history, err := h.DB.GetLocationHistory(vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get location history"})
		return
	}
	c.JSON(http.StatusOK, history)
}
