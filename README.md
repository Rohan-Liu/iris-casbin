学习Golang，Iris，Casbin
---
###介绍

为了系统学习Golang的Web开发，选择了Iris框架和Casbin做权限控制的框架。

本示例大量参考了网上各种代码后完善了一些功能。

####本示例使用

+ GO 1.13.8
+ Iris-GO Go的Web框架
+ Casbin 通用权限控制
+ Gorm Go的orm工具
+ MySql 数据库 持久化Casbin权限
+ Redis 存储登录信息
+ Jwt 中间件
+ Cors 中间件 跨域


####目录结构

├── cache  缓存Redis配置

├── configs  配置文件

├── databases  数据源

├── datamodels  数据对象

├── logs 日志存放目录

├── repositories  Repository

├── services  Service

└── web  网站

    ├── controllers  Controllers
    
    ├── dtos 传输对象
    
    └── middlewares 中间件



###参考资料：

[iris-go](https://iris-go.com/)

[casbin](https://casbin.org/)

[gorm](https://gorm.io/)

[IrisAdminApi](https://github.com/snowlyg/IrisAdminApi)