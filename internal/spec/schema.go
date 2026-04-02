package spec

import "errors"

// RFC2119Keyword represents RFC 2119 semantic keywords for requirements.
type RFC2119Keyword string

const (
	Must   RFC2119Keyword = "MUST"
	Should RFC2119Keyword = "SHOULD"
	May    RFC2119Keyword = "MAY"
)

// DeltaOp represents the type of change operation in a delta spec.
type DeltaOp string

const (
	DeltaAdded    DeltaOp = "ADDED"
	DeltaModified DeltaOp = "MODIFIED"
	DeltaRemoved  DeltaOp = "REMOVED"
	DeltaRenamed  DeltaOp = "RENAMED"
	DeltaNone     DeltaOp = ""
)

// ItemStatus represents the lifecycle status of a spec item.
type ItemStatus string

const (
	StatusPending    ItemStatus = "pending"
	StatusInProgress ItemStatus = "in_progress"
	StatusDone       ItemStatus = "done"
	StatusBlocked    ItemStatus = "blocked"
)

// SpecDirFlavor identifies which spec directory convention a project uses.
type SpecDirFlavor string

const (
	FlavorMySD     SpecDirFlavor = "mysd"
	FlavorOpenSpec SpecDirFlavor = "openspec"
	FlavorNone     SpecDirFlavor = "none"
)

// ChangeMeta holds metadata from .openspec.yaml at the change level.
type ChangeMeta struct {
	Schema  string `yaml:"schema"`
	Created string `yaml:"created"`
}

// ProposalFrontmatter holds the YAML frontmatter fields for proposal.md.
type ProposalFrontmatter struct {
	SpecVersion string `yaml:"spec-version"`
	ChangeName  string `yaml:"change"`
	Status      string `yaml:"status"`
	Created     string `yaml:"created"`
	Updated     string `yaml:"updated"`
}

// SpecFrontmatter holds the YAML frontmatter fields for spec.md files.
type SpecFrontmatter struct {
	Name        string     `yaml:"name,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	GeneratedBy string     `yaml:"generatedBy,omitempty"`
	SpecVersion string     `yaml:"spec-version"`
	Capability  string     `yaml:"capability"`
	Delta       DeltaOp    `yaml:"delta"`
	Status      ItemStatus `yaml:"status"`
}

// TasksFrontmatter holds the YAML frontmatter fields for tasks.md.
// Kept for backward compatibility with Phase 1 callers.
type TasksFrontmatter struct {
	SpecVersion string `yaml:"spec-version"`
	Total       int    `yaml:"total"`
	Completed   int    `yaml:"completed"`
}

// TaskEntry represents a single task with explicit status in tasks.md YAML frontmatter.
// Used by TasksFrontmatterV2 for YAML round-trip task tracking.
type TaskEntry struct {
	ID          int        `yaml:"id"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description,omitempty"`
	Status      ItemStatus `yaml:"status"`
	Spec        string     `yaml:"spec,omitempty"`      // spec directory name this task belongs to
	Depends     []int      `yaml:"depends,omitempty"`   // FSCHEMA-01: task dependency IDs
	Files       []string   `yaml:"files,omitempty"`     // FSCHEMA-02: files touched by this task
	Satisfies   []string   `yaml:"satisfies,omitempty"` // FSCHEMA-03: requirement IDs satisfied
	Skills      []string   `yaml:"skills,omitempty"`    // FSCHEMA-04: slash commands used
}

// TasksFrontmatterV2 extends TasksFrontmatter with a per-task Tasks slice,
// enabling YAML round-trip status updates via updater.go.
type TasksFrontmatterV2 struct {
	SpecVersion string      `yaml:"spec-version"`
	Total       int         `yaml:"total"`
	Completed   int         `yaml:"completed"`
	Tasks       []TaskEntry `yaml:"tasks,omitempty"`
}

// ProposalDoc is the parsed content of a proposal.md file.
type ProposalDoc struct {
	Frontmatter ProposalFrontmatter
	Body        string
}

// DesignDoc is the parsed content of a design.md file.
type DesignDoc struct {
	Body string
}

// Requirement represents a single requirement extracted from a spec file.
type Requirement struct {
	ID         string
	Text       string
	Keyword    RFC2119Keyword
	DeltaOp    DeltaOp
	Status     ItemStatus
	SourceFile string // basename of the source spec file, e.g. "spec.md"
}

// RenamedRequirement represents a rename operation pairing old and new names.
type RenamedRequirement struct {
	From string
	To   string
}

// Task represents a single task entry in tasks.md.
type Task struct {
	ID          int
	Name        string
	Description string
	Status      ItemStatus
	Skipped     bool
	SkipReason  string
}

// Change is the fully assembled representation of a change directory.
type Change struct {
	Name     string
	Dir      string
	Proposal ProposalDoc
	Specs    []Requirement
	Design   DesignDoc
	Tasks    []Task
	Meta     ChangeMeta
}

// Sentinel errors for the spec package.
var (
	ErrNoSpecDir        = errors.New("no spec directory found")
	ErrInvalidTransition = errors.New("invalid state transition")
)
