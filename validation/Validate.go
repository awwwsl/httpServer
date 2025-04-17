package validation

func Validate[T any](value T, options ValidateOptions, validators ...func(T) *ValidateError) (bool, []*ValidateError) {
	results := make([]*ValidateError, 0)

	for _, validator := range validators {
		if err := validator(value); err != nil {
			results = append(results, err)
			if options.ShortCircuit {
				return false, results
			}
		}
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, results
}
