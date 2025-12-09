package http

import (
	"net/http"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc}
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permission  []string `json:"permission"`
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	tenantID, exists := c.Get("tenat_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missig tenatn id"})
		return
	}
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permission,
		TenantID:    tenantID.(string),
	}
	if err := h.svc.CreateRole(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "role create successfully",
		"role":    role,
	})
}

func (h *RoleHandler) List(c *gin.Context) {
	tenantID, exits := c.Get("tenat_id")
	if !exits {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missig tenant id"})
		return
	}
	offset := 0
	limit := 20

	roles, err := h.svc.ListByTenant(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"total": len(roles),
	})
}
