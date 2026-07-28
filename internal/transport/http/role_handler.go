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
	Permissions []string `json:"permissions"`
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		TenantID:    tenantID.(string),
	}
	if err := h.svc.CreateRole(c.Request.Context(), role); err != nil {
		respondError(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "role create successfully",
		"role":    role,
	})
}

func (h *RoleHandler) List(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)

	roles, err := h.svc.ListByTenant(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"total": len(roles),
	})
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing role id")
		return
	}
	role, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role})
}

type UpdateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (h *RoleHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing role id")
		return
	}
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	role := &domain.Role{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		TenantID:    tenantID.(string),
	}
	if err := h.svc.Update(c.Request.Context(), tenantID.(string), role); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated successfully"})
}

func (h *RoleHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing role id")
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role softdeleted"})
}

func (h *RoleHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing role id")
		return
	}
	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role restored"})
}
