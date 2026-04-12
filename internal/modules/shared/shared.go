package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
)

func HashToken(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

var nonDigit = regexp.MustCompile(`\D`)

func NormalizePhone(raw, countryCode string) string {
	if countryCode == "" {
		countryCode = "62"
	}
	p := nonDigit.ReplaceAllString(raw, "")
	if strings.HasPrefix(p, "0") {
		p = countryCode + strings.TrimPrefix(p, "0")
	} else if !strings.HasPrefix(p, countryCode) {
		p = countryCode + p
	}
	return p
}

func StartOfDay(t time.Time, loc *time.Location) time.Time {
	in := t.In(loc)
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, loc)
}

func EndOfDay(t time.Time, loc *time.Location) time.Time {
	in := t.In(loc)
	return time.Date(in.Year(), in.Month(), in.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
}
