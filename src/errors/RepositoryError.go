package errors

type RepositoryError struct {
	message string
}

func NewRepositoryError() RepositoryError {
	return RepositoryError{}
}

func (r RepositoryError) FromMessage(message string) RepositoryError {
	r.message = message
	return r
}

func (r RepositoryError) FromError(err error) RepositoryError {
	r.message = err.Error()
	return r
}

func (r RepositoryError) Error() string {
	return r.message
}
