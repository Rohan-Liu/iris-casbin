package controllers

import (
	"../../datamodels"
	"../../services"
	"../dtos"
	"github.com/go-playground/validator"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

// 分页列表
type Lists struct {
	Items      interface{} `json:"items"`
	TotalCount uint64      `json:"totalCount"`
}

type PathName struct {
	Name   string
	Path   string
	Method string
}

func ToResult(status bool, objects interface{}, msg string) (r *Response) {
	r = &Response{Status: status, Data: objects, Msg: msg}
	return
}

func ToSuccess(objects interface{}, msg string) (r *Response) {
	return ToResult(true, objects, msg)
}

func ToFailure(msg string) (r *Response) {
	return ToResult(false, nil, msg)
}

type AppController struct {
	UserService services.IUserService
	RoleService services.IRoleService
	PermService services.IPermissionService
}

// 重置系统数据
// 管理端 管理员数据
// 账号，角色，权限
func (c *AppController) GetReset(ctx iris.Context) {

	c.RoleService.Clear()
	c.UserService.Clear()
	c.PermService.Clear()

	role := &datamodels.Role{
		Name:        "管理员",
		Permissions: getPermissions(ctx.Application().GetRoutesReadOnly()),
	}

	err := c.RoleService.Create(role)
	if err != nil {

		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure("创建角色失败：" + err.Error()))
	}

	user := datamodels.User{
		Name:     "rohan",
		Username: "rohan",
		Password: "111",
		//Password: utils.HashPassword("111"),
		Roles: []*datamodels.Role{
			role,
		},
	}
	err = c.UserService.Create(&user)

	if err != nil {

		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ToFailure("创建角色失败：" + err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ToSuccess(user, "重置数据成功"))
}

func getPathNames(i interface{}) []*PathName {
	var pns []*PathName
	if routeReadOnly, ok := i.([]context.RouteReadOnly); ok {
		for _, s := range routeReadOnly {
			pn := &PathName{
				Name:   s.Name(),
				Path:   s.Path(),
				Method: s.Method(),
			}
			pns = append(pns, pn)
		}
	} else if route, ok := i.([]*router.Route); ok {
		for _, s := range route {
			pn := &PathName{
				Name:   s.Name,
				Path:   s.Path,
				Method: s.Method,
			}
			pns = append(pns, pn)
		}
	}
	return pns
}

func getPermissions(i interface{}) (perms []*datamodels.Permission) {
	for _, s := range getPathNames(i) {
		if strings.HasPrefix(s.Path, "/v1/admin/") && isPermRoute(s.Method) {
			perms = append(perms, datamodels.NewPermission(0, s.Path, "", "", s.Method))
		}
	}
	return
}

func isPermRoute(name string) bool {
	exceptRouteName := []string{"GET", "POST", "DELETE", "PUT"}
	for _, er := range exceptRouteName {
		if strings.Contains(name, er) {
			return true
		}
	}
	return false
}

func (c AppController) PostLogin(ctx iris.Context) {
	req := new(dtos.ReqLogin)

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

	ctx.Application().Logger().Infof("%s 登录系统", req.Username)
	ctx.StatusCode(iris.StatusOK)
	user, err := c.UserService.GetByName(req.Username)
	if err != nil {

	}
	response, status, msg := checkLogin(&user, req.Password)

	_, _ = ctx.JSON(ToResult(status, response, msg))
	return

}

func (c AppController) GetLogout(ctx iris.Context) {

	value := ctx.Values().Get("jwt").(*jwt.Token)
	if value == nil || value.Raw == "" {
		_, _ = ctx.JSON(ToFailure("获取Token出错"))
		return
	}

	token := dtos.OauthToken{}
	token.GetOauthTokenByToken(value.Raw)
	if token.Token == "" {
		_, _ = ctx.JSON(ToFailure("获取Token出错"))
		return
	}

	token.RemoveOauthTokenByToken()

	ctx.Application().Logger().Infof("%d 退出系统", token.UserId)
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ToSuccess(nil, "退出"))
}

func checkLogin(user *datamodels.User, password string) (*dtos.Token, bool, string) {
	if user.ID == 0 {
		return nil, false, "用户不存在"
	} else {
		if ok := user.CheckPassword(password); ok {
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(100)).Unix(), // 失效时间
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))
			var roleIds []string
			var roleName string
			for _, role := range user.Roles {
				roleIds = append(roleIds, strconv.FormatUint(uint64(role.ID), 10))
				if len(roleIds) > 0 {
					roleName += ","
				}
				roleName += role.Name
			}

			oauthToken := new(dtos.OauthToken)
			oauthToken.Name = user.Name
			oauthToken.Token = tokenString
			oauthToken.UserId = user.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(100)).Unix()
			oauthToken.RoleIds = roleIds
			oauthToken.RoleName = roleName
			response := oauthToken.OauthTokenCreate()

			return response, true, "登陆成功"
		} else {
			return nil, false, "用户名或密码错误"
		}
	}
}
