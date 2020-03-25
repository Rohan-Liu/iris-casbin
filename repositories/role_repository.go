package repositories

import (
	"../datamodels"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
)

type IRoleRepository interface {
	IRepository
	FindByName(entity interface{}, name string) error
	RemoveCasbin(role *datamodels.Role) error
	UpdateCasbin(role *datamodels.Role, perms []*datamodels.Permission) error
	UpdateWithPermission(role *datamodels.Role, perms []*datamodels.Permission)
}

func RoleRepository(db *gorm.DB, enforcer *casbin.Enforcer) IRoleRepository {
	return &roleRepository{
		baseRepository: baseRepository{
			db:       db,
			enforcer: enforcer,
		},
	}
}

type roleRepository struct {
	baseRepository
}

func (r roleRepository) UpdateWithPermission(role *datamodels.Role, perms []*datamodels.Permission) {
	r.db.Model(&role).Association("Permissions").Replace(perms)

}

func (r roleRepository) FindByName(entity interface{}, name string) error {

	return r.FindOne(entity, "name = ?", name)
}

func (r roleRepository) RemoveCasbin(role *datamodels.Role) (err error) {

	_, err = r.enforcer.DeleteRole(role.GetCasbinName())
	if err != nil {
		return err
	}
	_, err = r.enforcer.DeletePermissionForUser(role.GetCasbinName())
	return err
}

func (r roleRepository) UpdateCasbin(role *datamodels.Role, perms []*datamodels.Permission) (err error) {

	_, err = r.enforcer.DeletePermissionsForUser(role.GetCasbinName())

	if err != nil {
		return err
	}

	for _, perm := range perms {
		_, err = r.enforcer.AddPermissionForUser(role.GetCasbinName(), perm.Name, perm.Action)
		if err != nil {
			return err
		}
	}
	return err
}
