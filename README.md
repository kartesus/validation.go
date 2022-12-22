# Validation

Prototype for a compositional library for validation in Go.

## Usage

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
