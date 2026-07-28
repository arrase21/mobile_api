package service

import (
	"context"
	"errors"

	"github.com/arrase21/mobileapi/internal/domain"
)

type SkinfoldService struct {
	skinfoldRepo   domain.SkinfoldRepo
	assessmentRepo domain.AssessmentRepo
}

func NewSkinfoldService(skinfoldRepo domain.SkinfoldRepo, assessmentRepo domain.AssessmentRepo) *SkinfoldService {
	return &SkinfoldService{skinfoldRepo: skinfoldRepo, assessmentRepo: assessmentRepo}
}

func (s *SkinfoldService) Create(ctx context.Context, tenantID string, skinfold *domain.Skinfold) error {
	if skinfold == nil {
		return errors.New("skinfold is nil")
	}
	if skinfold.AssessmentID == "" {
		return errors.New("assessment_id is required")
	}

	_, err := s.assessmentRepo.GetByID(ctx, tenantID, skinfold.AssessmentID)
	if err != nil {
		return errors.New("assessment not found in this tenant")
	}

	return s.skinfoldRepo.Create(ctx, tenantID, skinfold)
}

func (s *SkinfoldService) GetByID(ctx context.Context, tenantID, id string) (*domain.Skinfold, error) {
	if id == "" {
		return nil, errors.New("skinfold id is empty")
	}
	return s.skinfoldRepo.GetByID(ctx, tenantID, id)
}

func (s *SkinfoldService) ListByAssessment(ctx context.Context, tenantID, assessmentID string, offset, limit int) ([]*domain.Skinfold, error) {
	if assessmentID == "" {
		return nil, errors.New("assessment id is empty")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.skinfoldRepo.ListByAssessment(ctx, tenantID, assessmentID, offset, limit)
}

func (s *SkinfoldService) Update(ctx context.Context, tenantID string, skinfold *domain.Skinfold) error {
	if skinfold == nil {
		return errors.New("skinfold is nil")
	}
	if skinfold.ID == "" {
		return errors.New("skinfold id is empty")
	}
	return s.skinfoldRepo.Update(ctx, tenantID, skinfold)
}

func (s *SkinfoldService) SoftDelete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("skinfold id is empty")
	}
	return s.skinfoldRepo.SoftDelete(ctx, tenantID, id)
}

func (s *SkinfoldService) Restore(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("skinfold id is empty")
	}
	return s.skinfoldRepo.Restore(ctx, tenantID, id)
}

func (s *SkinfoldService) Delete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("skinfold id is empty")
	}
	return s.skinfoldRepo.Delete(ctx, tenantID, id)
}
