package http

import (
	"net/http"
	"time"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

type CreateTenantRequest struct {
	ID         string    `json:"id"`
	Name       string    `json:"name" binding:"required"`
	Plan       string    `json:"plan"`
	Domain     string    `json:"domain"`
	IsActive   bool      `json:"is_active"`
	MaxUsers   int       `json:"max_users"`
	Expiration time.Time `json:"expiration"`
}

type UpdateTenantRequest struct {
	Name       string    `json:"name"`
	Plan       string    `json:"plan"`
	Domain     string    `json:"domain"`
	IsActive   *bool     `json:"is_active"`
	MaxUsers   *int      `json:"max_users"`
	Expiration time.Time `json:"expiration"`
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	tenant := &domain.Tenant{
		ID:         req.ID,
		Name:       req.Name,
		Plan:       req.Plan,
		Domain:     req.Domain,
		IsActive:   req.IsActive,
		MaxUsers:   req.MaxUsers,
		Expiration: req.Expiration,
	}

	if err := h.svc.Create(c.Request.Context(), tenant); err != nil {
		respondError(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "tenant created successfully",
		"tenant":  tenant,
	})
}

func (h *TenantHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"tenant": tenant})
}

func (h *TenantHandler) List(c *gin.Context) {
	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)
	tenants, err := h.svc.List(c.Request.Context(), offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
		"total":   len(tenants),
	})
}

func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	tenant := &domain.Tenant{
		ID:         id,
		Name:       req.Name,
		Plan:       req.Plan,
		Domain:     req.Domain,
		Expiration: req.Expiration,
	}
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}
	if req.MaxUsers != nil {
		tenant.MaxUsers = *req.MaxUsers
	}

	if err := h.svc.Update(c.Request.Context(), tenant); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tenant updated successfully"})
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted successfully"})
}
