package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/mobileapi/internal/domain"
)

type TenantService struct {
	tenantRepo domain.TenantRepo
}

func NewTenantService(tenantRepo domain.TenantRepo) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

func (s *TenantService) Create(ctx context.Context, tenant *domain.Tenant) error {
	if tenant == nil {
		return errors.New("tenant is nil")
	}
	tenant.Name = strings.TrimSpace(tenant.Name)
	tenant.Plan = strings.TrimSpace(tenant.Plan)
	tenant.Domain = strings.TrimSpace(tenant.Domain)
	if tenant.Name == "" {
		return errors.New("tenant name cannot be empty")
	}
	existing, err := s.tenantRepo.GetByName(ctx, tenant.Name)
	if err != nil && !errors.Is(err, domain.ErrInvalidTenant) {
		return err
	}
	if existing != nil {
		return errors.New("tenant with this name already exists")
	}
	return s.tenantRepo.Create(ctx, tenant)
}

func (s *TenantService) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	if tenantID == "" {
		return nil, errors.New("tenant id is empty")
	}
	return s.tenantRepo.GetByID(ctx, tenantID)
}

func (s *TenantService) List(ctx context.Context, offset, limit int) ([]*domain.Tenant, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.tenantRepo.List(ctx, offset, limit)
}

func (s *TenantService) Update(ctx context.Context, tenant *domain.Tenant) error {
	if tenant == nil {
		return errors.New("tenant is nil")
	}
	if tenant.ID == "" {
		return errors.New("tenant id is empty")
	}
	return s.tenantRepo.Update(ctx, tenant)
}

func (s *TenantService) Delete(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return errors.New("tenant id is empty")
	}
	return s.tenantRepo.Delete(ctx, tenantID)
}
