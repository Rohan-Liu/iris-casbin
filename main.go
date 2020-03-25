package main

import (
	"./configs"
	"./databases"
	"./datamodels"
	"./repositories"
	"./services"
	"./web/controllers"
	"./web/middlewares"
	"fmt"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"os"
	"time"
)

func main() {
	f := NewLogFile()
	defer f.Close()

	yaag.Init(&yaag.Config{
		On:       true, //是否开启自动生成API文档功能
		DocTitle: "Iris-Casbin",
		DocPath:  "apidoc.html", //生成API文档名称存放路径
	})

	api := NewApp()
	//注册中间件

	api.Logger().SetOutput(f) //记录日志

	if err := api.Run(iris.Addr(":8085"), iris.WithConfiguration(configs.YamlConf)); err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}

func NewLogFile() *os.File {
	path := "./logs/"
	filename := path + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		color.Red(fmt.Sprintf("日志记录出错: %v", err))
	}

	return f
}

func NewApp() *iris.Application {
	api := iris.New()
	api.Logger().SetLevel("debug")

	db := databases.Db
	db.AutoMigrate(
		&datamodels.User{},
		&datamodels.Role{},
		&datamodels.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	RegisterApp(api) //注册 app 路由

	return api
}

var (
	userService = services.UserService(repositories.UserRepository(databases.Db, databases.Enforcer))
	roleService = services.RoleService(repositories.RoleRepository(databases.Db, databases.Enforcer))
	permService = services.PermissionService(repositories.PermissionRepository(databases.Db, databases.Enforcer))
)

func RegisterApp(api *iris.Application) {

	api.HandleDir("/static", "resources/app/static")
	api.Get("/", func(ctx iris.Context) { // 首页模块
		_ = ctx.View("app/index.html")
	})
	api.Use(irisyaag.New())
	app := api.Party("/", middleware.CorsAuth()).AllowMethods(iris.MethodOptions)
	{

		mvc.Configure(app.Party("/v1"), func(application *mvc.Application) {
			application.Register(userService)
			application.Register(roleService)
			application.Register(permService)
			application.Handle(new(controllers.AppController))

		})

		casbinMiddleware := middleware.New(databases.Enforcer)

		v1 := app.Party("/v1")
		{

			admin := v1.Party("/admin")
			{

				admin.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				mvc.Configure(admin.Party("/user"), func(application *mvc.Application) {

					application.Register(userService)
					application.Register(roleService)
					application.Handle(new(controllers.UserController))
				})

				mvc.Configure(admin.Party("/role"), func(application *mvc.Application) {

					application.Register(userService)
					application.Register(roleService)
					application.Register(permService)
					application.Handle(new(controllers.RoleController))
				})
				mvc.Configure(admin.Party("/permission"), func(application *mvc.Application) {

					application.Register(permService)
					application.Handle(new(controllers.PermissionController))
				})
			}
		}
	}
}
