package httputils

import (
	"encoding/json"
	"time"

	"github.com/edgarSucre/mye"
)

type Date struct {
	Value time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return mye.Wrap(
			err,
			mye.CodeInvalid,
			"date_unmarshal_failure",
			"failed to unmarshal date",
		).WithUserMsg("date incorrect format, use yyyy-mm-dd")
	}

	d.Value = t

	return nil
}

// just to comply with easyJson interface
func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value.Format(time.DateOnly))
}
