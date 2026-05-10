package helpers

import "strconv"

func GetString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GetFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

func ParseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}
