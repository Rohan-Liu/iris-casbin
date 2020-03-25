package controllers

import (
	"../../datamodels"
	"../../services"
	"../dtos"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type RoleController struct {
	Service     services.IRoleService
	PermService services.IPermissionService
	UserService services.IUserService
}

func (r *RoleController) Get(ctx iris.Context) {
	result, _ := r.Service.GetAll()

	_, _ = ctx.JSON(ToSuccess(rolesToResp(result), "操作成功"))
}

func (r *RoleController) Post(ctx iris.Context) {
	req := new(dtos.ReqRole)

	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	err := dtos.Validate.Struct(*req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(dtos.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ToFailure(e))
				return
			}
		}
	}

	perms, _ := r.PermService.GetByIds(req.PermissionsIds)

	role := &datamodels.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Permissions: perms,
	}

	err = r.Service.Create(role)

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if role.ID == 0 {
		_, _ = ctx.JSON(ToFailure("操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ToSuccess(nil, "操作成功"))
		return
	}
}

func (r *RoleController) GetBy(id uint, ctx iris.Context) {
	role, err := r.Service.Get(id)
	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)

	_, _ = ctx.JSON(ToSuccess(roleToResp(&role), "操作成功"))
}

func (r *RoleController) DeleteBy(ctx iris.Context, id uint) {
	err := r.Service.Delete(id)
	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)

	_, _ = ctx.JSON(ToSuccess(nil, "删除成功"))
}

func (r *RoleController) PutBy(ctx iris.Context, id uint) {
	req := new(dtos.ReqRole)

	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	err := dtos.Validate.Struct(*req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(dtos.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ToFailure(e))
				return
			}
		}
	}

	perms, _ := r.PermService.GetByIds(req.PermissionsIds)

	role := &datamodels.Role{
		Model: gorm.Model{
			ID: id,
		},
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	}

	err = r.Service.Update(role, perms)

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	_, _ = ctx.JSON(ToSuccess(nil, "操作成功"))
}

func (r *RoleController) GetPage(ctx iris.Context) {
	pageIndex := ctx.Values().GetUintDefault("pageIndex", 1)
	pageSize := ctx.Values().GetUintDefault("pageSize", 20)
	name := ctx.Values().GetStringDefault("name", "")
	orderBy := ctx.Values().GetStringDefault("orderBy", "")

	result, totalCount, err := r.Service.GetPage(name, orderBy, uint64(pageIndex), uint64(pageSize))

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	list := &Lists{
		Items:      rolesToResp(result),
		TotalCount: totalCount,
	}
	_, _ = ctx.JSON(ToSuccess(list, "操作成功"))

}

func rolesToResp(roles []*datamodels.Role) []*dtos.RespRole {
	var rs []*dtos.RespRole
	for _, role := range roles {
		r := roleToResp(role)
		rs = append(rs, r)
	}
	return rs
}

func roleToResp(role *datamodels.Role) *dtos.RespRole {
	r := &dtos.RespRole{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		DisplayName: role.DisplayName,
		Perms:       nil,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return r
}
