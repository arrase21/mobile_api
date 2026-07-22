package domain

import "time"

type Tenant struct {
	ID         string    `firestore:"id" json:"id"`
	Name       string    `firestore:"name" json:"name"`
	Plan       string    `firestore:"plan" json:"plan"`
	Domain     string    `firestore:"domain" json:"domain"`
	IsActive   bool      `firestore:"is_active" json:"is_active"`
	MaxUsers   int       `firestore:"max_users" json:"max_users"`
	Expiration time.Time `firestore:"expiration" json:"expiration"`
	CreatedAt  time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at" json:"updated_at"`
}

type Permission struct {
	ID          string    `firestore:"id" json:"id"`
	Name        string    `firestore:"name" json:"name"` // ej: "user.manage"
	Description string    `firestore:"description" json:"description"`
	Actions     []string  `firestore:"actions" json:"actions"`
	IsDeleted   bool      `firestore:"is_deleted" json:"is_deleted"`
	DeletedAt   time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updated_at"`
	TenantID    string    `firestore:"tenant_id" json:"tenant_id"`
}

type PermissionAction struct {
	ID          string    `firestore:"id" json:"id"`
	Name        string    `firestore:"name" json:"name"`
	Description string    `firestore:"description" json:"description"`
	IsDeleted   bool      `firestore:"is_deleted" json:"is_deleted"`
	DeletedAt   time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updated_at"`
	TenantID    string    `firestore:"tenant_id" json:"tenant_id"`
}

type Role struct {
	ID          string    `firestore:"id" json:"id"`
	Name        string    `firestore:"name" json:"name"`
	Description string    `firestore:"description" json:"description"`
	Permissions []string  `firestore:"permissions" json:"permissions"`
	IsDeleted   bool      `firestore:"is_deleted" json:"is_deleted"`
	DeletedAt   time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updated_at"`
	TenantID    string    `firestore:"tenant_id" json:"tenant_id"`
}

type User struct {
	ID           string    `firestore:"id" json:"id"`
	FirebaseUID  string    `firestore:"firebase_uid" json:"firebase_uid"`
	FirstName    string    `firestore:"first_name" json:"first_name"`
	LastName     string    `firestore:"last_name" json:"last_name"`
	Dni          string    `firestore:"dni" json:"dni"`
	Gender       string    `firestore:"gender" json:"gender"`
	Phone        string    `firestore:"phone" json:"phone"`
	Email        string    `firestore:"email" json:"email"`
	DateBirth    time.Time `firestore:"date_birth" json:"date_birth"`
	Nickname     string    `firestore:"nickname" json:"nickname"`
	IsSuperAdmin bool      `firestore:"is_super_admin" json:"is_super_admin"`
	IsDeleted    bool      `firestore:"is_deleted" json:"is_deleted"`
	DeletedAt    time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt    time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at" json:"updated_at"`
	TenantID     string    `firestore:"tenant_id" json:"tenant_id"`
	RoleID       string    `firestore:"role_id" json:"role_id"`
	CreatedBy    string    `firestore:"created_by" json:"created_by"`
}

type Assessment struct {
	ID        string    `firestore:"id" json:"id"`
	UserID    string    `firestore:"user_id" json:"user_id"`
	Date      time.Time `firestore:"date" json:"date"`
	Height    float64   `firestore:"height" json:"height"`
	Weight    float64   `firestore:"weight" json:"weight"`
	Humerus   float64   `firestore:"humerus" json:"humerus"`
	Femur     float64   `firestore:"femur" json:"femur"`
	DeletedAt time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}

type Skinfold struct {
	ID            string    `firestore:"id" json:"id"`
	AssessmentID  string    `firestore:"assessment_id" json:"assessment_id"`
	Tricipital    float64   `firestore:"tricipital" json:"tricipital"`
	Subscapular   float64   `firestore:"subscapular" json:"subscapular"`
	Suprailiaco   float64   `firestore:"suprailiaco" json:"suprailiaco"`
	Abdominal     float64   `firestore:"abdominal" json:"abdominal"`
	Thighs        float64   `firestore:"thighs" json:"thighs"`
	Leg           float64   `firestore:"leg" json:"leg"`
	IsDeleted     bool      `firestore:"is_deleted" json:"is_deleted"`
	DeletedAt     time.Time `firestore:"deleted_at" json:"deleted_at"`
	CreatedAt     time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt     time.Time `firestore:"updated_at" json:"updated_at"`
	TenantID      string    `firestore:"tenant_id" json:"tenant_id"`
}
