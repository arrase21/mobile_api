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

type FirePermissionRepo struct {
	client *firestore.Client
}

func NewFirePermissionRepo(client *firestore.Client) *FirePermissionRepo {
	return &FirePermissionRepo{client: client}
}

func (r *FirePermissionRepo) Create(ctx context.Context, tenantID string, permission *domain.Permission) error {
	if tenantID == "" || permission == nil {
		return fmt.Errorf("tenantID and permission cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").NewDoc()
	permission.ID = ref.ID
	permission.CreatedAt = time.Now().UTC()
	permission.UpdatedAt = time.Now().UTC()
	permission.IsDeleted = false
	permission.TenantID = tenantID
	_, err := ref.Set(ctx, permission)
	return err
}

func (r *FirePermissionRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.Permission, error) {
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, err
	}
	var permission domain.Permission
	if err := doc.DataTo(&permission); err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *FirePermissionRepo) GetByName(ctx context.Context, tenantID, name string) (*domain.Permission, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").
		Where("name", "==", name).
		Where("is_deleted", "==", false).
		Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrPermissionNotFound
	}
	if err != nil {
		return nil, err
	}
	var permission domain.Permission
	if err := doc.DataTo(&permission); err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *FirePermissionRepo) ListByTenant(ctx context.Context, tenantID string) ([]*domain.Permission, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").
		Where("is_deleted", "==", false).
		OrderBy("created_at", firestore.Asc).
		Documents(ctx)
	defer iter.Stop()

	var permissions []*domain.Permission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		p := new(domain.Permission)
		if err := doc.DataTo(p); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

func (r *FirePermissionRepo) Update(ctx context.Context, tenantID string, permission *domain.Permission) error {
	if tenantID == "" || permission == nil || permission.ID == "" {
		return fmt.Errorf("tenantID, permission, and permission.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").Doc(permission.ID)
	permission.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "name", Value: permission.Name},
		{Path: "description", Value: permission.Description},
		{Path: "actions", Value: permission.Actions},
		{Path: "updated_at", Value: permission.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionNotFound
		}
		return err
	}
	return nil
}

func (r *FirePermissionRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").Doc(id)
	now := time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: now},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionNotFound
		}
		return err
	}
	return nil
}

func (r *FirePermissionRepo) Restore(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permissions").Doc(id)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now().UTC()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionNotFound
		}
		return err
	}
	return nil
}
