package repositories

import (
	"../datamodels"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	IRepository
	FindByName(entity interface{}, name string) error
	UpdateCasbin(user *datamodels.User) error
	RemoveCasbin(user *datamodels.User) error
	UpdateWithRoles(user *datamodels.User, roles []*datamodels.Role)
}

func UserRepository(db *gorm.DB, enforcer *casbin.Enforcer) IUserRepository {
	return &userRepository{
		baseRepository: baseRepository{
			db:       db,
			enforcer: enforcer,
		},
	}
}

type userRepository struct {
	baseRepository
}

func (r userRepository) FindByName(entity interface{}, name string) error {

	return r.FindOne(entity, "name = ?", name)
}

func (r userRepository) UpdateWithRoles(user *datamodels.User, roles []*datamodels.Role) {
	r.db.Model(&user).Association("Roles").Replace(roles)
}

func (r userRepository) RemoveCasbin(user *datamodels.User) (err error) {

	_, err = r.enforcer.DeleteUser(user.GetCasbinName())
	if err != nil {
		return err
	}
	_, err = r.enforcer.DeleteRolesForUser(user.GetCasbinName())

	return err
}

func (r userRepository) UpdateCasbin(user *datamodels.User) (err error) {

	_, err = r.enforcer.DeleteRolesForUser(user.GetCasbinName())

	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		_, err = r.enforcer.AddRoleForUser(user.GetCasbinName(), role.GetCasbinName())
		if err != nil {
			return err
		}
	}
	return err
}
