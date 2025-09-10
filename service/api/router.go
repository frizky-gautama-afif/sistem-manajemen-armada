package api

import "github.com/gin-gonic/gin"

// SetupRouter configures and returns a Gin router with all API endpoints.
func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/vehicles")
	{
		api.GET("/:vehicle_id/location", h.GetLastLocation)
		api.GET("/:vehicle_id/history", h.GetLocationHistory)
	}
	return r
}
