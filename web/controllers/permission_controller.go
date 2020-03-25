package controllers

import (
	"../../datamodels"
	"../../services"
	"../dtos"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	//gf "github.com/snowlyg/gotransformer"
)

type PermissionController struct {
	Service services.IPermissionService
}

func (c *PermissionController) GetBy(ctx iris.Context, id uint) {
	permission, err := c.Service.Get(id)

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ToSuccess(permToDto(permission), "操作成功"))
}

func (c *PermissionController) Post(ctx iris.Context) {
	req := new(dtos.ReqPermission)
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

	perm := &datamodels.Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Action:      req.Action,
	}
	err = c.Service.Create(perm)

	ctx.StatusCode(iris.StatusOK)
	if perm.ID == 0 {
		_, _ = ctx.JSON(ToFailure("操作失败"))
	} else {
		_, _ = ctx.JSON(ToSuccess(perm, "操作成功"))
	}

}

// 更新Permission
// PUT /admin/permission/:id
func (c *PermissionController) PutBy(ctx iris.Context, id uint) {
	req := new(dtos.ReqPermission)

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

	perm := &datamodels.Permission{
		Model: gorm.Model{
			ID: id,
		},
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Action:      req.Action,
	}
	err = c.Service.Update(perm)

	ctx.StatusCode(iris.StatusOK)
	if perm.ID == 0 {
		_, _ = ctx.JSON(ToFailure("操作失败"))
	} else {
		_, _ = ctx.JSON(ToSuccess(perm, "操作成功"))
	}

}

func (c *PermissionController) DeleteBy(ctx iris.Context, id uint) {
	err := c.Service.Delete(id)
	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ToSuccess(nil, "删除成功"))
}

func (c *PermissionController) GetPage(ctx iris.Context) {
	pageIndex := ctx.Values().GetUintDefault("pageIndex", 1)
	pageSize := ctx.Values().GetUintDefault("pageSize", 20)
	name := ctx.Values().GetStringDefault("name", "")
	orderBy := ctx.Values().GetStringDefault("orderBy", "")

	permissions, totalCount, err := c.Service.GetPage(name, orderBy, uint64(pageIndex), uint64(pageSize))

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	list := Lists{
		Items:      permsToDto(permissions),
		TotalCount: totalCount,
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ToSuccess(list, "操作成功"))
}

func permsToDto(perms []*datamodels.Permission) (list []*dtos.RespPermission) {
	for _, perm := range perms {
		list = append(list, permToDto(perm))
	}
	return
}

func permToDto(perm *datamodels.Permission) (item *dtos.RespPermission) {
	item = &dtos.RespPermission{
		Id:          perm.ID,
		Name:        perm.Name,
		DisplayName: perm.DisplayName,
		Description: perm.Description,
		CreatedAt:   perm.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return
}
