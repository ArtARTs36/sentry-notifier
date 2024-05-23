package sentry

import "time"

type Time struct {
	time.Time
}

func (t Time) Human() string {
	return t.Time.Format("2006-01-02 15:04:05")
}
