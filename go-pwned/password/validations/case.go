package validations

import (
	"fmt"
	"unicode"
)

type Case struct {
	Enabled   bool `yaml:"enabled"`
	OnlyUpper bool `yaml:"onlyUpper"`
	OnlyLower bool `yaml:"onlyLower"`
	MinUpper  int  `yaml:"minUpper"`
	MinLower  int  `yaml:"minLower"`
}

func (cr *Case) Validate(password string) (bool, error) {

	if cr.Enabled {

		if cr.OnlyUpper {
			if !areAllUpper(password) {
				return false, fmt.Errorf("password should only contain uppercase characters")
			}
		}

		if cr.OnlyLower {
			if !areAllLower(password) {
				return false, fmt.Errorf("password should only contain lowercase characters")
			}
		}

		totalLower := countLower(password)

		if totalLower < cr.MinLower {
			return false, fmt.Errorf("password does not contain at least %d lower characters", cr.MinLower)
		}

		totalUpper := countUpper(password)

		if totalUpper < cr.MinUpper {
			return false, fmt.Errorf("password does not contain at least %d upper characters", cr.MinUpper)
		}

	}

	return true, nil
}

func areAllUpper(str string) bool {

	for _, char := range str {
		if unicode.IsLetter(char) {
			if !unicode.IsUpper(char) {
				return false
			}
		}
	}

	return true
}

func areAllLower(str string) bool {

	for _, char := range str {
		if unicode.IsLetter(char) {
			if !unicode.IsLower(char) {
				return false
			}
		}
	}

	return true
}

func countLower(str string) (total int) {

	for _, char := range str {
		if unicode.IsLower(char) {
			total++
		}
	}

	return
}

func countUpper(str string) (total int) {

	for _, char := range str {
		if unicode.IsUpper(char) {
			total++
		}
	}

	return
}
