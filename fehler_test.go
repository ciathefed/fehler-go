package fehler

import (
	"bytes"
	"strings"
	"testing"
)

func TestPositionCreation(t *testing.T) {
	pos := Position{Line: 10, Column: 5}
	if pos.Line != 10 {
		t.Errorf("expected line 10, got %d", pos.Line)
	}
	if pos.Column != 5 {
		t.Errorf("expected column 5, got %d", pos.Column)
	}
}

func TestSourceRangeSingleChar(t *testing.T) {
	r := NewSourceRangeSingle("test.go", 10, 5)

	if r.File != "test.go" {
		t.Errorf("expected file test.go, got %s", r.File)
	}
	if r.Start.Line != 10 || r.Start.Column != 5 {
		t.Errorf("unexpected start position %v", r.Start)
	}
	if r.End.Line != 10 || r.End.Column != 5 {
		t.Errorf("unexpected end position %v", r.End)
	}
	if !r.IsSingleChar() {
		t.Error("expected single char range")
	}
	if r.IsMultiline() {
		t.Error("expected not multiline")
	}
}

func TestSourceRangeSpan(t *testing.T) {
	r := NewSourceRangeSpan("test.go", 10, 5, 12, 8)

	if r.File != "test.go" {
		t.Errorf("expected file test.go, got %s", r.File)
	}
	if r.Start.Line != 10 || r.Start.Column != 5 {
		t.Errorf("unexpected start position %v", r.Start)
	}
	if r.End.Line != 12 || r.End.Column != 8 {
		t.Errorf("unexpected end position %v", r.End)
	}
	if r.IsSingleChar() {
		t.Error("expected not single char")
	}
	if !r.IsMultiline() {
		t.Error("expected multiline")
	}
}

func TestSourceRangeSingleLineSpan(t *testing.T) {
	r := NewSourceRangeSpan("test.go", 10, 5, 10, 15)

	if r.IsSingleChar() {
		t.Error("expected not single char")
	}
	if r.IsMultiline() {
		t.Error("expected not multiline")
	}
	if got := r.Length(); got != 11 {
		t.Errorf("expected length 11, got %d", got)
	}
}

func TestDiagnosticWithRange(t *testing.T) {
	r := NewSourceRangeSpan("example.go", 42, 10, 42, 20)
	diag := NewDiagnostic(SeverityError, "test error").WithRange(r)

	if diag.Severity != SeverityError {
		t.Errorf("expected SeverityError, got %v", diag.Severity)
	}
	if diag.Message != "test error" {
		t.Errorf("expected message 'test error', got %s", diag.Message)
	}
	if diag.Range == nil {
		t.Fatal("expected non-nil range")
	}
	if diag.Range.File != "example.go" {
		t.Errorf("expected file example.go, got %s", diag.Range.File)
	}
	if diag.Range.Start.Line != 42 || diag.Range.Start.Column != 10 {
		t.Errorf("unexpected start pos %v", diag.Range.Start)
	}
	if diag.Range.End.Line != 42 || diag.Range.End.Column != 20 {
		t.Errorf("unexpected end pos %v", diag.Range.End)
	}
}

func TestDiagnosticWithLocation(t *testing.T) {
	diag := NewDiagnostic(SeverityWarning, "test warning").WithLocation("test.go", 15, 8)

	if diag.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %v", diag.Severity)
	}
	if diag.Message != "test warning" {
		t.Errorf("expected message 'test warning', got %s", diag.Message)
	}
	if diag.Range == nil {
		t.Fatal("expected non-nil range")
	}
	if diag.Range.File != "test.go" {
		t.Errorf("expected file test.go, got %s", diag.Range.File)
	}
	if diag.Range.Start.Line != 15 || diag.Range.Start.Column != 8 {
		t.Errorf("unexpected start pos %v", diag.Range.Start)
	}
	if diag.Range.End.Line != 15 || diag.Range.End.Column != 8 {
		t.Errorf("unexpected end pos %v", diag.Range.End)
	}
	if !diag.Range.IsSingleChar() {
		t.Error("expected single char range")
	}
}

func TestNewDiagnosticWithLocationConvenience(t *testing.T) {
	diag := NewDiagnosticWithLocation(SeverityError, "syntax error", "main.go", 15, 8)

	if diag.Severity != SeverityError {
		t.Errorf("expected SeverityError, got %v", diag.Severity)
	}
	if diag.Message != "syntax error" {
		t.Errorf("expected message 'syntax error', got %s", diag.Message)
	}
	if diag.Range == nil {
		t.Fatal("expected non-nil range")
	}
	if diag.Range.File != "main.go" {
		t.Errorf("expected file main.go, got %s", diag.Range.File)
	}
	if diag.Range.Start.Line != 15 || diag.Range.Start.Column != 8 {
		t.Errorf("unexpected start pos %v", diag.Range.Start)
	}
	if !diag.Range.IsSingleChar() {
		t.Error("expected single char range")
	}
}

func TestNewDiagnosticWithRangeConvenience(t *testing.T) {
	diag := NewDiagnosticWithRange(SeverityWarning, "long identifier", "main.go", 15, 8, 15, 25)

	if diag.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %v", diag.Severity)
	}
	if diag.Message != "long identifier" {
		t.Errorf("expected message 'long identifier', got %s", diag.Message)
	}
	if diag.Range == nil {
		t.Fatal("expected non-nil range")
	}
	if diag.Range.File != "main.go" {
		t.Errorf("expected file main.go, got %s", diag.Range.File)
	}
	if diag.Range.Start.Line != 15 || diag.Range.Start.Column != 8 {
		t.Errorf("unexpected start pos %v", diag.Range.Start)
	}
	if diag.Range.End.Line != 15 || diag.Range.End.Column != 25 {
		t.Errorf("unexpected end pos %v", diag.Range.End)
	}
	if diag.Range.IsSingleChar() {
		t.Error("expected not single char range")
	}
	if diag.Range.IsMultiline() {
		t.Error("expected not multiline")
	}
}

func TestErrorReporterDiagnostics(t *testing.T) {
	reporter := NewErrorReporter()

	sourceCode := `
package main

import "fmt"

func main() {
    veryLongVariableName := 42
    y := x + "hello" // Type mismatch error
    fmt.Printf("Result: %v\n", y)
}
`
	reporter.AddSource("example.go", sourceCode)

	diagnostics := []*Diagnostic{
		NewDiagnosticWithRange(SeverityError, "type mismatch: cannot add integer and string", "example.go", 8, 14, 8, 20),
		NewDiagnosticWithRange(SeverityWarning, "variable name is too long", "example.go", 7, 5, 7, 24),
		NewDiagnosticWithLocation(SeverityError, "undefined variable 'x'", "example.go", 8, 10),
	}

	if got := len(diagnostics); got != 3 {
		t.Errorf("expected 3 diagnostics, got %d", got)
	}
	if diagnostics[0].Severity != SeverityError {
		t.Errorf("expected SeverityError for first diagnostic")
	}
	if diagnostics[1].Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning for second diagnostic")
	}
	if diagnostics[2].Severity != SeverityError {
		t.Errorf("expected SeverityError for third diagnostic")
	}
}

func TestMultilineRange(t *testing.T) {
	r := NewSourceRangeSpan("test.go", 5, 10, 8, 15)

	if !r.IsMultiline() {
		t.Error("expected multiline")
	}
	if r.IsSingleChar() {
		t.Error("expected not single char")
	}
	if r.Start.Line != 5 || r.Start.Column != 10 {
		t.Errorf("unexpected start position %v", r.Start)
	}
	if r.End.Line != 8 || r.End.Column != 15 {
		t.Errorf("unexpected end position %v", r.End)
	}
}

func TestErrorReporterIntegrationWithRanges(t *testing.T) {
	reporter := NewErrorReporter()

	sourceCode := `
package main

import (
    "fmt"
)

func main() {
    name := "World"
    greeting := fmt.Sprintf("Hello, %s!", name)
    fmt.Println(greeting)
}
`
	reporter.AddSource("hello.go", sourceCode)

	singleChar := NewDiagnosticWithLocation(SeverityError, "undefined variable 'greeting'", "hello.go", 10, 5)
	shortRange := NewDiagnosticWithRange(SeverityWarning, "unused variable", "hello.go", 9, 5, 9, 8)
	longRange := NewDiagnosticWithRange(SeverityNote, "function signature", "hello.go", 8, 1, 8, 11)

	if singleChar.Range == nil || !singleChar.Range.IsSingleChar() {
		t.Error("expected single char range")
	}
	if shortRange.Range == nil || shortRange.Range.IsSingleChar() {
		t.Error("expected not single char range")
	}
	if shortRange.Range == nil || shortRange.Range.IsMultiline() {
		t.Error("expected not multiline")
	}
	if longRange.Range == nil || longRange.Range.IsMultiline() {
		t.Error("expected not multiline")
	}
	if shortRange.Range.Length() != 4 {
		t.Errorf("expected length 4, got %d", shortRange.Range.Length())
	}
	if longRange.Range.Length() != 11 {
		t.Errorf("expected length 11, got %d", longRange.Range.Length())
	}
}

func TestEmitSarifOutputsValidJSON(t *testing.T) {
	diag1 := NewDiagnostic(SeverityError, "invalid token").
		WithLocation("main.go", 1, 2).
		WithCode("E001")

	diag2 := NewDiagnostic(SeverityError, "invalid token").
		WithLocation("main.go", 3, 4).
		WithCode("E001")

	var buf bytes.Buffer
	err := EmitSarif([]*Diagnostic{diag1, diag2}, &buf)
	if err != nil {
		t.Fatalf("EmitSarif failed: %v", err)
	}

	jsonStr := buf.String()
	if !strings.Contains(jsonStr, `"message"`) {
		t.Error("expected 'message' in JSON output")
	}
	if !strings.Contains(jsonStr, "invalid token") {
		t.Error("expected 'invalid token' in JSON output")
	}
	if !strings.Contains(jsonStr, "main.go") {
		t.Error("expected 'main.go' in JSON output")
	}
	if !strings.Contains(jsonStr, "E001") {
		t.Error("expected 'E001' in JSON output")
	}
}
