package httputils

import (
	"encoding/json"
	"time"
)

type Date struct {
	Value time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	d.Value = t

	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value.Format(time.DateOnly))
}
