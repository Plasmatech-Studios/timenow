# timenow

A simple, zero‑dependency Go CLI (and package) for fetching the current time in various formats and timezones. By default, `timenow` returns the Unix timestamp (seconds since epoch) and copies it to the clipboard. It can also be imported as a library to retrieve formatted time strings programmatically.

---

## Features

* **Default Unix timestamp** with clipboard copy
* Custom **timezone offsets** between UTC−12:00 and UTC+14:00 (supports `h:mm`, fractional `h.f`, and `h.mm` notation)
* Multiple **output formats**:

  * `unix` (default)
  * `rfc3339` (ISO‑8601)
  * `ansic` (Go's ANSIC layout)
  * `stamp` (basic timestamp layout)
  * `stampmilli` (millisecond precision)
* Lightweight library API for embedding in other Go programs
* Graceful error handling with sensible defaults

---

## Repository Layout

```bash
timenow/                   ← Go library package
  timenow.go               ← core logic (ParseTZ, GetTimeNow)
  timenow_test.go          ← unit tests (table‑driven)

cmd/
  timenow/                 ← CLI entrypoint
    main.go                ← flag parsing, clipboard, I/O

go.mod
README.md
```

---

## Installation

### Clone and Build

1. **Clone the repo** (or add it as a module dependency):

   ```bash
   git clone https://github.com/Plasmatech-Studios/timenow.git
   cd timenow
   ```
2. **Initialize and download dependencies**:

   ```bash
   go mod tidy
   ```
3. **Build the CLI**:

   ```bash
   go build -o timenow ./cmd/timenow
   ```

### Local Installation

After building, move the binary into your `$PATH`. For example:

```bash
# Move to ~/bin (ensure ~/bin is in your PATH):
mv timenow ~/bin/

# Or, install system-wide (optional, may require sudo):
sudo mv timenow /usr/local/bin/
```

---

## Usage

```bash
# Default: Unix timestamp → stdout & clipboard
$ timenow
1586515200

# Named output formats:
$ timenow --format=rfc3339
2025-05-04T14:07:12Z

# Shorthand flags:
$ timenow -f ansic
Sun May  4 14:07:12 2025

# Custom timezone offsets:
$ timenow --format=stampmilli --tz +5:30
May  4 19:37:12.123

# Ignored tz for unix format:
$ timenow -f unix -tz -8:30
warning: timezone flag "-8:30" ignored when format=unix
1586515200
```

### Flags

| Flag         | Alias | Description                                                                              | Default |
| ------------ | ----- | ---------------------------------------------------------------------------------------- | ------- |
| `--format`   | `-f`  | Output format: `unix`, `rfc3339`, `ansic`, `stamp`, `stampmilli`                         | `unix`  |
| `--timezone` | `-tz` | Timezone offset (`h:mm`, `h.f`, `h.mm`) in range UTC−12:00…UTC+14:00; ignored for `unix` | `UTC`   |

---

## Library API

Import:

```go
import "github.com/Plasmatech-Studios/timenow"
```

### Config struct

```go
// Config holds formatting options:

// TimezoneOffsetMinutes: minutes east of UTC (e.g. 330 for +05:30)
// Format: one of "unix", "rfc3339", "ansic", "stamp", "stampmilli"

type Config struct {
    TimezoneOffsetMinutes int
    Format                string
}
```

### GetTimeNow

```go
func GetTimeNow(cfg Config) (string, error)
```

* Returns the formatted time string.
* On unknown `Format`, returns Unix seconds and an error.

### ParseTZ

```go
func ParseTZ(input string) (int, error)
```

* Parses `input` like `-8:30`, `5.5`, `8.40` into minutes offset.
* Ensures the offset is between −12h and +14h.

---

## Testing

Run the table‑driven tests:

```bash
go test ./timenow
```

---

## Contributing

1. Fork the repo and create a feature branch.
2. Write tests for new behavior.
3. Open a pull request with a clear description.

---

## License

This project is released under the MIT License. See [LICENSE](LICENSE) for details.
