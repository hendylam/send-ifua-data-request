package apperror

// AppError interface represents service logical error behaviours which inherits from error interface
type AppError interface {
	error
	Type() int
}

// AppError struct
type appError struct {
	err     string
	errType int
}

// Newd func, create a AppError with 0 as a type
func Newd(err string) AppError {
	return New(err, 0)
}

// New func, create a AppError
func New(err string, errType int) AppError {
	return &appError{err: err, errType: errType}
}

// Error func
func (e *appError) Type() int {
	return e.errType
}

func (e *appError) Error() string {
	return e.err
}
