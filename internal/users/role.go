package users

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
