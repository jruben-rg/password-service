package validations

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Number struct {
	Enabled      bool `yaml:"enabled"`
	AllowNumbers bool `yaml:"allowNumbers"`
	Min          int  `yaml:"min"`
	OnlyNumbers  bool `yaml:"onlyNumbers"`
}

func (nr *Number) Validate(password string) (bool, error) {

	if nr.Enabled {
		totalNumbers := countNumbers(password)
		if nr.AllowNumbers {
			if totalNumbers < nr.Min {
				return false, fmt.Errorf("password should contain at least %d numbers", nr.Min)
			}
		} else {
			// it contains numbers and it shouldn't
			if totalNumbers > 0 {
				return false, fmt.Errorf("password should not contain numbers")
			}
		}

		if nr.OnlyNumbers {
			if totalNumbers != utf8.RuneCountInString(password) {
				return false, fmt.Errorf("password should only contain numbers")
			}
		}
	}

	return true, nil
}

func countNumbers(str string) (total int) {

	for _, char := range str {
		if unicode.IsNumber(char) {
			total++
		}
	}

	return
}
