package handler

import (
	"fmt"
)

type AddrPreference string

const (
	AddrPreferenceFirst AddrPreference = "first"
	AddrPreferenceLast  AddrPreference = "last"
	AddrPreferenceAll   AddrPreference = "all"
)

//goland:noinspection GoMixedReceiverTypes
func (p AddrPreference) String() string {
	return string(p)
}

//goland:noinspection GoMixedReceiverTypes
func (p *AddrPreference) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*p = ""
		return nil
	}

	s := string(text)
	switch s {
	case string(AddrPreferenceFirst):
		*p = AddrPreferenceFirst
	case string(AddrPreferenceLast):
		*p = AddrPreferenceLast
	case string(AddrPreferenceAll):
		*p = AddrPreferenceAll
	default:
		return fmt.Errorf("invalid address preference: %s", s)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (p AddrPreference) MarshalText() (text []byte, err error) {
	return []byte(p), nil
}
