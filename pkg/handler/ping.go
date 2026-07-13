package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ping godoc
// @Summary Health check
// @Tags public
// @Produce json
// @Success 200 {object} model.APIResponse
// @Failure 503 {object} model.APIResponse
// @Router /ping/ [get]
func (h *Handler) ping(c *gin.Context) {
	status := "ok"
	dbStatus := "ok"
	if err := h.services.Health.PingDB(); err != nil {
		dbStatus = "error"
		status = "degraded"
	}

	code := http.StatusOK
	if dbStatus == "error" {
		code = http.StatusServiceUnavailable
	}

	Success(c, code, gin.H{
		"status": status,
		"db":     dbStatus,
	})
}
