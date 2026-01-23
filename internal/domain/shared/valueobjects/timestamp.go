package valueobjects

import (
	"time"
)

type Timestamp time.Time

func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UTC())
}

func NewTimestampNow() Timestamp {
	return Timestamp(NewTimestamp(time.Now().UTC()))
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

func (t Timestamp) Format(layout string) string {
	return time.Time(t).Format(layout)
}
