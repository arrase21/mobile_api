package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/mobileapi/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(a domain.UserRepo) *UserService {
	return &UserService{userRepo: a}
}

// Validates
func (s *UserService) checkUserExist(ctx context.Context, tenantID, email, dni string) error {
	if email != "" {
		if u, err := s.userRepo.GetByEmail(ctx, tenantID, email); err == nil && u != nil {
			return errors.New("user with this email already exist")
		} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}
	}
	if dni != "" {
		if u, err := s.userRepo.GetByDni(ctx, tenantID, dni); err == nil && u != nil {
			return errors.New("user with this dni already exist")
		} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}
	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, usr *domain.User) error {
	if usr == nil {
		return errors.New("user is nil")
	}
	// clean data
	usr.Email = strings.ToLower(strings.TrimSpace(usr.Email))
	usr.FirstName = strings.TrimSpace(usr.FirstName)
	usr.LastName = strings.TrimSpace(usr.LastName)
	usr.Dni = strings.TrimSpace(usr.Dni)
	// Validate is not empty
	if usr.FirstName == "" || usr.LastName == "" || usr.Dni == "" || usr.Email == "" {
		return errors.New("first name, last name, dni, and email cannot be empty")
	}
	// validate if exist
	if err := s.checkUserExist(ctx, usr.TenantID, usr.Email, usr.Dni); err != nil {
		return err
	}
	// Create user
	return s.userRepo.CreateUser(ctx, usr.TenantID, usr)
}

func (s *UserService) GetByDni(ctx context.Context, tenantID, dni string) (*domain.User, error) {
	return s.userRepo.GetByDni(ctx, tenantID, dni)
}

func (s *UserService) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, tenantID, email)
}

func (s *UserService) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, error) {
	return s.userRepo.List(ctx, tenantID, offset, limit)
}

func (s *UserService) Update(ctx context.Context, tenantID string, usr *domain.User) error {
	if usr == nil {
		return errors.New("user is nil")
	}
	return s.userRepo.Update(ctx, tenantID, usr)
}

func (s *UserService) SoftDelete(ctx context.Context, tenantID, userID string) error {
	if userID == "" {
		return errors.New("user id is empty")
	}
	return s.userRepo.SoftDelete(ctx, tenantID, userID)
}

func (s *UserService) Delete(ctx context.Context, tenantID, userID string) error {
	if userID == "" {
		return errors.New("user id is empty")
	}
	return s.userRepo.Delete(ctx, tenantID, userID)
}

func (s *UserService) Restore(ctx context.Context, tenantID, userID string) error {
	if userID == "" {
		return errors.New("user id is empty")
	}
	return s.userRepo.Restore(ctx, tenantID, userID)
}

func (s *UserService) ListDeleted(ctx context.Context, tenantID string) ([]*domain.User, error) {
	return s.userRepo.ListDeleted(ctx, tenantID)
}
