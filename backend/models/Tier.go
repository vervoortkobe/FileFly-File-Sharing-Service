package models

import (
	"database/sql/driver"
	"fmt"
)

type Tier string

const (
	GuestTier   Tier = "guest"
	FreeTier    Tier = "free"
	PremiumTier Tier = "premium"
)

func (t Tier) Value() (driver.Value, error) {
	switch t {
	case GuestTier, FreeTier, PremiumTier:
		return string(t), nil
	default:
		return nil, fmt.Errorf("invalid Tier value: %s", t)
	}
}

func (t *Tier) Scan(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		byteVal, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to scan Tier: value is not string or []byte, got %T", value)
		}
		strVal = string(byteVal)
	}
	scannedTier := Tier(strVal)
	switch scannedTier {
	case GuestTier, FreeTier, PremiumTier:
		*t = scannedTier
		return nil
	default:
		return fmt.Errorf("invalid value for Tier from DB: %s", strVal)
	}
}

func (t Tier) String() string {
	return string(t)
}
