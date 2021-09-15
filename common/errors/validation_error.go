package errors

type ValidationError struct {
	Fields []ValidationErrorField
}

type ValidationErrorField struct {
	Error string
	Field string
}

func NewValidationError(fields []ValidationErrorField) *ValidationError {
	return &ValidationError{
		Fields: fields,
	}
}

func (v ValidationError) Error() string {
	return "Validation Error"
}
