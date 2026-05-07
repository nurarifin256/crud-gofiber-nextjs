package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func CellRef(row, col int) string {
	colName, _ := excelize.ColumnNumberToName(col) // 1->A, 2->B, ...
	return fmt.Sprintf("%s%d", colName, row)
}

func SafeName(s string) string {
	if s == "" {
		return "export"
	}
	// Simple cleaner; sesuaikan jika perlu
	forbidden := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	out := s
	for _, ch := range forbidden {
		out = strings.ReplaceAll(out, ch, "_")
	}
	return out
}
func ParseDateFlexible(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("invalid date")
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05",
		"01/02/2006", // MM/DD/YYYY
		"02/01/2006", // DD/MM/YYYY
		"2006/01/02",
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.UTC); err == nil {
			y, m, d := t.Date()
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, errors.New("invalid date")
}