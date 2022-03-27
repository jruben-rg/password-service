package validations

import (
	"fmt"
	"unicode/utf8"
)

type Length struct {
	Enabled bool `yaml:"enabled"`
	Min     int  `yaml:"min"`
	Max     int  `yaml:"max"`
}

func (lr *Length) Validate(password string) (bool, error) {

	if lr.Enabled {

		passLength := utf8.RuneCountInString(password)
		if passLength < lr.Min {
			return false, fmt.Errorf("password should be at least %d characters long", lr.Min)
		}

		if passLength > lr.Max {
			return false, fmt.Errorf("password maximum length allowed is %d characters", lr.Max)
		}
	}

	return true, nil
}
