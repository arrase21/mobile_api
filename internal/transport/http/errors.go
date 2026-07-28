package http

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, status int, msg string) {
	attrs := []any{
		slog.Int("status", status),
		slog.String("path", c.Request.URL.Path),
		slog.String("method", c.Request.Method),
		slog.String("error", msg),
		slog.String("client_ip", c.ClientIP()),
	}

	if tenantID, exists := c.Get("tenant_id"); exists {
		attrs = append(attrs, slog.String("tenant_id", tenantID.(string)))
	}
	if userID, exists := c.Get("user_id"); exists {
		attrs = append(attrs, slog.String("user_id", userID.(string)))
	}

	switch {
	case status >= 500:
		slog.Error("api_error", attrs...)
	case status >= 400:
		slog.Warn("api_error", attrs...)
	default:
		slog.Info("api_error", attrs...)
	}

	c.JSON(status, gin.H{"error": msg})
}

func respondSuccess(c *gin.Context, status int, data gin.H) {
	c.JSON(status, data)
}

func parseQueryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil || n < 0 {
		return defaultVal
	}
	return n
}
