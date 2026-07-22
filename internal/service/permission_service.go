package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/mobileapi/internal/domain"
)

type PermissionService struct {
	permissionRepo domain.PermissionRepo
}

func NewPermissionService(permissionRepo domain.PermissionRepo) *PermissionService {
	return &PermissionService{permissionRepo: permissionRepo}
}

func (s *PermissionService) Create(ctx context.Context, tenantID string, permission *domain.Permission) error {
	if permission == nil {
		return errors.New("permission is nil")
	}
	permission.Name = strings.TrimSpace(permission.Name)
	permission.Description = strings.TrimSpace(permission.Description)
	if permission.Name == "" {
		return errors.New("permission name cannot be empty")
	}
	existing, err := s.permissionRepo.GetByName(ctx, tenantID, permission.Name)
	if err != nil && !errors.Is(err, domain.ErrPermissionNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrPermissionAlreadyExists
	}
	return s.permissionRepo.Create(ctx, tenantID, permission)
}

func (s *PermissionService) GetByID(ctx context.Context, tenantID, id string) (*domain.Permission, error) {
	if id == "" {
		return nil, errors.New("permission id is empty")
	}
	return s.permissionRepo.GetByID(ctx, tenantID, id)
}

func (s *PermissionService) ListByTenant(ctx context.Context, tenantID string) ([]*domain.Permission, error) {
	return s.permissionRepo.ListByTenant(ctx, tenantID)
}

func (s *PermissionService) Update(ctx context.Context, tenantID string, permission *domain.Permission) error {
	if permission == nil {
		return errors.New("permission is nil")
	}
	if permission.ID == "" {
		return errors.New("permission id is empty")
	}
	existing, err := s.permissionRepo.GetByName(ctx, tenantID, permission.Name)
	if err != nil && !errors.Is(err, domain.ErrPermissionNotFound) {
		return err
	}
	if existing != nil && existing.ID != permission.ID {
		return domain.ErrPermissionAlreadyExists
	}
	return s.permissionRepo.Update(ctx, tenantID, permission)
}

func (s *PermissionService) SoftDelete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("permission id cannot be empty")
	}
	return s.permissionRepo.SoftDelete(ctx, tenantID, id)
}

func (s *PermissionService) Restore(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("permission id cannot be empty")
	}
	return s.permissionRepo.Restore(ctx, tenantID, id)
}
