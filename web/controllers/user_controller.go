package controllers

import (
	"../../datamodels"
	"../../services"
	"../dtos"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
)

type UserController struct {
	Service     services.IUserService
	RoleService services.IRoleService
}

func (c *UserController) Get(ctx iris.Context) {
	result, _ := c.Service.GetAll()
	_, _ = ctx.JSON(ToSuccess(usersToResp(result), "操作成功"))
}

func (c *UserController) GetProfile(ctx iris.Context) {
	userId := ctx.Values().Get("auth_user_id").(uint)
	user, err := c.Service.Get(userId)
	if err != nil {

		_, _ = ctx.JSON(ToFailure(err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)

	_, _ = ctx.JSON(ToSuccess(userToResp(&user), "操作成功"))
}

func (c *UserController) GetBy(ctx iris.Context, id uint) {
	user, err := c.Service.Get(id)
	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)

	_, _ = ctx.JSON(ToSuccess(userToResp(&user), "操作成功"))
}

func (c *UserController) Post(ctx iris.Context) {

	req := new(dtos.ReqUser)
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

	roles, _ := c.RoleService.GetByIds(req.RoleIds)

	user := &datamodels.User{
		Name:     req.Name,
		Username: req.Username,
		Roles:    roles,
	}
	_ = c.Service.Create(user)

	ctx.StatusCode(iris.StatusOK)
	if user.ID == 0 {
		_, _ = ctx.JSON(ToFailure("操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ToSuccess(nil, "操作成功"))
		return
	}

}

func (c *UserController) PutBy(ctx iris.Context, id uint) {
	req := new(dtos.ReqUser)

	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure(err.Error()))
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

	user, err := c.Service.Get(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	if user.Username == "rohan" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure("不能编辑管理员"))
		return
	}

	user.Name = req.Name
	user.Username = req.Username
	user.Roles = nil

	roles, _ := c.RoleService.GetByIds(req.RoleIds)

	err = c.Service.Update(&user, roles)

	ctx.StatusCode(iris.StatusOK)
	if user.ID == 0 {
		_, _ = ctx.JSON(ToFailure("操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ToSuccess(nil, "操作成功"))
		return
	}

}

func (c *UserController) DeleteBy(ctx iris.Context, id uint) {

	err := c.Service.Delete(id)
	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ToSuccess(nil, "删除成功"))
}

func (c *UserController) GetPage(ctx iris.Context) {
	pageIndex := ctx.Values().GetUintDefault("pageIndex", 1)
	pageSize := ctx.Values().GetUintDefault("pageSize", 20)
	name := ctx.Values().GetStringDefault("name", "")
	orderBy := ctx.Values().GetStringDefault("orderBy", "")

	result, totalCount, err := c.Service.GetPage(name, orderBy, uint64(pageIndex), uint64(pageSize))

	if err != nil {
		_, _ = ctx.JSON(ToFailure(err.Error()))
		return
	}

	list := Lists{
		Items:      usersToResp(result),
		TotalCount: totalCount,
	}
	_, _ = ctx.JSON(ToSuccess(list, "操作成功"))

}

func usersToResp(users []*datamodels.User) []*dtos.RespUser {
	var us []*dtos.RespUser
	for _, user := range users {
		u := userToResp(user)
		us = append(us, u)
	}
	return us
}

func userToResp(user *datamodels.User) *dtos.RespUser {
	u := &dtos.RespUser{}
	u.Name = user.Name
	u.Username = user.Username
	var ris []uint
	var roleName string
	for num, role := range user.Roles {
		ris = append(ris, role.ID)
		if num == len(ris)-1 {
			roleName += role.Name
		} else {
			roleName += role.Name + ","
		}
	}
	u.Id = user.ID
	u.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	u.RoleIds = ris
	u.RoleName = roleName
	return u
}
