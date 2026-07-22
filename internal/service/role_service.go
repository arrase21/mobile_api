package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/mobileapi/internal/domain"
)

type RoleService struct {
	roleRepo domain.RoleRepo
}

func NewRoleService(roleRepo domain.RoleRepo) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) CreateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}
	role.Name = strings.TrimSpace(role.Name)
	role.Description = strings.TrimSpace(role.Description)
	if role.Name == "" {
		return errors.New("role cannot be empty")
	}
	existing, err := s.roleRepo.GetByName(ctx, role.TenantID, role.Name)
	if err != nil && !errors.Is(err, domain.ErrRoleNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrRoleAlreadyExists
	}
	return s.roleRepo.Create(ctx, role.TenantID, role)
}

func (s *RoleService) GetByID(ctx context.Context, tenantID, id string) (*domain.Role, error) {
	return s.roleRepo.GetByID(ctx, tenantID, id)
}

func (s *RoleService) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Role, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.roleRepo.ListByTenant(ctx, tenantID, offset, limit)
}

func (s *RoleService) Update(ctx context.Context, tenantID string, role *domain.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}
	if role.ID == "" {
		return errors.New("role id is empty")
	}
	role.Name = strings.TrimSpace(role.Name)
	role.Description = strings.TrimSpace(role.Description)

	existing, err := s.roleRepo.GetByName(ctx, tenantID, role.Name)
	if err != nil && !errors.Is(err, domain.ErrRoleNotFound) {
		return err
	}
	if existing != nil && existing.ID != role.ID {
		return domain.ErrRoleAlreadyExists
	}

	return s.roleRepo.Update(ctx, tenantID, role)
}

func (s *RoleService) SoftDelete(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("role id cannot be empty")
	}
	return s.roleRepo.SoftDelete(ctx, tenantID, id)
}

func (s *RoleService) Restore(ctx context.Context, tenantID, id string) error {
	if id == "" {
		return errors.New("role id cannot be empty")
	}
	return s.roleRepo.Restore(ctx, tenantID, id)
}
