package main

import "fmt"

type Tier struct {
	name string
}

func (r Tier) String() string { return r.name }

var (
	UnknownTier = Tier{""}
	GuestTier   = Tier{"guest"}
	BasicTier   = Tier{"basic"}
	PremiumTier = Tier{"premium"}
)

func TierFromString(s string) (Tier, error) {
	switch s {
	case GuestTier.name:
		return GuestTier, nil
	case BasicTier.name:
		return BasicTier, nil
	case PremiumTier.name:
		return PremiumTier, nil
	default:
		return UnknownTier, fmt.Errorf("unknown tier: %s", s)
	}
}
