package handler

import (
	"fmt"
)

type IPv6URLCapitalization string

const (
	IPv6URLCapitalizationLower IPv6URLCapitalization = "lower"
	IPv6URLCapitalizationUpper IPv6URLCapitalization = "upper"
)

//goland:noinspection GoMixedReceiverTypes
func (c IPv6URLCapitalization) String() string {
	return string(c)
}

//goland:noinspection GoMixedReceiverTypes
func (c *IPv6URLCapitalization) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*c = ""
		return nil
	}

	s := string(text)
	switch s {
	case string(IPv6URLCapitalizationLower):
		*c = IPv6URLCapitalizationLower
	case string(IPv6URLCapitalizationUpper):
		*c = IPv6URLCapitalizationUpper
	default:
		return fmt.Errorf("invalid IPv6 URL capitalization: %s", s)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (c IPv6URLCapitalization) MarshalText() (text []byte, err error) {
	return []byte(c), nil
}
