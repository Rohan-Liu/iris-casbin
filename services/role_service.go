package services

import (
	"../datamodels"
	"../repositories"
)

type IRoleService interface {
	GetAll() ([]*datamodels.Role, error)
	Get(id uint) (datamodels.Role, error)
	Delete(id uint) error
	Create(role *datamodels.Role) error
	Update(role *datamodels.Role, perms []*datamodels.Permission) error
	GetPage(name, orderBy string, pageIndex, pageSize uint64) ([]*datamodels.Role, uint64, error)
	GetByIds(ids []uint) ([]*datamodels.Role, error)
	Clear()
}

type roleService struct {
	repo repositories.IRoleRepository
}

func (r roleService) GetByIds(ids []uint) (outEntities []*datamodels.Role, err error) {
	err = r.repo.Find(&outEntities, "id IN (?)", ids)
	return
}

func RoleService(repo repositories.IRoleRepository) IRoleService {
	return &roleService{repo: repo}
}

func (r roleService) Clear() {
	_ = r.repo.Clear(&datamodels.Role{})
	r.repo.Exec("DELETE FROM role_permissions;")
}

func (r roleService) GetAll() (outEntities []*datamodels.Role, err error) {
	err = r.repo.FindAll(&outEntities)
	return
}

func (r roleService) Get(id uint) (outEntity datamodels.Role, err error) {
	err = r.repo.FindOne(&outEntity, "id = ?", id)
	return
}

func (r roleService) Create(entity *datamodels.Role) (err error) {
	err = r.repo.Create(entity)
	if err == nil {
		err = r.repo.UpdateCasbin(entity, entity.Permissions)
	}
	return
}

func (r roleService) Delete(id uint) error {
	role, err := r.Get(id)
	if err != nil {
		return err
	}
	err = r.repo.Delete(role)
	if err != nil {
		return err
	}
	r.repo.ClearAssociation(&role, "Permissions")
	r.repo.ClearAssociation(&role, "Users")
	return r.repo.RemoveCasbin(&role)
}

func (r roleService) Update(entity *datamodels.Role, perms []*datamodels.Permission) (err error) {
	err = r.repo.Update(entity)
	if err == nil {
		//r.repo.UpdateWithPermission(entity, perms)
		r.repo.ReplaceAssociation(&entity, "Permissions", perms)
		err = r.repo.UpdateCasbin(entity, perms)
	}
	return
}

func (r roleService) GetPage(name, orderBy string, pageIndex, pageSize uint64) (outEntities []*datamodels.Role, totalCount uint64, err error) {
	var query interface{}
	var arg string

	orderBys := make([]string, 0)
	switch orderBy {
	case "":
		orderBy = "id desc"
	case "name":
		orderBy = "name asc"
	}
	orderBys = append(orderBys, orderBy)
	if name != "" {
		query = "name like ?"
		arg = "%" + name + "%"

		totalCount, err = r.repo.FindPage(&outEntities, pageIndex, pageSize, orderBys, query, arg)
	} else {
		totalCount, err = r.repo.FindPage(&outEntities, pageIndex, pageSize, orderBys, query)
	}
	return
}
