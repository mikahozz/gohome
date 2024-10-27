package spot

import "fmt"

// NoDataError represents an error when no data is available
type NoDataError struct {
	Code string
	Text string
}

// Error implements the error interface for NoDataError
func (e *NoDataError) Error() string {
	return fmt.Sprintf("No data available: %s (Code: %s)", e.Text, e.Code)
}
