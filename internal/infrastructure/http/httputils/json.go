package httputils

import (
	"io"

	"github.com/edgarSucre/mye"
	"github.com/mailru/easyjson"
)

func Unmarshal(r io.ReadCloser, v easyjson.Unmarshaler) (err error) {
	defer func() {
		if cErr := r.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	raw, err := io.ReadAll(r)
	if err != nil {
		return mye.Wrap(
			err,
			mye.CodeInvalid,
			"reading_request_failed",
			"error reading request",
		)
	}

	if err := easyjson.Unmarshal(raw, v); err != nil {
		return mye.Wrap(
			err,
			mye.CodeInvalid,
			"json_unmarshall_request_failed",
			"error unmarshall request",
		)
	}

	return nil
}

func Marshal(w io.Writer, v easyjson.Marshaler) error {
	raw, err := easyjson.Marshal(v)
	if err != nil {
		return mye.Wrap(
			err,
			mye.CodeInternal,
			"marshal_response_failed",
			"failed to marshall response",
		).WithAttribute("response", v)
	}

	_, err = w.Write(raw)
	return mye.Wrap(
		err,
		mye.CodeInternal,
		"writing_response_failed",
		"failed to write response",
	).WithAttribute("response", v)
}
