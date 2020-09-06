package mock

type PgError struct {
	FieldCode               string
	ErrorValue              string
	IntegrityViolationValue bool
}

func (e *PgError) Error() string {
	if len(e.ErrorValue) > 0 {
		return e.ErrorValue
	}

	return exampleErr.Error()
}

func (e *PgError) Field(field byte) string {
	if field == 'C' {
		return e.FieldCode
	}
	return ""
}

func (e *PgError) IntegrityViolation() bool {
	return e.IntegrityViolationValue
}
