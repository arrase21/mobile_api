package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/arrase21/mobileapi/internal/domain"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FireRoleRepo struct {
	client *firestore.Client
}

func NewFireRoleRepo(client *firestore.Client) *FireRoleRepo {
	return &FireRoleRepo{client: client}
}

func (r *FireRoleRepo) Create(ctx context.Context, tenantID string, role *domain.Role) error {
	if tenantID == "" || role == nil {
		return errors.New("tenant id and role cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("roles").NewDoc()
	role.ID = ref.ID
	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()
	role.IsDeleted = false
	role.TenantID = tenantID

	_, err := ref.Set(ctx, role)
	return err
}

func (r *FireRoleRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.Role, error) {
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("roles").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrRoleNotFound
		}
		return nil, err
	}

	var role domain.Role
	if err := doc.DataTo(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *FireRoleRepo) GetByName(ctx context.Context, tenantID, name string) (*domain.Role, error) {
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("roles").Where("name", "==", name).
		Where("is_deleted", "==", false).Limit(1).Documents(ctx)
	defer ref.Stop()
	doc, err := ref.Next()
	if err == iterator.Done {
		return nil, domain.ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	var role domain.Role
	if err := doc.DataTo(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *FireRoleRepo) ListByTenant(ctx context.Context, tenanID string, offset, limit int) ([]*domain.Role, error) {
	ref := r.client.Collection("tenants").Doc(tenanID).Collection("roles").Where("is_deleted", "==", false).
		OrderBy("created_at", firestore.Asc).Offset(offset).Limit(limit).Documents(ctx)
	defer ref.Stop()
	var roles []*domain.Role
	for {
		doc, err := ref.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		role := new(domain.Role)
		if err := doc.DataTo(role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *FireRoleRepo) Update(ctx context.Context, tenantID string, role *domain.Role) error {
	if tenantID == "" || role == nil || role.ID == "" {
		return fmt.Errorf("TennatID, role and role.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("roles").Doc(role.ID)
	role.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "name", Value: role.Name},
		{Path: "description", Value: role.Description},
		{Path: "permissions", Value: role.Permissions},
		{Path: "updated_at", Value: role.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrRoleNotFound
		}
		return err
	}
	return nil
}

func (r *FireRoleRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("TenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("roles").Doc(id)
	now := time.Now().UTC()

	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: now},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrRoleNotFound
		}
		return err
	}
	return nil
}

func (r *FireRoleRepo) Restore(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("roles").Doc(id)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now().UTC()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrRoleNotFound
		}
		return err
	}
	return nil
}
