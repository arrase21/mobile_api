package service

import (
	"context"
	"errors"

	"github.com/arrase21/mobileapi/internal/domain"
)

type AssessmentService struct {
	assessmentRepo domain.AssessmentRepo
	userRepo       domain.UserRepo
}

func NewAssessmentService(assessmentRepo domain.AssessmentRepo, userRepo domain.UserRepo) *AssessmentService {
	return &AssessmentService{assessmentRepo: assessmentRepo, userRepo: userRepo}
}

func (s *AssessmentService) Create(ctx context.Context, tenantID string, assessment *domain.Assessment) error {
	if assessment == nil {
		return errors.New("assessment is nil")
	}
	if assessment.UserID == "" {
		return errors.New("user_id is required")
	}

	_, err := s.userRepo.GetByID(ctx, tenantID, assessment.UserID)
	if err != nil {
		return errors.New("user not found in this tenant")
	}

	return s.assessmentRepo.Create(ctx, tenantID, assessment)
}

func (s *AssessmentService) GetByID(ctx context.Context, tenantID, id string) (*domain.Assessment, error) {
	if id == "" {
		return nil, errors.New("assessment id is empty")
	}
	return s.assessmentRepo.GetByID(ctx, tenantID, id)
}

func (s *AssessmentService) ListByUser(ctx context.Context, tenantID, userID string, offset, limit int) ([]*domain.Assessment, error) {
	if userID == "" {
		return nil, errors.New("user id is empty")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.assessmentRepo.ListByUser(ctx, tenantID, userID, offset, limit)
}

func (s *AssessmentService) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Assessment, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.assessmentRepo.ListByTenant(ctx, tenantID, offset, limit)
}

func (s *AssessmentService) Update(ctx context.Context, tenantID string, assessment *domain.Assessment) error {
	if assessment == nil {
		return errors.New("assessment is nil")
	}
	if assessment.ID == "" {
		return errors.New("assessment id is empty")
	}
	return s.assessmentRepo.Update(ctx, tenantID, assessment)
}

func (s *AssessmentService) SoftDelete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("assessment id is empty")
	}
	return s.assessmentRepo.SoftDelete(ctx, tenantID, id)
}

func (s *AssessmentService) Restore(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("assessment id is empty")
	}
	return s.assessmentRepo.Restore(ctx, tenantID, id)
}

func (s *AssessmentService) Delete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("assessment id is empty")
	}
	return s.assessmentRepo.Delete(ctx, tenantID, id)
}
