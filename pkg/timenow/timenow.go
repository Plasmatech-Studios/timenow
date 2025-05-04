package timenow

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Config holds the options for time formatting and timezone offset.
type Config struct {
	TimezoneOffsetMinutes int    // minutes east of UTC
	Format                string // "unix", "rfc3339", "ansic", "stamp", "stampmilli"
}

// GetTimeNow returns the current time per cfg.Format and cfg.TimezoneOffsetMinutes.
// On error, it falls back to unix.
func GetTimeNow(cfg Config) (string, error) {
	// We only apply the offset for non-unix outputs
	now := time.Now().UTC()
	if strings.ToLower(cfg.Format) != "unix" {
		off := time.Duration(cfg.TimezoneOffsetMinutes) * time.Minute
		now = now.Add(off)
	}

	switch strings.ToLower(cfg.Format) {
	case "unix":
		// ignore now + offset entirely
		return fmt.Sprintf("%d", time.Now().Unix()), nil
	case "rfc3339":
		return now.Format(time.RFC3339), nil
	case "ansic":
		return now.Format(time.ANSIC), nil
	case "stamp":
		return now.Format(time.Stamp), nil
	case "stampmilli":
		return now.Format(time.StampMilli), nil
	default:
		return fmt.Sprintf("%d", time.Now().Unix()), errors.New("unknown format, defaulting to unix")
	}
}

// ParseTZ parses strings like "-8:30", "8.45", "5.5", "5", "-5", into minutes offset.
// It supports:
//   - colon syntax (h:mm)
//   - dot syntax as fractional hours if single digit (h.f), or
//     minute-literal if two digits (h.mm).
//
// Range enforced: [-12h, +14h].
func ParseTZ(input string) (int, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return 0, nil
	}

	// extract sign
	sign := 1
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		sign = -1
		s = s[1:]
	}

	// fractional or minute-literal via dot
	if strings.Contains(s, ".") {
		parts := strings.SplitN(s, ".", 2)
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid tz format: %q", input)
		}
		hStr, frac := parts[0], parts[1]

		hours, err := strconv.Atoi(hStr)
		if err != nil {
			return 0, fmt.Errorf("invalid hour in %q", input)
		}

		var minutes int
		if len(frac) == 1 {
			// fractional hours
			f, err := strconv.ParseFloat("0."+frac, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid fraction in %q", input)
			}
			minutes = int(f*60 + 0.5)
		} else if len(frac) == 2 {
			// literal minutes
			m, err := strconv.Atoi(frac)
			if err != nil {
				return 0, fmt.Errorf("invalid minute in %q", input)
			}
			if m < 0 || m >= 60 {
				return 0, fmt.Errorf("minute must be 0–59 in %q", input)
			}
			minutes = m
		} else {
			return 0, fmt.Errorf("invalid tz format: %q", input)
		}

		total := sign * (hours*60 + minutes)
		if total < -12*60 || total > 14*60 {
			return 0, fmt.Errorf("timezone offset must be between -12h and +14h")
		}
		return total, nil
	}

	// colon syntax h:mm
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid tz format: %q", input)
		}
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid hour in %q", input)
		}
		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid minute in %q", input)
		}
		if minutes < 0 || minutes >= 60 {
			return 0, fmt.Errorf("minute must be 0–59 in %q", input)
		}
		total := sign * (hours*60 + minutes)
		if total < -12*60 || total > 14*60 {
			return 0, fmt.Errorf("timezone offset must be between -12h and +14h")
		}
		return total, nil
	}

	// bare hours only
	hours, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid hour in %q", input)
	}
	total := sign * hours * 60
	if total < -12*60 || total > 14*60 {
		return 0, fmt.Errorf("timezone offset must be between -12h and +14h")
	}
	return total, nil
}
