package datamodels

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Permission struct {
	gorm.Model
	Name        string  `gorm:"not null VARCHAR(191)"`
	DisplayName string  `gorm:"not null VARCHAR(191)"`
	Description string  `gorm:"not null VARCHAR(191)"`
	Roles       []*Role `gorm:"many2many:role_permissions;"`
	Action      string  `gorm:"not null VARCHAR(191)"`
}

func (Permission) TableName() string {
	return "permissions"
}

func NewPermission(id uint, name, displayName, description, action string) *Permission {
	return &Permission{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        name,
		Action:      action,
		Description: description,
		DisplayName: displayName,
	}
}
