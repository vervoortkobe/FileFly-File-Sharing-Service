package domain

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	GuestRole Role = "guest"
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

func (r Role) Value() (driver.Value, error) {
	switch r {
	case GuestRole, UserRole, AdminRole:
		return string(r), nil
	default:
		return nil, fmt.Errorf("invalid Role value: %s", r)
	}
}

func (r *Role) Scan(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		byteVal, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to scan Role: value is not string or []byte, got %T", value)
		}
		strVal = string(byteVal)
	}
	scannedRole := Role(strVal)
	switch scannedRole {
	case GuestRole, UserRole, AdminRole:
		*r = scannedRole
		return nil
	default:
		return fmt.Errorf("invalid value for Role from DB: %s", strVal)
	}
}

func (r Role) String() string {
	return string(r)
}
