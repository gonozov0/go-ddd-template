package parseutils

import "time"

const (
	dateFormat = "2006-01-02"
)

func TruncateTimeMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute).UTC()
}

func StringToDate(strDate string) time.Time {
	date, err := time.Parse(dateFormat, strDate)
	if err != nil {
		return time.Time{}
	}

	return date
}

func DateToString(date time.Time) string {
	if date.IsZero() {
		return ""
	}

	return date.Format(dateFormat)
}

func TruncateToDate(datetime time.Time) time.Time {
	return time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		0, 0, 0, 0,
		datetime.Location(),
	)
}
