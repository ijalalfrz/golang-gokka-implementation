package date

import (
	"time"
)

// Date constants.
const (
	SQLDateTimeFormat string = "2006-01-02 15:04:05.999"
	RFC3339Millis     string = "2006-01-02T15:04:05.999Z"
	HumanDateFormat   string = "_2 Jan 2006"
)

// SQLDateTime will return sql date time format
func SQLDateTime(layout string, value string, l *time.Location) string {
	t, _ := time.Parse(layout, value)
	return t.In(l).Format(SQLDateTimeFormat)
}
