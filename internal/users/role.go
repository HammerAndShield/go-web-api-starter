package users

import (
	"database/sql"
)

const (
	RoleRegularUser = "regular"
)

type Role struct {
	Name        string
	Permissions Permissions
}

var RegularRole = Role{
	Name:        RoleRegularUser,
	Permissions: Permissions{},
}

type RolePsqlRepo struct {
	DB *sql.DB
}
