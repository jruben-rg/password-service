package validations

import (
	"strings"
	"testing"
)

func TestSymbolDisabled(t *testing.T) {

	symbolRules := &Symbol{
		Enabled:        false,
		UseSymbol:      false,
		AllowedSymbols: "",
	}

	ok, result := symbolRules.Validate("passw0rd")

	if result != nil {
		t.Errorf("Symbol validator returned error: %q\n", result)
	}

	if ok != true {
		t.Error("Was expecting ok to be true")
	}
}

func TestValidateSymbolShouldFail(t *testing.T) {

	tests := []struct {
		scenario    string
		symbolRules Symbol
		passwords   []string
	}{
		{
			scenario: "Password does not contain any of the allowed symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      true,
				Min:            1,
				AllowedSymbols: "%@#()",
			},
			passwords: []string{"pass123!*&", "SOME_*+", "£&_?PASS"},
		},
		{
			scenario: "Password does not contain some of the allowed symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      true,
				Min:            1,
				AllowedSymbols: "!*&",
			},
			passwords: []string{"Pass123!*#", "PASS*_", "P4__W*RD"},
		},
		{
			scenario: "Password does not contain min valid symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      true,
				Min:            4,
				AllowedSymbols: "!*#",
			},
			passwords: []string{"Pass123!*#", "MyP4**!d", "N*Tv$l!#"},
		},
		{
			scenario: "Password should not contain any symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      false,
				Min:            0,
				AllowedSymbols: "",
			},
			passwords: []string{"pass123!*&", "N*tV4l!d", "PAZZ!"},
		},
	}

	for _, test := range tests {
		sr := test.symbolRules
		for _, password := range test.passwords {
			ok, expected := sr.Validate(password)

			if expected == nil {
				t.Errorf("Expected error for scenario '%s'\n", test.scenario)
			}
			if ok != false {
				t.Errorf("Was expecting ok to be false for scenario %s\n", test.scenario)
			}
		}

	}
}

func TestValidateSymbolShouldPass(t *testing.T) {

	tests := []struct {
		scenario    string
		symbolRules Symbol
		passwords   []string
	}{
		{
			scenario: "Password contains some allowed symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      true,
				Min:            0,
				AllowedSymbols: "%@#()*&^!",
			},
			passwords: []string{"passw0rd123!*&", "pass((#))", "!p^ssw*rd!"},
		},
		{
			scenario: "Password contains the min required symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      true,
				Min:            3,
				AllowedSymbols: "%@#()*&^!",
			},
			passwords: []string{"PaSSw0rd123%@#()*&^!", "my@(pass)", "#lets@Secure!"},
		},
		{
			scenario: "Password should not contain symbols",
			symbolRules: Symbol{
				Enabled:        true,
				UseSymbol:      false,
				Min:            0,
				AllowedSymbols: "",
			},
			passwords: []string{"PaSSw0rd123", "MyUnsecP4SS", "N0Ts3cur3"},
		},
	}

	for _, test := range tests {
		sr := test.symbolRules
		for _, password := range test.passwords {
			ok, error := sr.Validate(password)

			if error != nil {
				t.Errorf("Got unexpected error %s for scenario %s\n", error.Error(), test.scenario)
			}
			if ok == false {
				t.Errorf("Was expecting ok to be true for scenario %s\n", test.scenario)
			}
		}
	}
}

func TestShouldEscapePercentSymbol(t *testing.T) {

	tests := []struct {
		scenario       string
		symbols        string
		expectedResult string
	}{
		{
			scenario:       "Should not append any percentage symbol if there isn't any",
			symbols:        "~!@#$^&*()_-+={}[]|:;<>,.?/",
			expectedResult: "~!@#$^&*()_-+={}[]|:;<>,.?/",
		},
		{
			scenario:       "Should append a percentage symbol if there is any",
			symbols:        "%~!@#$^&*()_-+%={}[]|:;<>,.?/%",
			expectedResult: "%%~!@#$^&*()_-+%%={}[]|:;<>,.?/%%",
		},
	}

	for _, test := range tests {

		result := escapePercentSymbol(test.symbols)
		if strings.Compare(result, test.expectedResult) != 0 {
			t.Errorf("Scenario '%s', Expected: '%s', Got: '%s'\n", test.scenario, test.expectedResult, result)
		}

	}

}
