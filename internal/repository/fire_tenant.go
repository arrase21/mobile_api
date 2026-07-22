package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/arrase21/mobileapi/internal/domain"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FireTenantRepo struct {
	client *firestore.Client
}

func NewFireTenantRepo(client *firestore.Client) *FireTenantRepo {
	return &FireTenantRepo{client: client}
}

func (r *FireTenantRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	if tenant == nil {
		return fmt.Errorf("tenant cannot be nil")
	}
	ref := r.client.Collection("tenants").Doc(tenant.ID)
	doc, err := ref.Get(ctx)
	if err == nil && doc.Exists() {
		return fmt.Errorf("tenant with id %s already exists", tenant.ID)
	}

	if tenant.ID == "" {
		ref = r.client.Collection("tenants").NewDoc()
		tenant.ID = ref.ID
	}
	tenant.CreatedAt = time.Now().UTC()
	tenant.UpdatedAt = time.Now().UTC()
	_, err = ref.Set(ctx, tenant)
	return err
}

func (r *FireTenantRepo) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenantID cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrInvalidTenant
		}
		return nil, err
	}
	var tenant domain.Tenant
	if err := doc.DataTo(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *FireTenantRepo) GetByName(ctx context.Context, name string) (*domain.Tenant, error) {
	iter := r.client.Collection("tenants").
		Where("name", "==", name).
		Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrInvalidTenant
	}
	if err != nil {
		return nil, err
	}
	var tenant domain.Tenant
	if err := doc.DataTo(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *FireTenantRepo) List(ctx context.Context) ([]*domain.Tenant, error) {
	iter := r.client.Collection("tenants").
		OrderBy("created_at", firestore.Asc).
		Documents(ctx)
	defer iter.Stop()

	var tenants []*domain.Tenant
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		t := new(domain.Tenant)
		if err := doc.DataTo(t); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}

func (r *FireTenantRepo) Update(ctx context.Context, tenant *domain.Tenant) error {
	if tenant == nil || tenant.ID == "" {
		return fmt.Errorf("tenant and tenant.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenant.ID)
	tenant.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "name", Value: tenant.Name},
		{Path: "plan", Value: tenant.Plan},
		{Path: "domain", Value: tenant.Domain},
		{Path: "is_active", Value: tenant.IsActive},
		{Path: "max_users", Value: tenant.MaxUsers},
		{Path: "expiration", Value: tenant.Expiration},
		{Path: "updated_at", Value: tenant.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrInvalidTenant
		}
		return err
	}
	return nil
}

func (r *FireTenantRepo) Delete(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenantID cannot be empty")
	}
	_, err := r.client.Collection("tenants").Doc(tenantID).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrInvalidTenant
		}
		return err
	}
	return nil
}
