package spec

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// Scaffold creates a complete change directory structure under baseDir/changes/{name}/.
// It writes .openspec.yaml, proposal.md, specs/, design.md, and tasks.md with
// appropriate frontmatter and template bodies.
func Scaffold(name string, baseDir string) (Change, error) {
	today := time.Now().Format("2006-01-02")
	changeDir := filepath.Join(baseDir, "changes", name)

	if err := os.MkdirAll(changeDir, 0755); err != nil {
		return Change{}, fmt.Errorf("create change dir: %w", err)
	}

	// Create specs/ directory (empty; specs added later by `mysd spec`)
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return Change{}, fmt.Errorf("create specs dir: %w", err)
	}

	// Write .openspec.yaml
	metaContent := fmt.Sprintf("schema: spec-driven\ncreated: %s\n", today)
	if err := os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte(metaContent), 0644); err != nil {
		return Change{}, fmt.Errorf("write .openspec.yaml: %w", err)
	}

	// Write proposal.md
	proposalContent, err := renderTemplate(proposalTemplate, map[string]string{
		"Name":    name,
		"Today":   today,
		"Version": "1",
	})
	if err != nil {
		return Change{}, fmt.Errorf("render proposal template: %w", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		return Change{}, fmt.Errorf("write proposal.md: %w", err)
	}

	// Write design.md
	designContent, err := renderTemplate(designTemplate, nil)
	if err != nil {
		return Change{}, fmt.Errorf("render design template: %w", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "design.md"), []byte(designContent), 0644); err != nil {
		return Change{}, fmt.Errorf("write design.md: %w", err)
	}

	// Write tasks.md
	tasksContent, err := renderTemplate(tasksTemplate, map[string]string{
		"Today":   today,
		"Version": "1",
	})
	if err != nil {
		return Change{}, fmt.Errorf("render tasks template: %w", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		return Change{}, fmt.Errorf("write tasks.md: %w", err)
	}

	return Change{
		Name: name,
		Dir:  changeDir,
		Meta: ChangeMeta{
			Schema:  "spec-driven",
			Created: today,
		},
	}, nil
}

// renderTemplate executes a text/template string with the given data map.
func renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// proposalTemplate is the default proposal.md template.
const proposalTemplate = `---
spec-version: "{{.Version}}"
change: {{.Name}}
status: proposed
created: {{.Today}}
updated: {{.Today}}
---

## Summary

_Describe what this change does and why._

## Motivation

_Explain the problem or opportunity this change addresses._

## Scope

_List what is included and explicitly excluded from this change._
`

// designTemplate is the default design.md template.
const designTemplate = `## Architecture

_Describe the high-level design approach._

## Key Decisions

_Document significant design choices and their rationale._
`

// tasksTemplate is the default tasks.md template.
const tasksTemplate = `---
spec-version: "{{.Version}}"
total: 0
completed: 0
---

## Tasks

_Add tasks here as implementation is planned._
`
