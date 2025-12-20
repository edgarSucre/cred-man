package httputils

import (
	"encoding/json"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	d.Time = t

	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(time.DateOnly))
}
