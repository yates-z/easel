package rbac

import (
	"bufio"
	"encoding/csv"
	"os"
)

// Adapter defines the interface for a persistence adapter.
type Adapter interface {
	LoadPolicy() ([]Policy, error)
	SavePolicy(policies []Policy) error
}

// CSVAdapter implements the Adapter interface using a CSV file for persistence.
type CSVAdapter struct {
	filePath string
}

// NewCSVAdapter creates a new CSVAdapter with the given file path.
func NewCSVAdapter(filePath string) *CSVAdapter {
	return &CSVAdapter{filePath: filePath}
}

// LoadPolicy loads policies from a CSV file using buffered reading.
func (a *CSVAdapter) LoadPolicy() ([]Policy, error) {
	file, err := os.Open(a.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Buffered reader for handling large files.
	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var policies []Policy
	for _, record := range records {
		if len(record) == 3 {
			policies = append(policies, Policy{
				Role:     record[0],
				Resource: record[1],
				Action:   record[2],
			})
		}
	}
	return policies, nil
}

// SavePolicy saves policies to a CSV file.
func (a *CSVAdapter) SavePolicy(policies []Policy) error {
	file, err := os.Create(a.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, p := range policies {
		if err := writer.Write([]string{p.Role, p.Resource, p.Action}); err != nil {
			return err
		}
	}
	return nil
}
