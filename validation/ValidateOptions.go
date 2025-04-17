package validation

type ValidateOptions struct {
	ShortCircuit bool
}

var DefaultValidateOptions = ValidateOptions{
	ShortCircuit: false,
}
