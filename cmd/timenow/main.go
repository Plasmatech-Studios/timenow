package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Plasmatech-Studios/timenow/timenow/pkg/timenow"
	"github.com/atotto/clipboard"
)

func main() {
	// Flags
	tzFlag := flag.String("timezone", "", "Timezone offset, e.g. -8:30, 8.45, +5:30 (range -12..+14)")
	tzShort := flag.String("tz", "", "Shorthand for --timezone")
	fmtFlag := flag.String("format", "unix", "Output format: unix, rfc3339, ansic, stamp, stampmilli")
	fmtShort := flag.String("f", "unix", "Shorthand for --format")

	flag.Parse()

	// Pick up any tz input (long form wins)
	tzInput := *tzFlag
	if tzInput == "" {
		tzInput = *tzShort
	}

	// Determine format (long form wins)
	format := *fmtFlag
	if format == "unix" && *fmtShort != "unix" {
		format = *fmtShort
	}

	// Warn if user passed a timezone but is using unix‐seconds
	if format == "unix" && tzInput != "" {
		fmt.Fprintf(os.Stderr, "warning: timezone flag %q ignored when format=unix\n", tzInput)
	}

	// Parse timezone offset (only really used for non‐unix formats)
	offsetMin := 0
	if tzInput != "" {
		if m, err := timenow.ParseTZ(tzInput); err != nil {
			fmt.Fprintf(os.Stderr, "warning: invalid timezone %q, using UTC: %v\n", tzInput, err)
		} else {
			offsetMin = m
		}
	}

	cfg := timenow.Config{
		TimezoneOffsetMinutes: offsetMin,
		Format:                format,
	}

	// Generate and output
	out, err := timenow.GetTimeNow(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "note: %v\n", err)
	}
	fmt.Println(out)

	// Copy to clipboard
	if err := clipboard.WriteAll(out); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not copy to clipboard: %v\n", err)
	}
}
