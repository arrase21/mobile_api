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

type FirePermissionActionRepo struct {
	client *firestore.Client
}

func NewFirePermissionActionRepo(client *firestore.Client) *FirePermissionActionRepo {
	return &FirePermissionActionRepo{client: client}
}

func (r *FirePermissionActionRepo) Create(ctx context.Context, tenantID string, action *domain.PermissionAction) error {
	if tenantID == "" || action == nil {
		return fmt.Errorf("tenantID and action cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").NewDoc()
	action.ID = ref.ID
	action.CreatedAt = time.Now().UTC()
	action.UpdatedAt = time.Now().UTC()
	action.IsDeleted = false
	action.TenantID = tenantID
	_, err := ref.Set(ctx, action)
	return err
}

func (r *FirePermissionActionRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.PermissionAction, error) {
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrPermissionActionNotFound
		}
		return nil, err
	}
	var action domain.PermissionAction
	if err := doc.DataTo(&action); err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *FirePermissionActionRepo) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.PermissionAction, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").
		Where("is_deleted", "==", false).
		OrderBy("created_at", firestore.Asc).
		Offset(offset).Limit(limit).
		Documents(ctx)
	defer iter.Stop()

	var actions []*domain.PermissionAction
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		a := new(domain.PermissionAction)
		if err := doc.DataTo(a); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *FirePermissionActionRepo) Update(ctx context.Context, tenantID string, action *domain.PermissionAction) error {
	if tenantID == "" || action == nil || action.ID == "" {
		return fmt.Errorf("tenantID, action, and action.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").Doc(action.ID)
	action.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "name", Value: action.Name},
		{Path: "description", Value: action.Description},
		{Path: "updated_at", Value: action.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionActionNotFound
		}
		return err
	}
	return nil
}

func (r *FirePermissionActionRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").Doc(id)
	now := time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: now},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionActionNotFound
		}
		return err
	}
	return nil
}

func (r *FirePermissionActionRepo) Restore(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("permission_actions").Doc(id)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now().UTC()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrPermissionActionNotFound
		}
		return err
	}
	return nil
}
