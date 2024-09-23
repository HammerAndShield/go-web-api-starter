package users

import (
	"slices"
)

const (
	PermUsersManage = "users:manage"
)

type Permissions []string

// Includes checks if the given permission codes exist within the Permissions slice.
// For each code in the provided variadic argument, it checks if the code is present
// in the permission list. The function returns true only when all the given codes
// are present. If any single code is not found, it stops further checking and returns false.
func (p Permissions) Includes(codes ...string) bool {
	for _, perm := range codes {
		if !slices.Contains(p, perm) {
			return false
		}
	}
	return true
}
