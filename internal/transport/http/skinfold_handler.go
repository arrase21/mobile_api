package http

import (
	"fmt"
	"net/http"

	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

type SkinfoldHandler struct {
	svc *service.SkinfoldService
}

func NewSkinfoldHandler(svc *service.SkinfoldService) *SkinfoldHandler {
	return &SkinfoldHandler{svc}
}

type CreateSkinfoldRequest struct {
	AssessmentID  string  `json:"assessment_id" binding:"required"`
	Tricipital    float64 `json:"tricipital" binding:"required"`
	Subscapular   float64 `json:"subscapular" binding:"required"`
	Suprailiaco   float64 `json:"suprailiaco" binding:"required"`
	Abdominal     float64 `json:"abdominal" binding:"required"`
	Thighs        float64 `json:"thighs" binding:"required"`
	Leg           float64 `json:"leg" binding:"required"`
}

type UpdateSkinfoldRequest struct {
	Tricipital  float64 `json:"tricipital"`
	Subscapular float64 `json:"subscapular"`
	Suprailiaco float64 `json:"suprailiaco"`
	Abdominal   float64 `json:"abdominal"`
	Thighs      float64 `json:"thighs"`
	Leg         float64 `json:"leg"`
}

func (h *SkinfoldHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}

	var req CreateSkinfoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skinfold := &domain.Skinfold{
		AssessmentID: req.AssessmentID,
		Tricipital:   req.Tricipital,
		Subscapular:  req.Subscapular,
		Suprailiaco:  req.Suprailiaco,
		Abdominal:    req.Abdominal,
		Thighs:       req.Thighs,
		Leg:          req.Leg,
	}

	if err := h.svc.Create(c.Request.Context(), tenantID.(string), skinfold); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "skinfold created successfully",
		"skinfold": skinfold,
	})
}

func (h *SkinfoldHandler) GetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing skinfold id"})
		return
	}

	skinfold, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"skinfold": skinfold})
}

func (h *SkinfoldHandler) ListByAssessment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	assessmentID := c.Param("assessment_id")
	if assessmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing assessment id"})
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

	skinfolds, err := h.svc.ListByAssessment(c.Request.Context(), tenantID.(string), assessmentID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"skinfolds": skinfolds,
		"total":     len(skinfolds),
	})
}

func (h *SkinfoldHandler) Update(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing skinfold id"})
		return
	}

	var req UpdateSkinfoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skinfold := &domain.Skinfold{
		ID:          id,
		Tricipital:  req.Tricipital,
		Subscapular: req.Subscapular,
		Suprailiaco: req.Suprailiaco,
		Abdominal:   req.Abdominal,
		Thighs:      req.Thighs,
		Leg:         req.Leg,
	}

	if err := h.svc.Update(c.Request.Context(), tenantID.(string), skinfold); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold updated successfully"})
}

func (h *SkinfoldHandler) Delete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing skinfold id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantID.(string), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold deleted successfully"})
}
