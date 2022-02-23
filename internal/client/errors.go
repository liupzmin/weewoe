package client

// Error represents an error.
type Error string

// Error returns the error text.
func (e Error) Error() string {
	return string(e)
}
