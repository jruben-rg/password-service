package validations

import (
	"testing"
)

func TestValidateLengthDisabled(t *testing.T) {

	lengthRule := &Length{
		Enabled: false,
		Min:     0,
		Max:     0,
	}

	ok, result := lengthRule.Validate("passw0rd")

	if result != nil {
		t.Errorf("Length validator returned error: %q\n", result)
	}
	if ok != true {
		t.Error("Was expecting ok to be true")
	}

}

func TestValidateLengthShouldFail(t *testing.T) {

	tests := []struct {
		scenario   string
		lengthRule Length
		passwords  []string
	}{
		{
			scenario: "Invalid Min Length",
			lengthRule: Length{
				Enabled: true,
				Min:     5,
				Max:     10,
			},
			passwords: []string{"", "pass", "word", "_23*", "José"},
		},
		{
			scenario: "Invalid Max Length",
			lengthRule: Length{
				Enabled: true,
				Min:     5,
				Max:     10,
			},
			passwords: []string{"MyPasswórd1", "MyP&ssw*rd2"},
		},
	}

	for _, test := range tests {
		lr := test.lengthRule
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

func TestValidateLengthShouldPass(t *testing.T) {

	tests := []struct {
		scenario   string
		lengthRule Length
		passwords  []string
	}{
		{
			scenario: "Password with same Min and Max Length",
			lengthRule: Length{
				Enabled: true,
				Min:     5,
				Max:     5,
			},
			passwords: []string{"Ruben", "_1234", "+12PA"},
		},
		{
			scenario: "Password meets min and max length requirements",
			lengthRule: Length{
				Enabled: true,
				Min:     5,
				Max:     7,
			},
			passwords: []string{"Coffee", "(12*45)", "MyPwd", "!*+Pass"},
		},
	}

	for _, test := range tests {
		length := test.lengthRule
		for _, password := range test.passwords {
			ok, error := length.Validate(password)
			if error != nil {
				t.Errorf("Got unexpected error %s for scenario %s\n", error.Error(), test.scenario)
			}
			if ok != true {
				t.Errorf("Was expecting ok to be true for scenario %s\n", test.scenario)
			}
		}

	}

}
