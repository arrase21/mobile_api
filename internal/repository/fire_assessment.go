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

type FireAssessmentRepo struct {
	client *firestore.Client
}

func NewFireAssessmentRepo(client *firestore.Client) *FireAssessmentRepo {
	return &FireAssessmentRepo{client: client}
}

func (r *FireAssessmentRepo) Create(ctx context.Context, tenantID string, assessment *domain.Assessment) error {
	if tenantID == "" || assessment == nil {
		return fmt.Errorf("tenantID and assessment cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").NewDoc()
	assessment.ID = ref.ID
	assessment.IsDeleted = false
	assessment.CreatedAt = time.Now().UTC()
	assessment.UpdatedAt = time.Now().UTC()
	_, err := ref.Set(ctx, assessment)
	return err
}

func (r *FireAssessmentRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.Assessment, error) {
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id cannot be empty")
	}
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrAssessmentNotFound
		}
		return nil, err
	}
	var assessment domain.Assessment
	if err := doc.DataTo(&assessment); err != nil {
		return nil, err
	}
	return &assessment, nil
}

func (r *FireAssessmentRepo) ListByUser(ctx context.Context, tenantID, userID string, offset, limit int) ([]*domain.Assessment, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").
		Where("user_id", "==", userID).
		Where("is_deleted", "==", false).
		OrderBy("date", firestore.Desc).
		Offset(offset).Limit(limit).
		Documents(ctx)
	defer iter.Stop()

	var assessments []*domain.Assessment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		a := new(domain.Assessment)
		if err := doc.DataTo(a); err != nil {
			return nil, err
		}
		assessments = append(assessments, a)
	}
	return assessments, nil
}

func (r *FireAssessmentRepo) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Assessment, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").
		Where("is_deleted", "==", false).
		OrderBy("date", firestore.Desc).
		Offset(offset).Limit(limit).
		Documents(ctx)
	defer iter.Stop()

	var assessments []*domain.Assessment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		a := new(domain.Assessment)
		if err := doc.DataTo(a); err != nil {
			return nil, err
		}
		assessments = append(assessments, a)
	}
	return assessments, nil
}

func (r *FireAssessmentRepo) Update(ctx context.Context, tenantID string, assessment *domain.Assessment) error {
	if tenantID == "" || assessment == nil || assessment.ID == "" {
		return fmt.Errorf("tenantID, assessment, and assessment.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").Doc(assessment.ID)
	assessment.UpdatedAt = time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "date", Value: assessment.Date},
		{Path: "height", Value: assessment.Height},
		{Path: "weight", Value: assessment.Weight},
		{Path: "humerus", Value: assessment.Humerus},
		{Path: "femur", Value: assessment.Femur},
		{Path: "updated_at", Value: assessment.UpdatedAt},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrAssessmentNotFound
		}
		return err
	}
	return nil
}

func (r *FireAssessmentRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").Doc(id)
	now := time.Now().UTC()
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: now},
		{Path: "updated_at", Value: now},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrAssessmentNotFound
		}
		return err
	}
	return nil
}

func (r *FireAssessmentRepo) Restore(ctx context.Context, tenantID, id string) error {
	if tenantID == "" || id == "" {
		return fmt.Errorf("tenantID and id cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").Doc(id)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now().UTC()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrAssessmentNotFound
		}
		return err
	}
	return nil
}

func (r *FireAssessmentRepo) Delete(ctx context.Context, tenantID, assessmentID string) error {
	if tenantID == "" || assessmentID == "" {
		return fmt.Errorf("tenantID and assessmentID cannot be empty")
	}
	_, err := r.client.Collection("tenants").Doc(tenantID).Collection("assessments").Doc(assessmentID).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrAssessmentNotFound
		}
		return err
	}
	return nil
}
