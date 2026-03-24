package spec

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// ParseTasksV2 reads a tasks.md file, parses its YAML frontmatter into TasksFrontmatterV2,
// and returns the remaining body string. Enables YAML round-trip task status tracking.
func ParseTasksV2(tasksPath string) (TasksFrontmatterV2, string, error) {
	f, err := os.Open(tasksPath)
	if err != nil {
		return TasksFrontmatterV2{}, "", fmt.Errorf("open tasks file: %w", err)
	}
	defer f.Close()

	var fm TasksFrontmatterV2
	rest, err := frontmatter.Parse(f, &fm)
	if err != nil {
		// No valid frontmatter — return zero-value fm and full file as body
		content, readErr := os.ReadFile(tasksPath)
		if readErr != nil {
			return TasksFrontmatterV2{}, "", readErr
		}
		return TasksFrontmatterV2{}, string(content), nil
	}

	return fm, string(rest), nil
}

// UpdateTaskStatus updates a single task's status by ID, recomputes the Completed count,
// and writes the updated frontmatter + body back to the file.
func UpdateTaskStatus(tasksPath string, taskID int, newStatus ItemStatus) error {
	fm, body, err := ParseTasksV2(tasksPath)
	if err != nil {
		return fmt.Errorf("parse tasks: %w", err)
	}

	found := false
	for i := range fm.Tasks {
		if fm.Tasks[i].ID == taskID {
			fm.Tasks[i].Status = newStatus
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task %d not found", taskID)
	}

	// Recompute Completed count: number of tasks with StatusDone
	completed := 0
	for _, t := range fm.Tasks {
		if t.Status == StatusDone {
			completed++
		}
	}
	fm.Completed = completed

	return WriteTasks(tasksPath, fm, body)
}

// WriteTasks serializes TasksFrontmatterV2 as YAML frontmatter, prepends/appends the
// `---` delimiters, and appends the original body content. Preserves markdown body unchanged.
func WriteTasks(tasksPath string, fm TasksFrontmatterV2, body string) error {
	yamlBytes, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("marshal frontmatter: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlBytes)
	sb.WriteString("---\n")
	sb.WriteString(body)

	return os.WriteFile(tasksPath, []byte(sb.String()), 0644)
}
