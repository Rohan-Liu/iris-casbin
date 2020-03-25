package databases

import (
	"../configs"
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"path/filepath"
)

var (
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
)

func init() {
	var err error
	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		configs.GetConfigString("datasource.username"), configs.GetConfigString("datasource.password"),
		configs.GetConfigString("datasource.host"), configs.GetConfigString("datasource.database"))
	datasourceType := configs.GetConfigString("datasource.type")
	Db, err = gorm.Open(datasourceType, conn)
	if err != nil {
		color.Red(fmt.Sprintf("gorm open 错误: %v", err))
	}
	c, err := gormadapter.NewAdapter(datasourceType, conn, true) // Your driver and data source.
	if err != nil {
		color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	Enforcer, err = casbin.NewEnforcer(filepath.Join(configs.Root, "configs", "rbac_model.conf"), c)

	if err != nil {
		color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}
}
