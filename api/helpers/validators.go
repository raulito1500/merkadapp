package helpers

import (
	"errors"
	"strings"
)

func ValidateMandatory(d string) error {
	if d == "" {
		return errors.New("%s field is mandatory")
	}
	return nil
}

func ValidateIntNonZeroPositive(d uint16) error {
	if d <= 0 {
		return errors.New("%s field is mandatory and must be greater than zero")
	}
	return nil
}

func ValidateFloatNonZeroPositive(d float32) error {
	if d <= 0 {
		return errors.New("%s field is mandatory and must be greater than zero")
	}
	return nil
}

func ValidateInEnum(d string, enum []string) error {
	for _, e := range enum {
		if e == d {
			return nil
		}
	}
	return errors.New("%s field must be a valid data [" + strings.Join(enum, ",") + "]")
}
