# Validation

Prototype for a compositional library for validation in Go.

## Usage

A complex validator is composed of primitive validators.
In `NameValidator()` below, `All()` ensures that all validators must succeed for `NameValidator()` to succeed.
It than defines that a name must be a string with at least 3 charachters and at most 50.

```go
func NameValidator() Validator {
	return All(String(), Min(3), Max(50))
}
```

```go
func main() {
	requestData := map[string]any{"customerName": "John Doe", "customerEmail": "john@gmail.com", "fiscalNumber": "1234567890"}

	input := NewInput(requestData)
	input.ValidateWith("customerName", NameValidator())
	input.ValidateWith("customerEmail", EmailValidator())
	input.ValidateWith("fiscalNumber", FiscalNumberValidator())

	input.Validate(NewValidationResultPrinter())
}
```
