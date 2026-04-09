// Package verifier provides the verification engine for mysd.
// It builds verification context from spec requirements, parses verifier agent output,
// and writes gap and verification reports.
package verifier

import (
	"fmt"
	"hash/crc32"
	"path/filepath"
	"strings"

	"github.com/xenciscbc/mysd/internal/spec"
)

// VerificationContext is the JSON-serializable context passed to the verifier agent.
// It contains all MUST/SHOULD/MAY items from the spec and a summary of tasks.
type VerificationContext struct {
	SpecDir      string       `json:"spec_dir"`
	ChangeName   string       `json:"change_name"`
	ChangeDir    string       `json:"change_dir"`
	SpecsDir     string       `json:"specs_dir"`
	MustItems    []VerifyItem `json:"must_items"`
	ShouldItems  []VerifyItem `json:"should_items"`
	MayItems     []VerifyItem `json:"may_items"`
	TasksSummary []TaskItem   `json:"tasks_summary"`
}

// VerifyItem represents a single requirement in the verification context.
type VerifyItem struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
}

// TaskItem represents a task entry in the verification context summary.
type TaskItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// StableID generates a stable, deterministic ID for a requirement.
// Format: "{source_file}::{keyword_lower}-{hex_crc32}"
// Uses CRC32 IEEE checksum of the requirement text — not sequential counter.
// This ID is stable across re-parses as long as the text does not change.
func StableID(r spec.Requirement) string {
	hash := crc32.ChecksumIEEE([]byte(r.Text))
	keyword := strings.ToLower(string(r.Keyword))
	return fmt.Sprintf("%s::%s-%x", r.SourceFile, keyword, hash)
}

// BuildVerificationContextFromParts constructs a VerificationContext from pre-loaded data.
// This is the pure-function constructor for testing (no filesystem I/O).
func BuildVerificationContextFromParts(
	changeName, changeDir, specsDir string,
	reqs []spec.Requirement,
	tasks []spec.Task,
) VerificationContext {
	ctx := VerificationContext{
		ChangeName:   changeName,
		ChangeDir:    changeDir,
		SpecsDir:     specsDir,
		MustItems:    []VerifyItem{},
		ShouldItems:  []VerifyItem{},
		MayItems:     []VerifyItem{},
		TasksSummary: []TaskItem{},
	}

	// Classify requirements by RFC 2119 keyword
	for _, r := range reqs {
		item := VerifyItem{
			ID:      StableID(r),
			Text:    r.Text,
			Keyword: string(r.Keyword),
			Status:  string(r.Status),
		}
		switch r.Keyword {
		case spec.Must:
			ctx.MustItems = append(ctx.MustItems, item)
		case spec.Should:
			ctx.ShouldItems = append(ctx.ShouldItems, item)
		case spec.May:
			ctx.MayItems = append(ctx.MayItems, item)
		}
	}

	// Populate task summary
	for _, t := range tasks {
		ctx.TasksSummary = append(ctx.TasksSummary, TaskItem{
			ID:     t.ID,
			Name:   t.Name,
			Status: string(t.Status),
		})
	}

	return ctx
}

// BuildVerificationContext loads a change from disk and constructs a VerificationContext.
// specsDir is the root specs directory (e.g., ".specs") and changeName identifies the change.
// This is the entrypoint for `mysd verify --context-only`.
func BuildVerificationContext(specsDir, changeName string) (VerificationContext, error) {
	changeDir := filepath.Join(specsDir, "changes", changeName)

	change, err := spec.ParseChange(changeDir)
	if err != nil {
		return VerificationContext{}, fmt.Errorf("parse change: %w", err)
	}

	return BuildVerificationContextFromParts(changeName, changeDir, specsDir, change.Specs, change.Tasks), nil
}
