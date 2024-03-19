package helpers

import (
	"errors"
	"testing"
)

func TestValidateMandatory(t *testing.T) {
	tables := []struct {
		dato     string
		expected error
	}{{
		dato:     "Texto",
		expected: nil,
	}, {
		dato:     "text",
		expected: errors.New("%s field is mandatory"),
	}}

	for _, test := range tables {
		val := ValidateMandatory(test.dato)
		if test.expected != nil && val == nil {
			t.Errorf("Error, got nil expected %s", test.expected)
		} else if test.expected != nil && !errors.Is(val, test.expected) {
			t.Errorf("Error, got %s expected %s", val, test.expected)
		}
	}
}

func TestValidateIntNonZeroPositive(t *testing.T) {
	tables := []struct {
		dato     uint16
		expected error
	}{{
		dato:     1,
		expected: nil,
	}, {
		dato:     0,
		expected: errors.New("%s field is mandatory and must be greater than zero"),
	}}

	for _, test := range tables {
		val := ValidateIntNonZeroPositive(test.dato)
		if test.expected != nil && val == nil {
			t.Errorf("Error, got nil expected %v", test.expected)
		} else if test.expected != nil && !errors.Is(val, test.expected) {
			t.Errorf("Error, got %v expected %v", val, test.expected)
		}
	}
}
