package password

import (
	"fmt"
	"testing"
)

type testValidator struct {
	isValid bool
	err     error
}

func (f testValidator) Validate(str string) (bool, error) {
	return f.isValid, f.err
}

func TestToListShouldReturnAListOfValidators(t *testing.T) {

	validations := &Validations{}
	validators := validations.ToList()
	if len(validators) != 4 {
		t.Errorf("Expected validators %d, Got: %d\n", 4, len(validators))
	}
}

func TestValidateShouldReturnErrorIfAnyValidationReturnsError(t *testing.T) {

	tests := []struct {
		scenario   string
		validators []Validator
	}{
		{
			scenario: "One instance returns an error",
			validators: []Validator{
				&testValidator{false, fmt.Errorf("test error")},
				&testValidator{true, nil},
			},
		},
		{
			scenario: "Two instances return error",
			validators: []Validator{
				&testValidator{false, fmt.Errorf("test error 1")},
				&testValidator{false, fmt.Errorf("test error 2")},
			},
		},
	}

	for _, test := range tests {

		password := password{test.validators}
		isValid, err := password.Validate("APassw0rd!")

		if err == nil {
			t.Errorf("Scenario: %s. Was expecting an error", test.scenario)
		}

		if isValid != false {
			t.Errorf("Scenario: %s. Was expecting 'isValid' to be false.", test.scenario)
		}
	}

}

func TestValidateShouldSucceedIfAllValidationsSucceed(t *testing.T) {

	tests := []struct {
		scenario   string
		validators []Validator
	}{
		{
			scenario: "One instance returns an error",
			validators: []Validator{
				&testValidator{true, nil},
				&testValidator{true, nil},
			},
		},
		{
			scenario: "Two instances return error",
			validators: []Validator{
				&testValidator{true, nil},
				&testValidator{true, nil},
			},
		},
	}

	for _, test := range tests {

		password := password{test.validators}
		isValid, err := password.Validate("APassw0rd!")

		if err != nil {
			t.Errorf("Scenario: %s. Wasn't expecting an error. Got : '%s'", test.scenario, err)
		}

		if isValid != true {
			t.Errorf("Scenario: %s. Was expecting 'isValid' to be true.", test.scenario)
		}
	}

}
