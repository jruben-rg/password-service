package password

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jruben-rg/password-service/go-pwned/config"
	"github.com/jruben-rg/password-service/go-pwned/password/validations"
)

type (
	PasswordConfig struct {
		Validations `yaml:"password"`
	}

	Validations struct {
		Case    validations.Case   `yaml:"case"`
		Length  validations.Length `yaml:"length"`
		Symbols validations.Symbol `yaml:"symbols"`
		Numbers validations.Number `yaml:"numbers"`
	}

	password struct {
		validations []Validator
	}

	Validator interface {
		Validate(str string) (bool, error)
	}

	result struct {
		valid bool
		err   error
	}
)

func NewPasswordConfig(filePath string) *password {
	passwordConfig := &PasswordConfig{}
	err := config.Read(filePath, &passwordConfig)
	if err != nil {
		panic("Could not read configuration for password validations")
	}

	return &password{passwordConfig.Validations.ToList()}
}

func (p *Validations) ToList() []Validator {
	return []Validator{&p.Case, &p.Length, &p.Symbols, &p.Numbers}
}

func (p *password) Validate(password string) (bool, error) {

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(p.validations))

	// Validate password rules using fanIn goroutines
	validatorResult := make(chan result)

	go func() {

		for _, validator := range p.validations {
			go func(ruleValidator Validator) {
				validatorResult <- <-validateWithRule(password, ruleValidator)
				waitGroup.Done()
			}(validator) // Send the current validator as a parameter, otherwise it always process the same
		}

	}()

	go func() {
		//Wait for the validators to be processed
		waitGroup.Wait()

		close(validatorResult)
	}()

	isValid := true
	errorMessages := make([]string, 0)
	//Process validators when ready
	for validatorResult := range validatorResult {
		if !validatorResult.valid {
			isValid = false
			errorMessages = append(errorMessages, validatorResult.err.Error())
		}
	}

	if !isValid {
		return isValid, fmt.Errorf(strings.Join(errorMessages, "\n"))
	}

	return isValid, nil

}

func validateWithRule(password string, validator Validator) <-chan result {
	vr := make(chan result)
	go func() {
		ok, err := validator.Validate(password)
		vr <- result{ok, err}
	}()

	return vr
}
