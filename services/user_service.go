package services

import (
	"../datamodels"
	"../repositories"
)

type IUserService interface {
	GetAll() ([]*datamodels.User, error)
	Get(id uint) (datamodels.User, error)
	Delete(id uint) error
	Create(user *datamodels.User) error
	Update(user *datamodels.User, roles []*datamodels.Role) error
	GetPage(name, orderBy string, pageIndex, pageSize uint64) ([]*datamodels.User, uint64, error)
	GetByName(name string) (datamodels.User, error)
	Clear()
}

type userService struct {
	repo repositories.IUserRepository
}

func UserService(repo repositories.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (r userService) Clear() {
	_ = r.repo.Clear(&datamodels.User{})
	r.repo.Exec("DELETE FROM user_roles;")
}

func (r userService) Get(id uint) (outEntity datamodels.User, err error) {
	err = r.repo.FindOne(&outEntity, "id = ?", id)
	return
}

func (r userService) GetAll() (outEntities []*datamodels.User, err error) {
	err = r.repo.FindAll(&outEntities)
	return
}

func (r userService) GetByName(name string) (outEntity datamodels.User, err error) {
	err = r.repo.FindOne(&outEntity, "username = ?", name)
	return
}

func (r userService) Create(entity *datamodels.User) (err error) {
	err = r.repo.Create(entity)
	if err == nil {
		err = r.repo.UpdateCasbin(entity)
	}
	return
}

func (r userService) Delete(id uint) error {
	user, err := r.Get(id)
	if err != nil {
		return err
	}
	err = r.repo.Delete(user)

	if err == nil {
		r.repo.ClearAssociation(&user, "Roles")
		err = r.repo.RemoveCasbin(&user)
	}
	return err
}

func (r userService) Update(entity *datamodels.User, roles []*datamodels.Role) (err error) {
	err = r.repo.Update(entity)
	if err != nil {
		return err
	}
	r.repo.ReplaceAssociation(&entity, "Roles", roles)
	err = r.repo.UpdateCasbin(entity)

	return
}

func (r userService) GetPage(name, orderBy string, pageIndex, pageSize uint64) (outEntities []*datamodels.User, totalCount uint64, err error) {
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
