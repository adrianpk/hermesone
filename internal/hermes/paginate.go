package hermes

import (
	"fmt"
)

// Paginate returns a slice of data for the specified page and page size.
func Paginate[T any](data []T, page, pageSize int) ([]T, error) {
	if page < 1 {
		return nil, fmt.Errorf("page number must be greater than 0")
	}

	if pageSize < 1 {
		return nil, fmt.Errorf("page size must be greater than 0")
	}

	start := (page - 1) * pageSize
	if start >= len(data) {
		return nil, fmt.Errorf("page %d out of range", page)
	}

	end := start + pageSize
	if end > len(data) {
		end = len(data)
	}

	return data[start:end], nil
}

// TotalPages returns the total number of pages for the given data and page size.
func TotalPages[T any](data []T, pageSize int) int {
	if pageSize < 1 {
		return 0
	}

	return (len(data) + pageSize - 1) / pageSize
}
