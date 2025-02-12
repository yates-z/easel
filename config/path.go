package config

import "strings"

type Kind uint8

const (
	Index = iota
	Key
)

// Part represents a key or index in a path.
type Part struct {
	value string
	kind  Kind
}

// parsePath parses a dot-separated path into a slice of keys.
func parsePath(path string) []*Part {
	var keys []*Part
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
			// handle array index
			keys = append(keys, &Part{value: strings.Split(part, "[")[0], kind: Key})
			keys = append(keys, &Part{value: strings.TrimSuffix(strings.Split(part, "[")[1], "]"), kind: Index})
		} else {
			keys = append(keys, &Part{value: part, kind: Key})
		}
	}
	return keys
}
