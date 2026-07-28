package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arrase21/mobileapi/internal/domain"
)

type FireUserRepo struct {
	client *firestore.Client
}

func NewFireUserRepo(client *firestore.Client) *FireUserRepo {
	return &FireUserRepo{client: client}
}

func (r *FireUserRepo) CreateUser(ctx context.Context, tenantID string, user *domain.User) error {
	if tenantID == "" || user == nil {
		return fmt.Errorf("tenantID and user cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("users").NewDoc()
	user.ID = ref.ID
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.IsDeleted = false
	_, err := ref.Set(ctx, user)
	return err
}

func (r *FireUserRepo) GetByDni(ctx context.Context, tenantID string, dni string) (*domain.User, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("users").
		Where("dni", "==", dni).
		Where("is_deleted", "==", false).
		Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *FireUserRepo) GetByID(ctx context.Context, tenantID, id string) (*domain.User, error) {
	doc, err := r.client.Collection("tenants").Doc(tenantID).Collection("users").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *FireUserRepo) GetByFirebaseUID(ctx context.Context, tenantID, firebaseUID string) (*domain.User, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("users").
		Where("firebase_uid", "==", firebaseUID).
		Where("is_deleted", "==", false).
		Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *FireUserRepo) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("users").
		Where("email", "==", email).Where("is_deleted", "==", false).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *FireUserRepo) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, error) {
	iter := r.client.Collection("tenants").Doc(tenantID).Collection("users").Where("is_deleted", "==", false).
		OrderBy("created_at", firestore.Asc).Offset(offset).Limit(limit).Documents(ctx)
	defer iter.Stop()

	var users []*domain.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		u := new(domain.User)
		if err := doc.DataTo(u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *FireUserRepo) Update(ctx context.Context, tenantID string, user *domain.User) error {
	if tenantID == "" || user == nil || user.ID == "" {
		return fmt.Errorf("tenantID, user, and user.ID cannot be empty")
	}
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("users").Doc(user.ID)
	user.UpdatedAt = time.Now().UTC()
	updates := []firestore.Update{
		{Path: "updated_at", Value: user.UpdatedAt},
	}

	if user.Email != "" {
		updates = append(updates, firestore.Update{Path: "email", Value: user.Email})
	}
	if user.FirstName != "" {
		updates = append(updates, firestore.Update{Path: "first_name", Value: user.FirstName})
	}
	if user.LastName != "" {
		updates = append(updates, firestore.Update{Path: "last_name", Value: user.LastName})
	}
	if user.Dni != "" {
		updates = append(updates, firestore.Update{Path: "dni", Value: user.Dni})
	}
	if user.Gender != "" {
		updates = append(updates, firestore.Update{Path: "gender", Value: user.Gender})
	}
	if user.Phone != "" {
		updates = append(updates, firestore.Update{Path: "phone", Value: user.Phone})
	}
	if !user.DateBirth.IsZero() {
		updates = append(updates, firestore.Update{Path: "date_birth", Value: user.DateBirth})
	}
	if user.Nickname != "" {
		updates = append(updates, firestore.Update{Path: "nickname", Value: user.Nickname})
	}
	if user.RoleID != "" {
		updates = append(updates, firestore.Update{Path: "role_id", Value: user.RoleID})
	}
	if user.FirebaseUID != "" {
		updates = append(updates, firestore.Update{Path: "firebase_uid", Value: user.FirebaseUID})
	}
	_, err := ref.Update(ctx, updates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *FireUserRepo) SoftDelete(ctx context.Context, tenantID, userID string) error {
	docRef := r.client.Collection("tenants").Doc(tenantID).Collection("users").Doc(userID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: true},
		{Path: "deleted_at", Value: time.Now()},
		{Path: "updated_at", Value: time.Now()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to soft delete user: %w", err)
	}
	return nil
}

func (r *FireUserRepo) Delete(ctx context.Context, tenantID, userID string) error {
	if tenantID == "" || userID == "" {
		return fmt.Errorf("tenantID and userID cannot be empty")
	}
	_, err := r.client.Collection("tenants").Doc(tenantID).Collection("users").Doc(userID).Delete(ctx)
	return err
}

func (r *FireUserRepo) Restore(ctx context.Context, tenantID, userID string) error {
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("users").Doc(userID)
	_, err := ref.Update(ctx, []firestore.Update{
		{Path: "is_deleted", Value: false},
		{Path: "deleted_at", Value: firestore.Delete},
		{Path: "updated_at", Value: time.Now()},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *FireUserRepo) ListDeleted(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, error) {
	ref := r.client.Collection("tenants").Doc(tenantID).Collection("users").
		Where("is_deleted", "==", true).OrderBy("deleted_at", firestore.Desc).
		Offset(offset).Limit(limit).
		Documents(ctx)

	defer ref.Stop()

	var users []*domain.User
	for {
		doc, err := ref.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		u := new(domain.User)
		if err := doc.DataTo(u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *FireUserRepo) CountByTenant(ctx context.Context, tenantID string) (int, error) {
	if tenantID == "" {
		return 0, fmt.Errorf("tenantID cannot be empty")
	}
	q := r.client.Collection("tenants").Doc(tenantID).Collection("users").
		Where("is_deleted", "==", false)
	aggregationQuery := q.NewAggregationQuery().WithCount("all")

	results, err := aggregationQuery.Get(ctx)
	if err != nil {
		return 0, err
	}

	count, ok := results["all"]
	if !ok {
		return 0, fmt.Errorf("count not found in aggregation results")
	}
	return int(count.(int64)), nil
}
