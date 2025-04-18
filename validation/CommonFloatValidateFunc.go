package validation

import "strconv"

type FloatValidateFunc func(float64) *ValidateError
type floatValidator struct{}

// Float is a struct that provides validation functions for floats.
var Float floatValidator

func (floatValidator) NotEqualTo(value float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v == value {
			return &ValidateError{Reason: "Value can't equal to " + strconv.FormatFloat(value, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotGreaterThan(value float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v > value {
			return &ValidateError{Reason: "Value can't be greater than " + strconv.FormatFloat(value, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotGreaterOrEqualTo(value float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v >= value {
			return &ValidateError{Reason: "Value can't be greater or equal to " + strconv.FormatFloat(value, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotLessThan(value float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v < value {
			return &ValidateError{Reason: "Value can't be less than " + strconv.FormatFloat(value, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotLessOrEqualTo(value float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v <= value {
			return &ValidateError{Reason: "Value can't be less or equal to " + strconv.FormatFloat(value, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotEqualToAny(values ...float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		for _, value := range values {
			if v == value {
				return &ValidateError{Reason: "Value can't be in " + strconv.FormatFloat(value, 'f', -1, 64)}
			}
		}
		return nil
	}
}

func (floatValidator) Between(min, max float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v < min || v > max {
			return &ValidateError{Reason: "Value must be between " + strconv.FormatFloat(min, 'f', -1, 64) + " and " + strconv.FormatFloat(max, 'f', -1, 64)}
		}
		return nil
	}
}
func (floatValidator) NotBetween(min, max float64) FloatValidateFunc {
	return func(v float64) *ValidateError {
		if v >= min && v <= max {
			return &ValidateError{Reason: "Value must not be between " + strconv.FormatFloat(min, 'f', -1, 64) + " and " + strconv.FormatFloat(max, 'f', -1, 64)}
		}
		return nil
	}
}
