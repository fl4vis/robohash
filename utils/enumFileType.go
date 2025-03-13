package utils

import (
	"fmt"
	"strings"
)

type ExtensionsType struct {
	ValidValues []string
	Value       string
}

func (e *ExtensionsType) String() string {
	return e.Value
}

func (e *ExtensionsType) Set(s string) error {
	s = strings.ToLower(s)

	for _, valid := range e.ValidValues {
		if s == valid {
			e.Value = s
			return nil
		}
	}

	return fmt.Errorf("invalid extension type: %s, must be one of: %v", s, e.ValidValues)
}
