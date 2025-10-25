package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInvalidTenant           = errors.New("invalid tenant")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrRoleNotFound            = errors.New("role not found")
	ErrRoleAlreadyExists       = errors.New("role already exists")
	ErrAssessmentNotFound      = errors.New("assessment not found")
	ErrSkinfoldNotFound        = errors.New("skinfold not found")
	ErrVolumetryNotFound       = errors.New("volumetry not found")
)

type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, tenantID string) (*Tenant, error)
	GetByName(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, tenantID string) error
}
type PermissionRepo interface {
	Create(ctx context.Context, tenantID string, permission *Permission) error
	GetByID(ctx context.Context, tenantID, id string) (*Permission, error)
	GetByName(ctx context.Context, tenantID, name string) (*Permission, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*Permission, error)
	Update(ctx context.Context, tenantID string, permission *Permission) error
	SoftDelete(ctx context.Context, tenantID, id string) error
	Restore(ctx context.Context, tenantID, id string) error
}

type PermissionActionRepo interface {
	Create(ctx context.Context, tenantID string, action *PermissionAction) error
	GetByID(ctx context.Context, tenantID, id string) (*PermissionAction, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*PermissionAction, error)
	Update(ctx context.Context, tenantID string, action *PermissionAction) error
	SoftDelete(ctx context.Context, tenantID, id string) error
	Restore(ctx context.Context, tenantID, id string) error
}

type RoleRepo interface {
	Create(ctx context.Context, tenantID string, role *Role) error
	GetByID(ctx context.Context, tenantID, id string) (*Role, error)
	GetByName(ctx context.Context, tenantID, name string) (*Role, error)
	ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*Role, error)
	Update(ctx context.Context, tenantID string, role *Role) error
	SoftDelete(ctx context.Context, tenantID, id string) error
	Restore(ctx context.Context, tenantID, id string) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, tenantID string, user *User) error
	GetByDni(ctx context.Context, tenantID, dni string) (*User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*User, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]*User, error)
	Update(ctx context.Context, tenantID string, user *User) error
	SoftDelete(ctx context.Context, tenantID, userID string) error
	Delete(ctx context.Context, tenantID, userID string) error
	Restore(ctx context.Context, tenantID, userID string) error
	ListDeleted(ctx context.Context, tenantID string) ([]*User, error)
}
