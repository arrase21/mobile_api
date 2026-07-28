package http

import (
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
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}

	var req CreateSkinfoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
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
		respondError(c, http.StatusBadRequest, err.Error())
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
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing skinfold id")
		return
	}

	skinfold, err := h.svc.GetByID(c.Request.Context(), tenantID.(string), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"skinfold": skinfold})
}

func (h *SkinfoldHandler) ListByAssessment(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	assessmentID := c.Param("assessment_id")
	if assessmentID == "" {
		respondError(c, http.StatusBadRequest, "missing assessment id")
		return
	}

	offset := parseQueryInt(c.Query("offset"), 0)
	limit := parseQueryInt(c.Query("limit"), 20)

	skinfolds, err := h.svc.ListByAssessment(c.Request.Context(), tenantID.(string), assessmentID, offset, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
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
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing skinfold id")
		return
	}

	var req UpdateSkinfoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
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
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold updated successfully"})
}

func (h *SkinfoldHandler) SoftDelete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing skinfold id")
		return
	}

	if err := h.svc.SoftDelete(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold soft deleted successfully"})
}

func (h *SkinfoldHandler) Restore(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing skinfold id")
		return
	}

	if err := h.svc.Restore(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold restored successfully"})
}

func (h *SkinfoldHandler) Delete(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		respondError(c, http.StatusBadRequest, "missing tenant id")
		return
	}
	id := c.Param("id")
	if id == "" {
		respondError(c, http.StatusBadRequest, "missing skinfold id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantID.(string), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skinfold deleted permanently"})
}
