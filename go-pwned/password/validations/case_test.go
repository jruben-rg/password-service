package validations

import (
	"testing"
)

func TestCaseDisabled(t *testing.T) {

	tests := []struct {
		scenario    string
		caseVal     Case
		password    string
		expectedErr error
	}{
		{
			scenario: "When the case rule is disabled, password case is not validated",
			caseVal: Case{
				Enabled:   false,
				OnlyUpper: true,
				OnlyLower: true,
				MinLower:  3,
				MinUpper:  3,
			},
			password:    "Passw0rd",
			expectedErr: nil,
		},
	}

	for _, test := range tests {

		cv := &test.caseVal
		ok, actualErr := cv.Validate(test.password)
		if actualErr != test.expectedErr {
			t.Errorf("Scenario: '%q'. Want: '%q' Got: '%q'\n", test.scenario, test.expectedErr, actualErr)
		}
		if ok != true {
			t.Error("Ok was expected to be true.")
		}
	}

}

func TestValidateCaseShouldFail(t *testing.T) {

	tests := []struct {
		scenario  string
		caseVal   Case
		passwords []string
	}{
		{
			scenario: "Should validate only uppercase",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: true,
				OnlyLower: false,
				MinLower:  0,
				MinUpper:  0,
			},
			passwords: []string{"Passw0rd", "MyPASS_123", "PASSWoRD"},
		},
		{
			scenario: "Should validate max uppercase length",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: false,
				MinLower:  0,
				MinUpper:  3,
			},
			passwords: []string{"Passw0rd", "02aBCd94", "AbcdeF"},
		},
		{
			scenario: "Should validate only lowercase",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: true,
				MinLower:  0,
				MinUpper:  0,
			},
			passwords: []string{"Passw0rd", "passW0rd", "myPassw0rd"},
		},
		{
			scenario: "Should validate min lowercase length",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: false,
				MinLower:  8,
				MinUpper:  0,
			},
			passwords: []string{"passw0rd", "123ab*cd_ef9", "ab8723_78ef"},
		},
	}

	for _, test := range tests {

		cv := &test.caseVal
		for _, password := range test.passwords {
			ok, actualErr := cv.Validate(password)
			if actualErr == nil {
				t.Errorf("Wanted an error for scenario: %q\n", test.scenario)
			}

			if ok != false {
				t.Error("Ok was expected to be false.")
			}
		}
	}

}

func TestValidateCaseShouldPass(t *testing.T) {

	tests := []struct {
		scenario  string
		caseVal   Case
		passwords []string
	}{
		{
			scenario: "Should validate a valid uppercase password",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: true,
				OnlyLower: false,
				MinLower:  0,
				MinUpper:  0,
			},
			passwords: []string{"PASSW0RD", "MYP4SSW0RD", "PASSWORD"},
		},
		{
			scenario: "Should validate the length of uppercase password",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: true,
				OnlyLower: false,
				MinLower:  0,
				MinUpper:  7,
			},
			passwords: []string{"PASSW0RD", "MY_23_PASSWD", "A_PASSWD"},
		},
		{
			scenario: "Should validate a lowercase password",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: true,
				MinLower:  7,
				MinUpper:  0,
			},
			passwords: []string{"passw0rd", "myp4&sswrd", "password"},
		},
		{
			scenario: "Should validate a the length of lowercase password",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: true,
				MinLower:  4,
				MinUpper:  0,
			},
			passwords: []string{"passw0rd", "pa123^ss", "87pass23"},
		},
		{
			scenario: "Should validate a mixed case password",
			caseVal: Case{
				Enabled:   true,
				OnlyUpper: false,
				OnlyLower: false,
				MinLower:  4,
				MinUpper:  4,
			},
			passwords: []string{"PASSword", "myPASSw0rd", "$paSSWOrd"},
		},
	}

	for _, test := range tests {

		cv := &test.caseVal
		for _, password := range test.passwords {
			ok, actualErr := cv.Validate(password)
			if actualErr != nil {
				t.Errorf("Scenario '%q' Got error: '%q' for password: '%s'\n", test.scenario, actualErr, password)
			}
			if ok != true {
				t.Error("Ok was expected to be true.")
			}
		}

	}

}
