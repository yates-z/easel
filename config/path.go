package config

import "strings"

type Kind uint8

const (
	Index = iota
	Key
)

// Part represents a key or index in a path.
type Part struct {
	value  string
	kind   Kind
	isLeaf bool

	parent *Part
}

// parsePath parses a dot-separated path into a slice of keys.
func parsePath(path string) []*Part {
	var keys []*Part
	parts := strings.Split(path, ".")
	for _, part := range parts {
		for strings.Contains(part, "[") && strings.HasSuffix(part, "]") {

			indexStart := strings.Index(part, "[")
			indexEnd := strings.Index(part, "]")
			key := part[:indexStart]
			if key != "" {
				keys = append(keys, &Part{value: key, kind: Key})
			}
			index := part[indexStart+1 : indexEnd]
			keys = append(keys, &Part{value: index, kind: Index})
			part = part[indexEnd+1:]
		}
		if part != "" {
			keys = append(keys, &Part{value: part, kind: Key})
		}
	}
	keys[len(keys)-1].isLeaf = true

	for i := len(keys) - 1; i > 0; i-- {
		keys[i].parent = keys[i-1]
	}

	return keys
}
