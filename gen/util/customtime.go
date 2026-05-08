package util

import (
	"strconv"
	"time"
)

type DTTime time.Time

func (t *DTTime) UnmarshalJSON(b []byte) error {
	millis, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = DTTime(time.UnixMilli(millis))
	return nil
}
