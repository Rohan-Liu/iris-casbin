package services

import (
	"../datamodels"
	"../repositories"
)

type IPermissionService interface {
	GetAll() ([]*datamodels.Permission, error)
	Get(id uint) (*datamodels.Permission, error)
	Delete(id uint) error
	Create(permission *datamodels.Permission) error
	Update(permission *datamodels.Permission) error
	GetPage(name, orderBy string, pageIndex, pageSize uint64) ([]*datamodels.Permission, uint64, error)
	GetByIds(ids []uint) ([]*datamodels.Permission, error)
	Clear()
}

type permissionService struct {
	repo repositories.IPermissionRepository
}

func (r permissionService) Clear() {
	_ = r.repo.Clear(&datamodels.Permission{})
	r.repo.ClearCasbin()
}

func (r permissionService) GetByIds(ids []uint) (outEntities []*datamodels.Permission, err error) {
	err = r.repo.Find(&outEntities, "id IN (?)", ids)

	return
}

func PermissionService(repo repositories.IPermissionRepository) IPermissionService {
	return &permissionService{repo: repo}
}

func (r permissionService) GetAll() (outEntities []*datamodels.Permission, err error) {
	err = r.repo.FindAll(&outEntities)
	return
}

func (r permissionService) Get(id uint) (outEntity *datamodels.Permission, err error) {
	outEntity = &datamodels.Permission{}
	err = r.repo.FindOne(outEntity, "id = ?", id)
	return
}

func (r permissionService) Create(entity *datamodels.Permission) (err error) {
	err = r.repo.Create(entity)
	return
}

func (r permissionService) Delete(id uint) error {
	permission, err := r.Get(id)
	if err != nil {
		return err
	}
	err = r.repo.Delete(permission)
	r.repo.ClearAssociation(&permission, "Roles")
	if err != nil {
		return err
	}
	return r.repo.RemoveCasbin(permission)

}

func (r permissionService) Update(entity *datamodels.Permission) (err error) {
	err = r.repo.Update(entity)
	return
}

func (r permissionService) GetPage(name, orderBy string, pageIndex, pageSize uint64) (outEntities []*datamodels.Permission, totalCount uint64, err error) {
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
