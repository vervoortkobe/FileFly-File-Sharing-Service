package main

import "fmt"

type Role struct {
	name string
}

func (r Role) String() string { return r.name }

var (
	UnknownRole = Role{""}
	GuestRole   = Role{"guest"}
	MemberRole  = Role{"member"}
	AdminRole   = Role{"admin"}
)

func RoleFromString(s string) (Role, error) {
	switch s {
	case GuestRole.name:
		return GuestRole, nil
	case MemberRole.name:
		return MemberRole, nil
	case AdminRole.name:
		return AdminRole, nil
	default:
		return UnknownRole, fmt.Errorf("unknown role: %s", s)
	}
}
