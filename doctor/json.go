// Package doctor — JSON emitter for `movie doctor --json`.
// Output schema is stable for scripting/CI consumption.
package doctor

import (
	"encoding/json"
	"fmt"
)

// JsonReport is the stable wire-format for `movie doctor --json`.
// Field names are snake_case for cross-language consumers.
// Field order optimized for govet fieldalignment (strings, slice, bools last).
type JsonReport struct {
	Schema    string        `json:"schema"`
	Source    string        `json:"deploy_source"`
	Target    string        `json:"active_binary"`
	DeployDir string        `json:"deploy_dir"`
	Findings  []JsonFinding `json:"findings"`
	HasErr    bool          `json:"has_errors"`
	HasFix    bool          `json:"has_fixable"`
}

// JsonFinding mirrors Finding with snake_case JSON tags.
type JsonFinding struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Severity  string `json:"severity"`
	Detail    string `json:"detail"`
	FixHint   string `json:"fix_hint"`
	IsFixable bool   `json:"is_fixable"`
}

// PrintJson writes the report as indented JSON to stdout.
func (r *Report) PrintJson() error {
	payload := r.toJson()
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func (r *Report) toJson() JsonReport {
	return JsonReport{
		Schema:    "movie-doctor/v1",
		Source:    r.Source,
		Target:    r.Target,
		DeployDir: r.DeployDir,
		HasErr:    r.HasErrors(),
		HasFix:    r.HasFixable(),
		Findings:  toJsonFindings(r.Findings),
	}
}

func toJsonFindings(findings []Finding) []JsonFinding {
	out := make([]JsonFinding, 0, len(findings))
	for _, f := range findings {
		out = append(out, JsonFinding{
			Id:        f.ID,
			Title:     f.Title,
			Severity:  string(f.Severity),
			Detail:    f.Detail,
			FixHint:   f.FixHint,
			IsFixable: f.IsFixable,
		})
	}
	return out
}
