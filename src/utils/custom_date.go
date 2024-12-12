package utils

import (
	"fmt"
	"strings"
	"time"
)

/*
Type Definition for Custom Date
*/

type CustomDate struct {
	*time.Time
}

/*
Implement UnmarshalJSON for CustomDate type
*/

func (ct *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // Erase quotes
	if s == "null" || s == "" {       // Handle null or empty string
		ct.Time = &time.Time{}
		return nil
	}

	var t time.Time
	// Try to parse with YYYY-MM-DD format
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		ct.Time = &t
		return nil
	}
	// If Failed, Try to parse with RFC3339 format (default json.Unmarshal)
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		ct.Time = &t
		return nil
	}

	return fmt.Errorf("cannot parse date string '%s': %w", s, err)
}

/*
Optional
*/

func (ct CustomDate) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, ct.Time.Format("2006-01-02")), nil
}
