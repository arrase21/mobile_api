package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type AssessmentHandler struct {
	svc *service.AssessmentService
}

func NewAssessmentHandler(svc *service.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{svc}
}

type CreateAssessmentRequest struct {
	UserID  string    `json:"user_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
	Height  float64   `json:"height" binding:"required"`
	Weight  float64   `json:"weight" binding:"required"`
	Humerus float64   `json:"humerus" binding:"required"`
	Femur   float64   `json:"femur" binding:"required"`
}

type UpdateAssessmentRequest struct {
	Date    time.Time `json:"date"`
	Height  float64   `json:"height"`
	Weight  float64   `json:"weight"`
	Humerus float64   `json:"humerus"`
	Femur   float64   `json:"femur"`
}

func (h *AssessmentHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}

	var req CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment := &domain.Assessment{
		UserID:  req.UserID,
		Date:    req.Date,
		Height:  req.Height,
		Weight:  req.Weight,
		Humerus: req.Humerus,
		Femur:   req.Femur,
	}

	if err := h.svc.Create(c.Request.Context(), tenantID.(string), assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "assessment created successfully",
		"assessment": assessment,
	})
}

func (h *AssessmentHandler) GetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing assessment id"})
		return
	}

	assessment, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"assessment": assessment})
}

func (h *AssessmentHandler) ListByUser(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
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

	assessments, err := h.svc.ListByUser(c.Request.Context(), tenantID.(string), userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
		"total":       len(assessments),
	})
}

func (h *AssessmentHandler) ListByTenant(c *gin.Context) {
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

	assessments, err := h.svc.ListByTenant(c.Request.Context(), tenantID.(string), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
		"total":       len(assessments),
	})
}

func (h *AssessmentHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing assessment id"})
		return
	}

	var req UpdateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment := &domain.Assessment{
		ID:      id,
		Date:    req.Date,
		Height:  req.Height,
		Weight:  req.Weight,
		Humerus: req.Humerus,
		Femur:   req.Femur,
	}

	if err := h.svc.Update(c.Request.Context(), tenantID.(string), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assessment updated successfully"})
}

func (h *AssessmentHandler) Delete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing assessment id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantID.(string), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assessment deleted successfully"})
}
