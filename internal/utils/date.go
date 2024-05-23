package utils

import "time"

func DateFormat(layout string, d int64) string {
	intTime := int64(d)
	t := time.Unix(intTime, 0)
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}
