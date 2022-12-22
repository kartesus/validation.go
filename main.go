package main

import (
	"fmt"
	"regexp"
)

type Validator func(value any) (any, []error)

func String() Validator {
	return func(value any) (any, []error) {
		if str, ok := value.(string); ok {
			return str, nil
		}

		return nil, []error{fmt.Errorf("%v is not a string", value)}
	}
}

func Min(min int) Validator {
	return func(value any) (any, []error) {
		if str, ok := value.(string); ok {
			if len(str) < min {
				return nil, []error{fmt.Errorf("%s must have at lease %d characters", value, min)}
			}

			return str, nil
		}

		return nil, []error{fmt.Errorf("%v is not a string", value)}
	}
}

func Max(max int) Validator {
	return func(value any) (any, []error) {
		if str, ok := value.(string); ok {
			if len(str) > max {
				return nil, []error{fmt.Errorf("%s must have at most %d characters", value, max)}
			}

			return str, nil
		}

		return nil, []error{fmt.Errorf("%v is not a string", value)}
	}
}

func All(validators ...Validator) Validator {
	return func(value any) (any, []error) {
		var errs []error
		for _, validator := range validators {
			value, errs = validator(value)
			if errs != nil {
				errs = append(errs, errs...)
			}
		}

		if errs != nil {
			return nil, errs
		}

		return value, nil
	}
}

func Match(pattern string) Validator {
	return func(value any) (any, []error) {
		if str, ok := value.(string); ok {
			if !regexp.MustCompile(pattern).MatchString(str) {
				return nil, []error{fmt.Errorf("%s does not match %s", value, pattern)}
			}

			return str, nil
		}

		return nil, []error{fmt.Errorf("%v is not a string", value)}
	}
}

func NameValidator() Validator {
	return All(String(), Min(3), Max(50))
}

func EmailValidator() Validator {
	return All(String(), Match(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`))
}

func FiscalNumberValidator() Validator {
	return All(String(), Match(`^[0-9]{10}$`))
}

type Input struct {
	data       map[string]any
	validators map[string]Validator
}

func NewInput(data map[string]any) *Input {
	return &Input{
		data:       data,
		validators: make(map[string]Validator),
	}
}

func (i *Input) ValidateWith(name string, validator Validator) {
	i.validators[name] = validator
}

func (i *Input) Validate(result ValidationResult) {
	input := make(map[string]any)
	errors := make(map[string][]error)

	for name, validator := range i.validators {
		if value, ok := i.data[name]; ok {
			value, errs := validator(value)
			if errs != nil {
				errors[name] = errs
			} else {
				input[name] = value
			}
		} else {
			errors[name] = []error{fmt.Errorf("missing %s", name)}
		}
	}

	if len(errors) > 0 {
		result.ValidationFailed(errors)
	} else {
		result.ValidationSucceeded(input)
	}
}

type ValidationResult interface {
	ValidationFailed(map[string][]error)
	ValidationSucceeded(map[string]any)
}

type PrintValidationResult struct{}

func (r *PrintValidationResult) ValidationFailed(errs map[string][]error) {
	fmt.Printf("Errors: %v\n", errs)
}

func (r *PrintValidationResult) ValidationSucceeded(request map[string]any) {
	fmt.Printf("Request: %v\n", request)
}

func NewValidationResultPrinter() *PrintValidationResult {
	return &PrintValidationResult{}
}

func main() {
	requestData := map[string]any{"customerName": "John Doe", "customerEmail": "john@gmail.com", "fiscalNumber": "1234567890"}

	input := NewInput(requestData)
	input.ValidateWith("customerName", NameValidator())
	input.ValidateWith("customerEmail", EmailValidator())
	input.ValidateWith("fiscalNumber", FiscalNumberValidator())

	input.Validate(NewValidationResultPrinter())
}
