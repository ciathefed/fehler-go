package fehler

import (
	"fmt"
	"strings"
)

const (
	colorReset   = "\x1b[0m"
	colorRed     = "\x1b[31m"
	colorYellow  = "\x1b[33m"
	colorBlue    = "\x1b[34m"
	colorMagenta = "\x1b[35m"
	colorCyan    = "\x1b[36m"
	colorWhite   = "\x1b[37m"
	colorBold    = "\x1b[1m"
	colorDim     = "\x1b[2m"
)

type OutputFormat int

const (
	FormatFehler OutputFormat = iota
	FormatGCC
	FormatMSVC
)

// Represents a position in source code with line and column information.
type Position struct {
	Line   int
	Column int
}

// Represents a range in source code with start and end positions.
type SourceRange struct {
	File  string
	Start Position
	End   Position
}

// Creates a single-character range at the specified position.
func NewSourceRangeSingle(file string, line int, column int) SourceRange {
	return SourceRange{
		File:  file,
		Start: Position{Line: line, Column: column},
		End:   Position{Line: line, Column: column},
	}
}

// Creates a range spanning from start to end positions.
func NewSourceRangeSpan(file string, startLine int, startColumn int, endLine int, endColumn int) SourceRange {
	return SourceRange{
		File:  file,
		Start: Position{Line: startLine, Column: startColumn},
		End:   Position{Line: endLine, Column: endColumn},
	}
}

// Returns true if this range spans multiple lines.
func (s SourceRange) IsMultiline() bool {
	return s.Start.Line != s.End.Line
}

// Returns true if this range is a single character.
func (s SourceRange) IsSingleChar() bool {
	return s.Start.Line == s.End.Line && s.Start.Column == s.End.Column
}

// Returns the length of the range on a single line (only valid for single-line ranges).
func (s SourceRange) Length() int {
	if s.IsMultiline() {
		return 0
	}
	if s.End.Column >= s.Start.Column {
		return s.End.Column - s.Start.Column + 1
	}
	return 1
}

// Severity levels for diagnostics, determining color and label presentation.
type Severity int

const (
	SeverityFatal Severity = iota
	SeverityError
	SeverityWarning
	SeverityNote
	SeverityTodo
	SeverityUnimplemented
)

// Returns the ANSI color code associated with this severity level.
func (s Severity) Color() string {
	switch s {
	case SeverityFatal, SeverityError:
		return colorRed
	case SeverityWarning:
		return colorYellow
	case SeverityNote:
		return colorBlue
	case SeverityTodo:
		return colorMagenta
	case SeverityUnimplemented:
		return colorCyan
	default:
		return ""
	}
}

// Returns the human-readable label for this severity level.
func (s Severity) Label() string {
	switch s {
	case SeverityFatal:
		return "fatal"
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityNote:
		return "note"
	case SeverityTodo:
		return "todo"
	case SeverityUnimplemented:
		return "unimplemented"
	default:
		return "unknown"
	}
}

// A diagnostic message with optional source range and help text.
// This is the primary data structure for representing compiler errors, warnings, and notes.
type Diagnostic struct {
	Severity Severity
	Message  string
	Range    *SourceRange
	Help     *string
	Code     *string
	Url      *string
}

// Creates a new diagnostic with the specified severity and message.
// Additional properties can be added using the fluent interface methods.
func NewDiagnostic(severity Severity, message string) *Diagnostic {
	return &Diagnostic{
		Severity: severity,
		Message:  message,
	}
}

// Returns a copy of this diagnostic with the specified source range.
// This method follows the builder pattern for fluent construction of diagnostics.
func (d *Diagnostic) WithRange(r SourceRange) *Diagnostic {
	d.Range = &r
	return d
}

// Returns a copy of this diagnostic with a single-character range.
// This method follows the builder pattern for fluent construction of diagnostics.
func (d *Diagnostic) WithLocation(file string, line int, column int) *Diagnostic {
	r := NewSourceRangeSingle(file, line, column)
	d.Range = &r
	return d
}

// Returns a copy of this diagnostic with the specified help text.
// This method follows the builder pattern for fluent construction of diagnostics.
func (d *Diagnostic) WithHelp(help string) *Diagnostic {
	d.Help = &help
	return d
}

// Returns a copy of this diagnostic with the specified error code.
// The code can be used to look up error documentation.
func (d *Diagnostic) WithCode(code string) *Diagnostic {
	d.Code = &code
	return d
}

// Returns a copy of this diagnostic with the specified documentation URL.
// Useful for linking to online resources about this error.
func (d *Diagnostic) WithUrl(url string) *Diagnostic {
	d.Url = &url
	return d
}

// A comprehensive error reporting system that manages source files and formats diagnostics.
// This reporter can store multiple source files and display rich error messages with
// source code context, similar to modern compiler error output.
type ErrorReporter struct {
	Sources map[string]string
	Format  OutputFormat
}

// Initializes a new ErrorReporter with the given allocator.
// The reporter starts with no source files registered.
// Uses the default output format (Fehler).
func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{
		Sources: make(map[string]string),
		Format:  FormatFehler,
	}
}

// Returns a copy of this reporter with the specified output format.
func (e *ErrorReporter) WithFormat(format OutputFormat) *ErrorReporter {
	e.Format = format
	return e
}

// Adds a source file to the reporter for later reference in diagnostics.
// The content is duplicated and owned by the reporter.
func (e *ErrorReporter) AddSource(filename string, content string) {
	e.Sources[filename] = content
}

// Reports a single diagnostic to stdout with color formatting.
// If the diagnostic has a range and the source file is available,
// displays a source code snippet with the error range highlighted.
func (e *ErrorReporter) Report(diagnostic *Diagnostic) {
	switch e.Format {
	case FormatFehler:
		e.printFehler(diagnostic)
	case FormatGCC:
		e.printGcc(diagnostic)
	case FormatMSVC:
		e.printMsvc(diagnostic)
	}
}

// Reports multiple diagnostics in sequence.
// Each diagnostic is printed with the same formatting as `report()`.
func (e *ErrorReporter) ReportMany(diagnostics []*Diagnostic) {
	for _, diagnostic := range diagnostics {
		e.Report(diagnostic)
	}
}

func (e *ErrorReporter) printFehler(diagnostic *Diagnostic) {
	if diagnostic.Code != nil {
		fmt.Printf("%s%s%s[%s]%s: %s\n",
			diagnostic.Severity.Color(),
			colorBold,
			diagnostic.Severity.Label(),
			*diagnostic.Code,
			colorReset,
			diagnostic.Message,
		)
	} else {
		fmt.Printf("%s%s%s%s: %s\n",
			diagnostic.Severity.Color(),
			colorBold,
			diagnostic.Severity.Label(),
			colorReset,
			diagnostic.Message,
		)
	}

	if diagnostic.Range != nil {
		r := *diagnostic.Range
		fmt.Printf("  %s%s%s:%d:%d%s\n",
			colorCyan,
			colorBold,
			r.File,
			r.Start.Line,
			r.Start.Column,
			colorReset,
		)

		color := diagnostic.Severity.Color()
		e.printSourceSnippet(r, color)
	}

	if diagnostic.Help != nil {
		fmt.Printf("  %s%shelp%s: %s\n", colorCyan, colorBold, colorReset, *diagnostic.Help)
	}

	if diagnostic.Url != nil {
		fmt.Printf("  %s%ssee%s: %s\n", colorCyan, colorBold, colorReset, *diagnostic.Url)
	}

	fmt.Println()
}

func (e *ErrorReporter) printGcc(diagnostic *Diagnostic) {
	color := diagnostic.Severity.Color()
	if diagnostic.Range != nil {
		r := *diagnostic.Range
		fmt.Printf("%s%s:%d:%d: %s%s: %s%s%s%s\n",
			colorBold,
			r.File,
			r.Start.Line,
			r.Start.Column,
			color,
			diagnostic.Severity.Label(),
			colorReset,
			colorBold,
			diagnostic.Message,
			colorReset,
		)
	} else {
		fmt.Printf("%s%s%s: %s%s%s%s\n",
			colorBold,
			color,
			diagnostic.Severity.Label(),
			colorReset,
			colorBold,
			diagnostic.Message,
			colorReset,
		)
	}
}

func (e *ErrorReporter) printMsvc(diagnostic *Diagnostic) {
	if diagnostic.Range != nil {
		code := "unknown"
		if diagnostic.Code != nil {
			code = *diagnostic.Code
		}
		r := *diagnostic.Range
		fmt.Printf("%s(%d, %d): %s %s: %s\n",
			r.File,
			r.Start.Line,
			r.Start.Column,
			diagnostic.Severity.Label(),
			code,
			diagnostic.Message,
		)
	} else {
		fmt.Printf("%s: %s\n",
			diagnostic.Severity.Label(),
			diagnostic.Message,
		)
	}
}

// Prints a source code snippet showing the context around a diagnostic range.
// Shows 2 lines before and after the error location, with the error range highlighted
// using carets (^) for single characters or tildes (~) for ranges.
func (e *ErrorReporter) printSourceSnippet(r SourceRange, color string) {
	source, ok := e.Sources[r.File]
	if !ok {
		return
	}

	lines := strings.Split(source, "\n")
	contextStart := 1
	if r.Start.Line > 2 {
		contextStart = r.Start.Line - 2
	}

	contextEnd := r.Start.Line + 2
	if r.IsMultiline() {
		contextEnd = r.End.Line + 2
	}
	if contextEnd > len(lines) {
		contextEnd = len(lines)
	}

	for currentLine := contextStart; currentLine <= contextEnd; currentLine++ {
		line := lines[currentLine-1]
		lineNumWidth := 4
		isErrorLine := currentLine >= r.Start.Line && currentLine <= r.End.Line

		if isErrorLine {
			fmt.Printf("  %s%s%4d |%s %s\n",
				colorRed,
				colorBold,
				currentLine,
				colorReset,
				line,
			)

			e.printUnderline(r, currentLine, lineNumWidth, color)
		} else {
			fmt.Printf("  %s%4d |%s %s\n",
				colorDim,
				currentLine,
				colorReset,
				line,
			)
		}
	}
}

// Prints the underline (carets or tildes) for a specific line in a range.
func (e *ErrorReporter) printUnderline(r SourceRange, lineNum int, lineNumWidth int, color string) {
	fmt.Print("  ", color)
	fmt.Print(strings.Repeat(" ", lineNumWidth+1))
	fmt.Print("  ")

	if r.IsMultiline() {
		if lineNum == r.Start.Line {
			fmt.Print(strings.Repeat(" ", r.Start.Column-1))
			fmt.Print("~")
			fmt.Print(strings.Repeat("~", 80-(r.Start.Column)))
		} else if lineNum == r.End.Line {
			fmt.Print(strings.Repeat("~", r.End.Column))
		} else if lineNum > r.Start.Line && lineNum < r.End.Line {
			fmt.Print(strings.Repeat("~", 80))
		}
	} else {
		fmt.Print(strings.Repeat(" ", r.Start.Column-1))
		if r.IsSingleChar() {
			fmt.Print("^")
		} else {
			fmt.Print(strings.Repeat("~", r.Length()))
		}
	}

	fmt.Println(colorReset)
}

// Convenience function to create a diagnostic with single-character location information.
func NewDiagnosticWithLocation(severity Severity, message, file string, line, column int) *Diagnostic {
	return NewDiagnostic(severity, message).WithLocation(file, line, column)
}

// Convenience function to create a diagnostic with range information.
func NewDiagnosticWithRange(severity Severity, message, file string, startLine, startColumn, endLine, endColumn int) *Diagnostic {
	return NewDiagnostic(severity, message).WithRange(NewSourceRangeSpan(file, startLine, startColumn, endLine, endColumn))
}
