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
	Rejection
)

func (t ErrorType) New(code, detail string) Error {
	return Error{
		Code:   code,
		Detail: detail,
		T:      t,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf(`{"code": "%s", "detail": "%s"}`, e.Code, e.Detail)
}

func (e Error) HttpStatus() int {
	switch e.T {
	case Timeout:
		return http.StatusRequestTimeout
	case NotFound:
		return http.StatusNotFound
	case Validation:
		return http.StatusBadRequest
	case Rejection:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (e Error) Raw() []byte {
	return []byte(e.Error())
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

func ToInternal(e error) error {
	ue := Unwrap(e)

	if err, ok := ue.(Error); ok {
		return Internal.New(err.Code, err.Detail)
	}

	return Internal.New("internal-error", e.Error())
}
