package parser

import (
	"fmt"
	"time"
)

// DateLongFormat is same as time.RFC3339Nano
const DateLongFormat = time.RFC3339Nano

// DateShortFormat is a short date format with human-readable month element
const DateShortFormat = "2006-Jan-02"

// DateNumberFormat is a fully numeric short date format
const DateNumberFormat = "2006-01-02"

// ParseDate intends to parse all SQL date formats
func ParseDate(data string) (*time.Time, error) {
	t, err := time.Parse(DateLongFormat, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(time.RFC3339, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse("2006-01-02 15:04:05.999999 -0700 MST", data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(DateShortFormat, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(DateNumberFormat, data)
	if err == nil {
		return &t, nil
	}

	if data == "null" || data == "NULL" {
		return nil, nil
	}

	return nil, fmt.Errorf("not a date")
}
