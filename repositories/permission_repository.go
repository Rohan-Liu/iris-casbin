package repositories

import (
	"../datamodels"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
)

type IPermissionRepository interface {
	IRepository
	FindByName(entity interface{}, name string) error
	ClearCasbin()
	RemoveCasbin(perm *datamodels.Permission) error
}

func PermissionRepository(db *gorm.DB, enforcer *casbin.Enforcer) IPermissionRepository {
	return &permissionRepository{
		baseRepository: baseRepository{
			db:       db,
			enforcer: enforcer,
		},
	}
}

type permissionRepository struct {
	baseRepository
}

func (r permissionRepository) FindByName(entity interface{}, name string) error {

	return r.FindOne(entity, "name = ?", name)
}

func (r permissionRepository) ClearCasbin() {
	r.enforcer.ClearPolicy()
}

func (r permissionRepository) RemoveCasbin(perm *datamodels.Permission) error {
	_, err := r.enforcer.DeletePermission(perm.Name, perm.Action)
	return err
}
