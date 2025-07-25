# fehler-go

A diagnostic and error reporting library for Go, modeled after modern compilers like `rustc` and Zig's compiler. It produces clear, color-coded error messages with source code context and supports multiple output formats.

This is a port of [fehler](https://github.com/ciathefed/fehler) written in Zig.

![Go](https://img.shields.io/badge/Go-1.24.5-blue?style=flat-square%5C&logo=go)
![Tests](https://img.shields.io/github/actions/workflow/status/ciathefed/fehler-go/go.yml?label=Tests%20%F0%9F%A7%AA&style=flat-square)

## Features

* ✨ **Colorful Output**: ANSI color-coded diagnostics with severity-based coloring
* 📊 **Source Highlighting**: Highlights exact location of errors using carets (^) or tildes (\~)
* 📄 **Source Context**: Displays surrounding lines of code for better readability
* 🔍 **Smart Underlining**: Single-character and multi-character highlighting logic
* ⚖️ **Multiple Formats**: Supports `Fehler` (default), `GCC`, and `MSVC` styles
* 🎓 **Fluent API**: Builder pattern for constructing diagnostics
* ✅ **Convenient Helpers**: Shorthand functions for common cases

## Install

```bash
go get github.com/ciathefed/fehler-go
```

## Quick Start

```go
package main

import (
    "github.com/ciathefed/fehler-go"
)

func main() {
    reporter := fehler.NewErrorReporter()

    source := `package main

func main() {
    x := 42
    y := x + "hello" // Error here!
    println(y)
}`

    reporter.AddSource("example.go", source)

    diag := fehler.NewDiagnostic(fehler.SeverityError, "type mismatch: cannot add int and string").
        WithRange(fehler.NewSourceRangeSpan("example.go", 5, 14, 5, 20)).
        WithHelp("consider converting the string using strconv.Itoa or similar").
        WithCode("E0001").
        WithUrl("https://docs.example.org/errors/E0001")

    reporter.Report(diag)
}
```

## Output Example

Default (Fehler) style:

```
error[E0001]: type mismatch: cannot add int and string
  example.go:5:12
  3 | func main() {
  4 |     x := 42
  5 |     y := x + "hello" // Error here!
                   ~~~~~~~
  6 |     println(y)
  7 | }
  help: consider converting the string using strconv.Itoa or similar
  see: https://docs.example.org/errors/E0001
```

## API Overview

### Severity

```go
type Severity int

const (
    SeverityFatal Severity = iota
    SeverityError
    SeverityWarning
    SeverityNote
    SeverityTodo
    SeverityUnimplemented
)
```

Use `.Label()` and `.Color()` methods to access readable labels or ANSI color codes.

### SourceRange

```go
func NewSourceRangeSingle(file string, line, column int) SourceRange
func NewSourceRangeSpan(file string, startLine, startColumn, endLine, endColumn int) SourceRange
```

Use `.IsSingleChar()`, `.IsMultiline()`, and `.Length()` methods to inspect the range.

### Diagnostic

```go
type Diagnostic struct {
    Severity Severity
    Message  string
    Range    *SourceRange
    Help     *string
    Code     *string
    Url      *string
}
```

Builder methods:

```go
func NewDiagnostic(severity Severity, message string) *Diagnostic
func (d *Diagnostic) WithRange(r SourceRange) *Diagnostic
func (d *Diagnostic) WithLocation(file string, line, column int) *Diagnostic
func (d *Diagnostic) WithHelp(help string) *Diagnostic
func (d *Diagnostic) WithCode(code string) *Diagnostic
func (d *Diagnostic) WithUrl(url string) *Diagnostic
```

Convenience:

```go
func NewDiagnosticWithLocation(...) *Diagnostic
func NewDiagnosticWithRange(...) *Diagnostic
```

### ErrorReporter

```go
type ErrorReporter struct

func NewErrorReporter() *ErrorReporter
func (e *ErrorReporter) WithFormat(format OutputFormat) *ErrorReporter
func (e *ErrorReporter) AddSource(filename string, content string)
func (e *ErrorReporter) Report(d *Diagnostic)
func (e *ErrorReporter) ReportMany(diagnostics []*Diagnostic)
```

### OutputFormat

```go
type OutputFormat int

const (
    FormatFehler OutputFormat = iota
    FormatGCC
    FormatMSVC
)
```

Use `WithFormat()` to switch output style.

## Format Examples

### GCC

```go
reporter := fehler.NewErrorReporter().WithFormat(fehler.FormatGCC)
reporter.Report(diag)
```

Output:

```
example.go:5:12: error: type mismatch: cannot add int and string
```

### MSVC

```go
reporter := fehler.NewErrorReporter().WithFormat(fehler.FormatMSVC)
reporter.Report(diag)
```

Output:

```
example.go(5, 12): error E0001: type mismatch: cannot add int and string
```

## Contributing

1. Fork the repo
2. Make a new branch
3. Keep changes small and focused
4. Add tests if needed
5. Open a pull request

Commit messages should follow [Conventional Commits](https://www.conventionalcommits.org/):

* `fix: handle missing source context`
* `feat: add MSVC-style output`
* `refactor: unify underline printer`

## License

This project is licensed under the [MIT License](./LICENSE)
