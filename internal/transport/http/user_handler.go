package http

import (
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
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}

	creatorID, _ := c.Get("user_id")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
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

	var creatorIDStr string
	if creatorID != nil {
		creatorIDStr = creatorID.(string)
	}

	if err := h.svc.CreateUser(c.Request.Context(), user, creatorIDStr); err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			respondError(c, http.StatusConflict, err.Error())
		case domain.ErrTenantUserLimitReached:
			respondError(c, http.StatusForbidden, err.Error())
		default:
			respondError(c, http.StatusBadRequest, err.Error())
		}
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
		respondError(c, http.StatusBadRequest, "tenant id not found")
		return
	}

	dni := c.Param("dni")
	if dni == "" {
		respondError(c, http.StatusBadRequest, "missing dni")
		return
	}

	user, err := h.svc.GetByDni(c.Request.Context(), tenantID.(string), dni)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) List(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)

	users, err := h.svc.List(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
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
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}

	userID := c.Param("user")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "missing user id")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
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
		switch err {
		case domain.ErrTenantUserLimitReached:
			respondError(c, http.StatusForbidden, err.Error())
		default:
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	userID := c.Param("user")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "missing user id")
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user softdeleted"})
}

func (h *UserHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	userID := c.Param("user")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "missing user id")
		return
	}
	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user restored"})
}

func (h *UserHandler) ListDeleted(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)
	users, err := h.svc.ListDeleted(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": len(users),
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	userID := c.Param("user")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "missing user id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), tenantID.(string), userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted permanently"})
}
