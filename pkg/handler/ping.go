package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
