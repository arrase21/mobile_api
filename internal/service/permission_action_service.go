package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/mobileapi/internal/domain"
)

type PermissionActionService struct {
	actionRepo domain.PermissionActionRepo
}

func NewPermissionActionService(actionRepo domain.PermissionActionRepo) *PermissionActionService {
	return &PermissionActionService{actionRepo: actionRepo}
}

func (s *PermissionActionService) Create(ctx context.Context, tenantID string, action *domain.PermissionAction) error {
	if action == nil {
		return errors.New("permission action is nil")
	}
	action.Name = strings.TrimSpace(action.Name)
	action.Description = strings.TrimSpace(action.Description)
	if action.Name == "" {
		return errors.New("permission action name cannot be empty")
	}
	return s.actionRepo.Create(ctx, tenantID, action)
}

func (s *PermissionActionService) GetByID(ctx context.Context, tenantID, id string) (*domain.PermissionAction, error) {
	if id == "" {
		return nil, errors.New("permission action id is empty")
	}
	return s.actionRepo.GetByID(ctx, tenantID, id)
}

func (s *PermissionActionService) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.PermissionAction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.actionRepo.ListByTenant(ctx, tenantID, offset, limit)
}

func (s *PermissionActionService) Update(ctx context.Context, tenantID string, action *domain.PermissionAction) error {
	if action == nil {
		return errors.New("permission action is nil")
	}
	if action.ID == "" {
		return errors.New("permission action id is empty")
	}
	return s.actionRepo.Update(ctx, tenantID, action)
}

func (s *PermissionActionService) SoftDelete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("permission action id cannot be empty")
	}
	return s.actionRepo.SoftDelete(ctx, tenantID, id)
}

func (s *PermissionActionService) Restore(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("permission action id cannot be empty")
	}
	return s.actionRepo.Restore(ctx, tenantID, id)
}
