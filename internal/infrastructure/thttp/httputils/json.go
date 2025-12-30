package httputils

import (
	"io"

	"github.com/mailru/easyjson"
)

func Unmarshal(r io.ReadCloser, v easyjson.Unmarshaler) error {
	defer r.Close()

	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := easyjson.Unmarshal(raw, v); err != nil {
		return err
	}

	return nil
}

func Marshal(w io.Writer, v easyjson.Marshaler) error {
	raw, err := easyjson.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(raw)
	return err
}
