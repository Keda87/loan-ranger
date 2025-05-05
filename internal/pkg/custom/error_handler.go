package custom

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
)

func NewErrorHandler(err error, c echo.Context) {
	var (
		code    = http.StatusInternalServerError
		message = "Internal Server Error"
		errs    []pkgerr.FieldError
	)

	switch t := err.(type) {
	case pkgerr.CustomError:
		message = t.Message
		code = t.StatusCode
	case validator.ValidationErrors:
		code = http.StatusBadRequest
		message = "Validation Error"
		for _, valErr := range t {
			errs = append(errs, pkgerr.FieldError{Field: valErr.Field(), Reason: messageForTag(valErr)})
		}
	case *echo.HTTPError:
		code = t.Code
		message = fmt.Sprint(t.Message)
	default:
		message = t.Error()
	}

	errResp := payload.ResponseError[pkgerr.CustomError]{
		Error: pkgerr.CustomError{
			StatusCode: code,
			Message:    message,
		},
	}

	if len(errs) > 0 {
		errResp.Error.Errors = errs
	}

	_ = c.JSON(code, errResp)
}

func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "max":
		return fmt.Sprintf("Reach maximum %s", fe.Param())
	case "min":
		return fmt.Sprintf("Does not meet minimum %s", fe.Param())
	}
	return fe.Error()
}
