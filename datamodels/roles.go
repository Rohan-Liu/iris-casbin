package datamodels

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Role struct {
	gorm.Model
	Name        string        `gorm:"not null VARCHAR(191)"`
	DisplayName string        `gorm:"not null VARCHAR(191)"`
	Description string        `gorm:"not null VARCHAR(191)"`
	Users       []*User       `gorm:"many2many:user_roles;"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

func (role *Role) GetCasbinName() string {
	return "role:" + strconv.FormatUint(uint64(role.ID), 10)
}

func (Role) TableName() string {
	return "roles"
}
