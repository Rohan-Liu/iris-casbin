package datamodels

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

type User struct {
	gorm.Model
	Name     string  `gorm:"not null VARCHAR(191)"`
	Username string  `gorm:"not null VARCHAR(191)"`
	Password string  `gorm:"not null VARCHAR(191)"`
	Roles    []*Role `gorm:"many2many:user_roles;"`
}

func (user *User) GetCasbinName() string {
	return "user:" + strconv.FormatUint(uint64(user.ID), 10)
}

func (user *User) CheckPassword(password string) bool {
	plainPwd := []byte(password)
	byteHash := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	password := hashAndSalt(user.Password)

	return scope.SetColumn("Password", password)
}

func (user User) TableName() string {
	return "users"
}

func hashAndSalt(password string) string {

	pwd := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
