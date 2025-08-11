package models

import (
	"time"

	"gorm.io/gorm"
)

type {{.ModelName}} struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func ({{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}

// Create creates a new {{.ModelName}}
func (m *{{.ModelName}}) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// GetAll retrieves all {{.ModelName}}s
func GetAll{{.ModelName}}s(db *gorm.DB) ([]{{.ModelName}}, error) {
	var models []{{.ModelName}}
	err := db.Find(&models).Error
	return models, err
}

// GetByID retrieves a {{.ModelName}} by ID
func Get{{.ModelName}}ByID(db *gorm.DB, id uint) (*{{.ModelName}}, error) {
	var model {{.ModelName}}
	err := db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// Update updates a {{.ModelName}}
func (m *{{.ModelName}}) Update(db *gorm.DB) error {
	return db.Save(m).Error
}

// Delete deletes a {{.ModelName}}
func (m *{{.ModelName}}) Delete(db *gorm.DB) error {
	return db.Delete(m).Error
}
