package validations

import (
	"fmt"
	"strings"
	"unicode"
)

type Symbol struct {
	Enabled        bool   `yaml:"enabled"`
	UseSymbol      bool   `yaml:"allowSymbols"`
	Min            int    `yaml:"min"`
	AllowedSymbols string `yaml:"allowedSymbols"`
}

func (s *Symbol) Validate(password string) (bool, error) {

	if s.Enabled {
		totalSymbols := countSymbols(password)

		if !s.UseSymbol && totalSymbols > 0 {
			return false, fmt.Errorf("password should not contain any symbols")
		}

		if s.UseSymbol {
			if totalSymbols > 0 {
				areValidSymbols, totalValid, invalidSymbols := validateSymbols(escapePercentSymbol(s.AllowedSymbols), password)
				if !areValidSymbols {
					return false, fmt.Errorf("password contains invalid symbols '%s'", invalidSymbols)
				}

				if totalValid < s.Min {
					return false, fmt.Errorf("password does not contain at least %d valid symbols (%s)", s.Min, escapePercentSymbol(s.AllowedSymbols))
				}
			} else {
				return false, fmt.Errorf("password does not contain any of the allowed symbols '%s'", escapePercentSymbol(s.AllowedSymbols))
			}
		}
	}

	return true, nil
}

func countSymbols(str string) (total int) {

	for _, char := range str {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			total++
		}
	}

	return
}

func validateSymbols(symbols, password string) (bool, int, string) {

	invalidChars := ""
	validCount := 0
	for _, char := range password {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			if !strings.ContainsAny(symbols, string(char)) {
				invalidChars += string(char)
			} else {
				validCount++
			}
		}
	}

	if invalidChars != "" {
		return false, validCount, invalidChars
	}

	return true, validCount, ""
}

// Golang needs the percentage symbol to be escaped, this is by appending it
func escapePercentSymbol(symbols string) string {

	if len(symbols) > 0 && strings.Contains(symbols, "%") {

		newSymbols := ""
		for _, symbol := range symbols {
			if symbol == '%' {
				newSymbols += "%"
			}
			newSymbols += string(symbol)
		}

		return newSymbols
	}

	return symbols
}
