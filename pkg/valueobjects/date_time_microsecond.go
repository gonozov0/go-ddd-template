package valueobjects

import (
	"time"

	"go-ddd-template/pkg/parseutils"

	"errors"
	"fmt"
)

type DateTimeMicrosecond time.Time

var ErrDateTimeMicrosecondValidation = errors.New("date time validation error")

var EmptyDateTimeMicrosecond = DateTimeMicrosecond(time.Time{})

func NewDateTimeMicrosecondNow() DateTimeMicrosecond {
	return NewDateTimeMicrosecond(time.Now().UTC())
}

func NewRequiredDateTimeMicrosecond(t time.Time) (DateTimeMicrosecond, error) {
	if t.IsZero() {
		return DateTimeMicrosecond{}, fmt.Errorf("%w: empty value", ErrDateTimeMicrosecondValidation)
	}

	return DateTimeMicrosecond(t.Truncate(time.Microsecond)), nil
}

func NewDateTimeMicrosecond(t time.Time) DateTimeMicrosecond {
	return DateTimeMicrosecond(t.Truncate(time.Microsecond))
}

func (dtm DateTimeMicrosecond) ToTime() time.Time {
	return time.Time(dtm)
}

func (dtm DateTimeMicrosecond) Equal(other DateTimeMicrosecond) bool {
	return dtm.ToTime().Equal(other.ToTime())
}

func (dtm DateTimeMicrosecond) IsZero() bool {
	return dtm.ToTime().IsZero()
}

func (dtm DateTimeMicrosecond) Before(other DateTimeMicrosecond) bool {
	return dtm.ToTime().Before(other.ToTime())
}

func (dtm DateTimeMicrosecond) After(other DateTimeMicrosecond) bool {
	return dtm.ToTime().After(other.ToTime())
}

func (dtm DateTimeMicrosecond) Sub(other DateTimeMicrosecond) time.Duration {
	return dtm.ToTime().Sub(other.ToTime())
}

func (dtm DateTimeMicrosecond) Add(d time.Duration) DateTimeMicrosecond {
	return NewDateTimeMicrosecond(dtm.ToTime().Add(d))
}

func (dtm DateTimeMicrosecond) EqualInHours(equalTo DateTimeMicrosecond) bool {
	t1, t2 := dtm.ToTime().UTC(), equalTo.ToTime().UTC()

	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() && t1.Hour() == t2.Hour()
}

func (dtm DateTimeMicrosecond) AddDate(years, months, days int) DateTimeMicrosecond {
	return DateTimeMicrosecond(dtm.ToTime().AddDate(years, months, days))
}

func (dtm DateTimeMicrosecond) TimePointer() *time.Time {
	if dtm.IsZero() {
		return nil
	}

	return parseutils.ToPointer(dtm.ToTime())
}
