package http

import (
	"net/http"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	svc *service.PermissionService
}

func NewPermissionHandler(svc *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

type CreatePermissionRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
}

type UpdatePermissionRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
}

func (h *PermissionHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}

	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
		Actions:     req.Actions,
		TenantID:    tenantID.(string),
	}

	if err := h.svc.Create(c.Request.Context(), tenantID.(string), permission); err != nil {
		respondError(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":    "permission created successfully",
		"permission": permission,
	})
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing permission id")
		return
	}
	permission, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

func (h *PermissionHandler) List(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)
	permissions, err := h.svc.ListByTenant(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
		"total":       len(permissions),
	})
}

func (h *PermissionHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing permission id")
		return
	}
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	permission := &domain.Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Actions:     req.Actions,
		TenantID:    tenantID.(string),
	}
	if err := h.svc.Update(c.Request.Context(), tenantID.(string), permission); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission updated successfully"})
}

func (h *PermissionHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing permission id")
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission softdeleted"})
}

func (h *PermissionHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing permission id")
		return
	}
	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission restored"})
}
