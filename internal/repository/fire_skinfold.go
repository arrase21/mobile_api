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

type FireSkinfoldRepo struct {
	client *firestore.Client
}

func NewFireSkinfoldRepo(client *firestore.Client) *FireSkinfoldRepo {
	return &FireSkinfoldRepo{client: client}
}

func (r *FireSkinfoldRepo) Create(ctx context.Context, tenantID string, skinfold *domain.Skinfold) error {
	if tenantID == "" || skinfold == nil {
		return fmt.Errorf("tenantID and skinfold cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").NewDoc()
	skinfold.ID = ref.ID
	skinfold.CreatedAt = time.Now().UTC()
	skinfold.UpdatedAt = time.Now().UTC()
	skinfold.IsDeleted = false
	skinfold.TenantID = tenantID
	_, err := ref.Set(ctx, skinfold)
	return err
}

func (r *FireSkinfoldRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.Skinfold, error) {
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrSkinfoldNotFound
		}
		return nil, err
	}
	var skinfold domain.Skinfold
	if err := doc.DataTo(&skinfold); err != nil {
		return nil, err
	}
	return &skinfold, nil
}

func (r *FireSkinfoldRepo) ListByAssessment(ctx context.Context, tenantID, assessmentID string, offset, limit int) ([]*domain.Skinfold, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").
		Where("assessment_id", "==", assessmentID).
		Where("is_deleted", "==", false).
		OrderBy("created_at", firestore.Asc).
		Offset(offset).Limit(limit).
		Documents(ctx)
	defer iter.Stop()

	var skinfolds []*domain.Skinfold
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		s := new(domain.Skinfold)
		if err := doc.DataTo(s); err != nil {
			return nil, err
		}
		skinfolds = append(skinfolds, s)
	}
	return skinfolds, nil
}

func (r *FireSkinfoldRepo) Update(ctx context.Context, tenantID string, skinfold *domain.Skinfold) error {
	if tenantID == "" || skinfold == nil || skinfold.ID == "" {
		return fmt.Errorf("tenantID, skinfold, and skinfold.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").Doc(skinfold.ID)
	skinfold.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "tricipital", Value: skinfold.Tricipital},
		{Path: "subscapular", Value: skinfold.Subscapular},
		{Path: "suprailiaco", Value: skinfold.Suprailiaco},
		{Path: "abdominal", Value: skinfold.Abdominal},
		{Path: "thighs", Value: skinfold.Thighs},
		{Path: "leg", Value: skinfold.Leg},
		{Path: "updated_at", Value: skinfold.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrSkinfoldNotFound
		}
		return err
	}
	return nil
}

func (r *FireSkinfoldRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").Doc(id)
	now := time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: now},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrSkinfoldNotFound
		}
		return err
	}
	return nil
}

func (r *FireSkinfoldRepo) Restore(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").Doc(id)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now().UTC()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrSkinfoldNotFound
		}
		return err
	}
	return nil
}

func (r *FireSkinfoldRepo) Delete(ctx context.Context, tenantID, skinfoldID string) error {
	if tenantID == "" || skinfoldID == "" {
		return fmt.Errorf("tenantID and skinfoldID cannot be empty")
	}
	_, err := r.client.Collection("tenants").Doc(tenantID).Collection("skinfolds").Doc(skinfoldID).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrSkinfoldNotFound
		}
		return err
	}
	return nil
}
