package http

import (
	"net/http"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type PermissionActionHandler struct {
	svc *service.PermissionActionService
}

func NewPermissionActionHandler(svc *service.PermissionActionService) *PermissionActionHandler {
	return &PermissionActionHandler{svc: svc}
}

type CreatePermissionActionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionActionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *PermissionActionHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}

	var req CreatePermissionActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := &domain.PermissionAction{
		Name:        req.Name,
		Description: req.Description,
		TenantID:    tenantID.(string),
	}

	if err := h.svc.Create(c.Request.Context(), tenantID.(string), action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "permission action created successfully",
		"action":  action,
	})
}

func (h *PermissionActionHandler) GetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing permission action id"})
		return
	}
	action, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"action": action})
}

func (h *PermissionActionHandler) List(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	actions, err := h.svc.ListByTenant(c.Request.Context(), tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"actions": actions,
		"total":   len(actions),
	})
}

func (h *PermissionActionHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing permission action id"})
		return
	}
	var req UpdatePermissionActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	action := &domain.PermissionAction{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		TenantID:    tenantID.(string),
	}
	if err := h.svc.Update(c.Request.Context(), tenantID.(string), action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission action updated successfully"})
}

func (h *PermissionActionHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing permission action id"})
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission action softdeleted"})
}

func (h *PermissionActionHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing permission action id"})
		return
	}
	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission action restored"})
}
