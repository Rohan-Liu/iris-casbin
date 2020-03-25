package dtos

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	zh_translations "github.com/go-playground/validator/translations/zh"
	"reflect"
)

type ReqUser struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
	RoleIds  []uint `json:"role_ids"  validate:"required" comment:"角色"`
}

type ReqLogin struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

type ReqRole struct {
	Name           string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName    string `json:"display_name" comment:"显示名称"`
	Description    string `json:"description" comment:"描述"`
	PermissionsIds []uint `json:"permissions_ids" comment:"权限"`
}

type ReqPermission struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `json:"display_name" comment:"显示名称"`
	Description string `json:"description" comment:"描述"`
	Action      string `json:"action" comment:"Action"`
}

var (
	uni           *ut.UniversalTranslator
	Validate      *validator.Validate
	ValidateTrans ut.Translator
)

func init() {
	zh2 := zh.New()
	uni = ut.New(zh2, zh2)
	ValidateTrans, _ = uni.GetTranslator("zh")
	Validate = validator.New()
	// 收集结构体中的comment标签，用于替换英文字段名称，这样返回错误就能展示中文字段名称了
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
	if err := zh_translations.RegisterDefaultTranslations(Validate, ValidateTrans); err != nil {
		fmt.Println(fmt.Sprintf("RegisterDefaultTranslations %v", err))
	}
}
