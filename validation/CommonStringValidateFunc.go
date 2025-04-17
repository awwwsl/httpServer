package validation

import (
	"regexp"
	"strconv"
	"unicode"
)

type StringValidateFunc func(string) *ValidateError
type stringValidator struct{}

var String stringValidator

func (stringValidator) NotEmptyOrWhiteSpace() StringValidateFunc {
	return func(value string) *ValidateError {
		for _, r := range value {
			if !unicode.IsSpace(r) {
				return nil
			}
		}
		return &ValidateError{Reason: "Value is empty or whitespace"}
	}
}

func (stringValidator) NotLongerThan(maxLength int) StringValidateFunc {
	return func(value string) *ValidateError {
		if len(value) > maxLength {
			return &ValidateError{Reason: "MaxLength is " + strconv.Itoa(maxLength) + " but was " + strconv.Itoa(len(value))}
		}
		return nil
	}
}

func (stringValidator) NotShorterThan(minLength int) StringValidateFunc {
	return func(value string) *ValidateError {
		if len(value) < minLength {
			return &ValidateError{Reason: "MinLength is " + strconv.Itoa(minLength) + " but was " + strconv.Itoa(len(value))}
		}
		return nil
	}
}

func (stringValidator) ByRegex(regexp regexp.Regexp) StringValidateFunc {
	return func(value string) *ValidateError {
		if !regexp.MatchString(value) {
			return &ValidateError{Reason: "Value does not match regex " + regexp.String()}
		}
		return nil
	}
}

func (stringValidator) ByRegexString(regex string) StringValidateFunc {
	regularExp := regexp.MustCompile(regex)
	return String.ByRegex(*regularExp)
}

func (stringValidator) NotEqualTo(value string) StringValidateFunc {
	return func(v string) *ValidateError {
		if v == value {
			return &ValidateError{Reason: "Value can't equal to " + value}
		}
		return nil
	}
}

func (stringValidator) NotEqualToAny(values ...string) StringValidateFunc {
	return func(value string) *ValidateError {
		for _, v := range values {
			if value == v {
				return &ValidateError{Reason: "Value can't equal to " + v}
			}
		}
		return nil
	}
}

func (stringValidator) NotContains(value string) StringValidateFunc {
	return func(v string) *ValidateError {
		if len(v) == 0 {
			return nil
		}
		if len(value) == 0 {
			return nil
		}
		if !regexp.MustCompile(value).MatchString(v) {
			return &ValidateError{Reason: "Value can't contain " + value}
		}
		return nil
	}
}

func (stringValidator) NotContainsAny(values ...string) StringValidateFunc {
	return func(value string) *ValidateError {
		for _, v := range values {
			if len(value) == 0 {
				return nil
			}
			if len(v) == 0 {
				return nil
			}
			if regexp.MustCompile(v).MatchString(value) {
				return &ValidateError{Reason: "Value can't contain " + v}
			}
		}
		return nil
	}
}
