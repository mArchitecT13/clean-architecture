package entities

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity in the domain
type User struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Email     string         `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user instance
func NewUser(email, name string) *User {
	now := time.Now()
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}
