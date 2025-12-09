package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc}
}

type CreateUserRequest struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Dni       string    `json:"dni" binding:"required"`
	Gender    string    `json:"gender" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	DateBirth time.Time `json:"date_birth" binding:"required"`
	Nickname  string    `json:"nickname" binding:"required"`
}

type UpdateUserRequest struct {
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Dni       string    `json:"dni"`
	Gender    string    `json:"gender"`
	Phone     string    `json:"phone"`
	DateBirth time.Time `json:"date_birth"`
	Nickname  string    `json:"nickname"`
}

func (h *UserHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dni:       req.Dni,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		DateBirth: req.DateBirth,
		Nickname:  req.Nickname,
		TenantID:  tenantID.(string),
	}

	if err := h.svc.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    user,
	})
}
func (h *UserHandler) GetByDni(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant id not found"})
		return
	}

	dni := c.Param("dni")
	if dni == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing dni"})
		return
	}

	user, err := h.svc.GetByDni(c.Request.Context(), tenantID.(string), dni)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) List(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	offset := 0
	limit := 20
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	users, err := h.svc.List(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": len(users),
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}

	userID := c.Param("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := &domain.User{
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dni:       req.Dni,
		Gender:    req.Gender,
		Phone:     req.Phone,
		DateBirth: req.DateBirth,
		Nickname:  req.Nickname,
		TenantID:  tenantID.(string),
	}
	if err := h.svc.Update(c.Request.Context(), tenantID.(string), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	userID := c.Param("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user softdeleted"})
}

func (h *UserHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	userID := c.Param("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user restored"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	userID := c.Param("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), tenantID.(string), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted permanently"})
}
