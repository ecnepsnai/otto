// Package environ is a custom environment variable utilities for otto
package environ

import "fmt"

// Variable describes a single variable
type Variable struct {
	Key    string
	Value  string
	Secret bool
}

// New create a new variable
func New(key string, value string) Variable {
	return Variable{
		Key:   key,
		Value: value,
	}
}

// FromMap return an array of variables form the given map
func FromMap(m map[string]string) []Variable {
	vars := []Variable{}
	for k, v := range m {
		vars = append(vars, New(k, v))
	}
	return vars
}

// Merge will merge the given two variable slices. Objects from `original` will be replaces by any from `adding`
// if there are duplicate keys
func Merge(original []Variable, adding []Variable) []Variable {
	keyIdxMap := map[string]int{}
	for i, v := range original {
		keyIdxMap[v.Key] = i
	}

	newVars := original
	for _, v := range adding {
		idx, existing := keyIdxMap[v.Key]
		if existing {
			newVars[idx] = v
		} else {
			newVars = append(newVars, v)
		}
	}

	return newVars
}

// Map will return a mapping of key => value from the given slice of variables
func Map(vars []Variable) map[string]string {
	m := map[string]string{}
	for _, v := range vars {
		m[v.Key] = v.Value
	}
	return m
}

// ReservedKeys keys reserved by the otto system
var ReservedKeys = []string{
	"OTTO_SERVER_VERSION",
	"OTTO_SERVER_URL",
	"OTTO_HOST_ADDRESS",
	"OTTO_HOST_PORT",
	"OTTO_HOST_PSK",
}

// Validate will return if any of the variables are invalid
func Validate(vars []Variable) error {
	for _, v := range vars {
		for _, key := range ReservedKeys {
			if v.Key == key {
				return fmt.Errorf("key is reserved by the Otto system")
			}
		}
	}
	return nil
}
