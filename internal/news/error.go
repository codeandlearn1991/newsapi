package news

// CustomError represents the error state of
// database error.
type CustomError struct {
	err        error
	httpStatus int
}

// NewCustomError returns an instance of customer error.
func NewCustomError(err error, httpStatus int) *CustomError {
	return &CustomError{
		err:        err,
		httpStatus: httpStatus,
	}
}

// Error implements the error interface.
func (ce CustomError) Error() string {
	return ce.err.Error()
}

// Unwrap the underlying error.
func (ce CustomError) Unwrap() error {
	return ce.err
}

// HTTPStatusCode representing the error.
func (ce CustomError) HTTPStatusCode() int {
	return ce.httpStatus
}
