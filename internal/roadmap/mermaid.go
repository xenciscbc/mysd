package roadmap

import (
	"bytes"
	"text/template"
	"time"
)

// TrackingFile is the root structure of tracking.yaml.
type TrackingFile struct {
	SchemaVersion string         `yaml:"schema_version"`
	UpdatedAt     time.Time      `yaml:"updated_at"`
	Changes       []ChangeRecord `yaml:"changes"`
}

// ChangeRecord tracks the lifecycle of a single change.
type ChangeRecord struct {
	Name           string     `yaml:"name"`
	Status         string     `yaml:"status"`
	StartedAt      *time.Time `yaml:"started_at,omitempty"`
	CompletedAt    *time.Time `yaml:"completed_at,omitempty"`
	TotalTasks     int        `yaml:"total_tasks"`
	CompletedTasks int        `yaml:"completed_tasks"`
	MustTotal      int        `yaml:"must_total"`
	MustPassed     int        `yaml:"must_passed"`
}

const mermaidTmpl = `gantt
    title Roadmap
    dateFormat YYYY-MM-DD
{{- range .Changes}}
    section {{.Name}}
        {{.Name}} :{{ganttStatus .}}, {{formatStartDate .}}, {{formatEndDate .}}
{{- end}}
`

// GenerateMermaid produces a Mermaid gantt chart string from a TrackingFile.
// If there are no changes, it returns a valid gantt header with no sections.
func GenerateMermaid(tf TrackingFile) string {
	funcMap := template.FuncMap{
		"ganttStatus": func(cr ChangeRecord) string {
			switch cr.Status {
			case "archived", "verified":
				return "done"
			case "executed":
				return "active"
			default:
				return ""
			}
		},
		"formatStartDate": func(cr ChangeRecord) string {
			if cr.StartedAt != nil {
				return cr.StartedAt.Format("2006-01-02")
			}
			return tf.UpdatedAt.Format("2006-01-02")
		},
		"formatEndDate": func(cr ChangeRecord) string {
			if cr.CompletedAt != nil {
				return cr.CompletedAt.Format("2006-01-02")
			}
			// Use today + 1 day as a future placeholder for in-progress items
			return time.Now().Add(24 * time.Hour).Format("2006-01-02")
		},
	}

	tmpl, err := template.New("mermaid").Funcs(funcMap).Parse(mermaidTmpl)
	if err != nil {
		// Fallback: return minimal valid gantt
		return "gantt\n    title Roadmap\n    dateFormat YYYY-MM-DD\n"
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tf); err != nil {
		return "gantt\n    title Roadmap\n    dateFormat YYYY-MM-DD\n"
	}

	return buf.String()
}
