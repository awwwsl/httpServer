package validation

import "strconv"

type IntegerValidateFunc func(int64) *ValidateError
type integerValidator struct{}

// Integer is a struct that provides validation functions for integers.
var Integer integerValidator

func (integerValidator) NotEqualTo(value int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v == value {
			return &ValidateError{Reason: "Value can't equal to " + strconv.FormatInt(value, 10)}
		}
		return nil
	}
}

func (integerValidator) NotGreaterThan(value int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v > value {
			return &ValidateError{Reason: "Value can't be greater than " + strconv.FormatInt(value, 10)}
		}
		return nil
	}
}
func (integerValidator) NotGreaterOrEqualTo(value int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v >= value {
			return &ValidateError{Reason: "Value can't be greater or equal to " + strconv.FormatInt(value, 10)}
		}
		return nil
	}
}

func (integerValidator) NotLessThan(value int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v < value {
			return &ValidateError{Reason: "Value can't be less than " + strconv.FormatInt(value, 10)}
		}
		return nil
	}
}
func (integerValidator) NotLessOrEqualTo(value int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v <= value {
			return &ValidateError{Reason: "Value can't be less or equal to " + strconv.FormatInt(value, 10)}
		}
		return nil
	}
}
func (integerValidator) NotEqualToAny(values ...int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		for _, value := range values {
			if v == value {
				return &ValidateError{Reason: "Value can't be in " + strconv.FormatInt(value, 10)}
			}
		}
		return nil
	}
}

func (integerValidator) Between(min, max int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v < min || v > max {
			return &ValidateError{Reason: "Value must be between " + strconv.FormatInt(min, 10) + " and " + strconv.FormatInt(max, 10)}
		}
		return nil
	}
}

func (integerValidator) NotBetween(min, max int64) IntegerValidateFunc {
	return func(v int64) *ValidateError {
		if v >= min && v <= max {
			return &ValidateError{Reason: "Value can't be between " + strconv.FormatInt(min, 10) + " and " + strconv.FormatInt(max, 10)}
		}
		return nil
	}
}
