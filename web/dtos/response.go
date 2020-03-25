package dtos

type RespUser struct {
	Id        uint
	Name      string
	Username  string
	RoleName  string
	RoleIds   []uint
	CreatedAt string
}

type RespRole struct {
	Id          uint
	Name        string
	DisplayName string
	Description string
	Perms       []*RespPermission
	CreatedAt   string
}

type RespPermission struct {
	Id          uint
	Name        string
	DisplayName string
	Description string
	CreatedAt   string
}
