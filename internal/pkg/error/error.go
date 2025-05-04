package error

import "net/http"

type FieldError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (fe FieldError) Error() string {
	return fe.Reason
}

type CustomError struct {
	StatusCode int          `json:"status_code"`
	Message    string       `json:"message"`
	Errors     []FieldError `json:"errors,omitempty"`
}

func (ce CustomError) Error() string {
	return ce.Message
}

func Err422(msg string) CustomError {
	return CustomError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    msg,
	}
}

func Err400(msg string) CustomError {
	return CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
	}
}

func Err401(msg string) CustomError {
	return CustomError{
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
	}
}

func Err500(msg string) CustomError {
	return CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
	}
}
