package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// Create creates a new user
func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

// GetAll retrieves all users
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

// GetByID retrieves a user by ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
}

// Delete deletes a user
func (u *User) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}
