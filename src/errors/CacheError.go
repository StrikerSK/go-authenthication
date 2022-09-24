package errors

type CacheError struct {
	message string
}

func NewCacheError() CacheError {
	return CacheError{}
}

func (r CacheError) FromMessage(message string) CacheError {
	r.message = message
	return r
}

func (r CacheError) FromError(err error) CacheError {
	r.message = err.Error()
	return r
}

func (r CacheError) Error() string {
	return r.message
}
