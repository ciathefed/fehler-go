package fehler

import (
	"encoding/json"
	"io"
)

type SarifReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SarifRun `json:"runs"`
}

type SarifRun struct {
	Tool    SarifTool     `json:"tool"`
	Results []SarifResult `json:"results"`
}

type SarifTool struct {
	Driver SarifDriver `json:"driver"`
}

type SarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SarifRule `json:"rules,omitempty"`
}

type SarifRule struct {
	ID                   string              `json:"id"`
	ShortDescription     SarifMessage        `json:"shortDescription"`
	DefaultConfiguration *SarifConfiguration `json:"defaultConfiguration,omitempty"`
	HelpURI              string              `json:"helpUri,omitempty"`
}

type SarifConfiguration struct {
	Level string `json:"level"`
}

type SarifResult struct {
	Message   SarifMessage    `json:"message"`
	Level     string          `json:"level"`
	RuleID    *string         `json:"ruleId,omitempty"`
	Locations []SarifLocation `json:"locations,omitempty"`
	Kind      string          `json:"kind,omitempty"`
}

type SarifMessage struct {
	Text string `json:"text"`
}

type SarifLocation struct {
	PhysicalLocation SarifPhysicalLocation `json:"physicalLocation"`
}

type SarifPhysicalLocation struct {
	ArtifactLocation SarifArtifactLocation `json:"artifactLocation"`
	Region           SarifRegion           `json:"region"`
}

type SarifArtifactLocation struct {
	URI string `json:"uri"`
}

type SarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndLine     int `json:"endLine"`
	EndColumn   int `json:"endColumn"`
}

func sarifLevel(sev Severity) string {
	switch sev {
	case SeverityFatal, SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityNote:
		return "note"
	case SeverityTodo, SeverityUnimplemented:
		return "none"
	default:
		return "none"
	}
}

// Emits all diagnostics in SARIF format to the given writer.
// Supports version 2.1.0. Includes rule metadata if code is set.
func EmitSarif(diagnostics []*Diagnostic, w io.Writer) error {
	const sarifVersion = "2.1.0"
	const sarifSchema = "https://json.schemastore.org/sarif-2.1.0.json"

	ruleMap := make(map[string]SarifRule)
	for _, d := range diagnostics {
		if d.Code != nil {
			code := *d.Code
			if _, exists := ruleMap[code]; !exists {
				ruleMap[code] = SarifRule{
					ID: code,
					ShortDescription: SarifMessage{
						Text: d.Message,
					},
					DefaultConfiguration: &SarifConfiguration{
						Level: sarifLevel(d.Severity),
					},
					HelpURI: func() string {
						if d.Url != nil {
							return *d.Url
						}
						return ""
					}(),
				}
			}
		}
	}

	rules := make([]SarifRule, 0, len(ruleMap))
	for _, r := range ruleMap {
		rules = append(rules, r)
	}

	results := make([]SarifResult, 0, len(diagnostics))
	for _, d := range diagnostics {
		res := SarifResult{
			Message: SarifMessage{
				Text: d.Message,
			},
			Level: sarifLevel(d.Severity),
			Kind:  "fail",
		}
		if d.Code != nil {
			res.RuleID = d.Code
		}
		if d.Range != nil {
			loc := SarifLocation{
				PhysicalLocation: SarifPhysicalLocation{
					ArtifactLocation: SarifArtifactLocation{
						URI: d.Range.File,
					},
					Region: SarifRegion{
						StartLine:   d.Range.Start.Line,
						StartColumn: d.Range.Start.Column,
						EndLine:     d.Range.End.Line,
						EndColumn:   d.Range.End.Column,
					},
				},
			}
			res.Locations = []SarifLocation{loc}
		}
		results = append(results, res)
	}

	report := SarifReport{
		Version: sarifVersion,
		Schema:  sarifSchema,
		Runs: []SarifRun{{
			Tool: SarifTool{
				Driver: SarifDriver{
					Name:           "fehler",
					Version:        "0.5.0",
					InformationURI: "https://github.com/ciathefed/fehler",
					Rules:          rules,
				},
			},
			Results: results,
		}},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(report)
}
