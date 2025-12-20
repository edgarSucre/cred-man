package terror

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	ErrorType int

	Error struct {
		T      ErrorType
		Code   string
		Detail string
	}
)

const (
	Timeout ErrorType = iota
	Internal
	NotFound
	Validation
)

func (t ErrorType) New(code, detail string) Error {
	return Error{
		Code:   code,
		Detail: detail,
		T:      t,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %s, detail: %s", e.Code, e.Detail)
}

func (e Error) HttpStatus() int {
	switch e.T {
	case Timeout:
		return http.StatusRequestTimeout
	case NotFound:
		return http.StatusNotFound
	case Validation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (e Error) Raw() []byte {
	template := `{"code": "%s", "detail": "%s"}`

	msg := fmt.Sprintf(template, e.Code, e.Detail)

	return []byte(msg)
}

func Unwrap(err error) error {

	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped != nil {
			err = unwrapped
			continue
		}

		return err
	}
}
