//Casbin权限控制中间件

package middleware

import (
	"net/http"
	"strconv"
	"time"

	"../dtos"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/context"

	"github.com/casbin/casbin"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}

func (c *Casbin) ServeHTTP(ctx context.Context) {
	value := ctx.Values().Get("jwt").(*jwt.Token)
	token := dtos.OauthToken{}

	token.GetOauthTokenByToken(value.Raw) //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		//_, _ = ctx.Writef("token 失效，请重新登录") // 输出到前端
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.StopExecution()
		return
	} else if !c.Check(ctx.Request(), strconv.FormatUint(uint64(token.UserId), 10)) && !c.CheckRoles(ctx.Request(), token.RoleIds) {
		ctx.StatusCode(http.StatusForbidden) // Status Forbidden
		ctx.StopExecution()
		return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	ctx.Next()
}

// Casbin is the auth services which contains the casbin enforcer.
type Casbin struct {
	enforcer *casbin.Enforcer
}

// Check checks the username, request's method and path and
// returns true if permission granted otherwise false.
func (c *Casbin) Check(r *http.Request, name string) bool {
	method := r.Method
	path := r.URL.Path
	item := "user:" + name
	ok, _ := c.enforcer.Enforce(item, path, method)
	return ok
}

func (c *Casbin) CheckRoles(r *http.Request, roles []string) bool {
	method := r.Method
	path := r.URL.Path
	for _, role := range roles {
		name := "role:" + role
		if ok, _ := c.enforcer.Enforce(name, path, method); ok {
			return true
		}
	}
	return false
}
