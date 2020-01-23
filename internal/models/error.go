package models

type HttpError struct {
	StatusCode int
	StringErr string
}

func (e *HttpError) Error () string{
	return e.StringErr
}