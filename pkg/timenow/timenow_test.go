package timenow

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseTZ(t *testing.T) {
	tests := []struct {
		input   string
		wantMin int
		wantErr bool
	}{
		{"", 0, false},                 // 0h
		{"0", 0, false},                // 0h
		{"+0", 0, false},               // 0h
		{"-0", 0, false},               // 0h
		{"5", 5 * 60, false},           // 5h
		{"-5", -5 * 60, false},         // -5h
		{"5:30", 5*60 + 30, false},     // 5h30m
		{"5.5", 5*60 + 30, false},      // 5h30m
		{"8.4", 8*60 + 24, false},      // 8h24m
		{"8.40", 8*60 + 40, false},     // 8h40m
		{"-8.4", -(8*60 + 24), false},  // -8h24m
		{"-8.40", -(8*60 + 40), false}, // -8h40m
		{"-8:30", -(8*60 + 30), false}, // -8h30m
		{"14:00", 14 * 60, false},      // 14h
		{"-12:00", -12 * 60, false},    // -12h
		{"14:01", 0, true},             // above +14h
		{"-12:01", 0, true},            // below -12h
		{"foo", 0, true},               // unparsable
		{"5:60", 0, true},              // invalid minutes
		{"5:ab", 0, true},              // invalid minutes
	}

	for _, tt := range tests {
		got, err := ParseTZ(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseTZ(%q): expected error, got none (got %d)", tt.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseTZ(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.wantMin {
			t.Errorf("ParseTZ(%q): want %d minutes, got %d", tt.input, tt.wantMin, got)
		}
	}
}

func TestGetTimeNowUnixAndFallback(t *testing.T) {
	// 1) Unix format
	s, err := GetTimeNow(Config{Format: "unix", TimezoneOffsetMinutes: 123})
	if err != nil {
		t.Errorf("GetTimeNow(unix): unexpected error: %v", err)
	}
	if _, err := strconv.ParseInt(s, 10, 64); err != nil {
		t.Errorf("GetTimeNow(unix) returned non‐numeric %q", s)
	}

	// 2) Fallback on unknown format
	s2, err2 := GetTimeNow(Config{Format: "not_a_format", TimezoneOffsetMinutes: 0})
	if err2 == nil {
		t.Error("GetTimeNow(unknown) expected error, got nil")
	}
	if !strings.Contains(err2.Error(), "unknown format") {
		t.Errorf("GetTimeNow(unknown) error = %q; want contains \"unknown format\"", err2.Error())
	}
	if _, err := strconv.ParseInt(s2, 10, 64); err != nil {
		t.Errorf("GetTimeNow(unknown) returned non‐numeric %q", s2)
	}
}

func TestGetTimeNowFormats(t *testing.T) {
	tests := []struct {
		format string
		layout string
	}{
		{"rfc3339", time.RFC3339},
		{"ansic", time.ANSIC},
		{"stamp", time.Stamp},
		{"stampmilli", time.StampMilli},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			cfg := Config{Format: tt.format, TimezoneOffsetMinutes: 0}
			out, err := GetTimeNow(cfg)
			if err != nil {
				t.Fatalf("GetTimeNow(%q): unexpected error: %v", tt.format, err)
			}
			// Ensure the output parses under the expected layout
			if _, perr := time.Parse(tt.layout, out); perr != nil {
				t.Errorf("GetTimeNow(%q) = %q; failed to parse with layout %q: %v",
					tt.format, out, tt.layout, perr)
			}
		})
	}
}
