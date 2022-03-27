package validations

import (
	"testing"
)

func TestDisabled(t *testing.T) {

	numberRules := &Number{
		Enabled:      false,
		AllowNumbers: false,
		Min:          0,
		OnlyNumbers:  false,
	}

	ok, result := numberRules.Validate("passw0rd")

	if result != nil {
		t.Errorf("Number validator returned error: %q\n", result)
	}
	if ok != true {
		t.Error("Was expecting ok to be true")
	}

}

func TestNumberValiationShouldFail(t *testing.T) {

	tests := []struct {
		scenario    string
		numberRules Number
		passwords   []string
	}{
		{
			scenario: "Password should only contain numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: true,
				Min:          4,
				OnlyNumbers:  true,
			},
			passwords: []string{"pass01212", "12Pass12", "*123456", "_123456"},
		},
		{
			scenario: "Password does not contain min required numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: true,
				Min:          4,
				OnlyNumbers:  false,
			},
			passwords: []string{"pass012", "1passw0rd9", "myP4ssw0rd"},
		},
		{
			scenario: "Password should not contain numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: false,
				Min:          0,
				OnlyNumbers:  false,
			},
			passwords: []string{"pass012", "456myPass", "passw0rd"},
		},
	}

	for _, test := range tests {
		lr := test.numberRules
		for _, password := range test.passwords {
			ok, expected := lr.Validate(password)
			if expected == nil {
				t.Errorf("Expected error for scenario %s\n", test.scenario)
			}
			if ok != false {
				t.Errorf("Was expecting ok to be false for scenario %s\n", test.scenario)
			}
		}

	}
}

func TestValidationShouldPass(t *testing.T) {

	tests := []struct {
		scenario    string
		numberRules Number
		passwords   []string
	}{
		{
			scenario: "Password should only contain numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: true,
				Min:          0,
				OnlyNumbers:  true,
			},
			passwords: []string{"123459774", "1234", "09498678373764658"},
		},
		{
			scenario: "Password contains min length of numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: true,
				Min:          5,
				OnlyNumbers:  false,
			},
			passwords: []string{"ab123cd34", "12345", "12345a"},
		},
		{
			scenario: "Password should not contain numbers",
			numberRules: Number{
				Enabled:      true,
				AllowNumbers: false,
				Min:          0,
				OnlyNumbers:  false,
			},
			passwords: []string{"abcdE", "_MyPass=", "password+"},
		},
	}

	for _, test := range tests {
		lr := test.numberRules
		for _, password := range test.passwords {
			ok, expected := lr.Validate(password)
			if expected != nil {
				t.Errorf("Got unexpected error for scenario '%s'. Error is: %s\n", test.scenario, expected)
			}
			if ok != true {
				t.Errorf("Was expecting a correct validation for scenario %s\n", test.scenario)
			}
		}
	}
}
