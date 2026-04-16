---
spec-version: "1.0"
capability: spec-health-check
delta: RENAMED
status: pending
name: Spec Health Check
description: >
  Renamed from "artifact-analysis" to "spec-health-check" to better reflect
  the capability's purpose of analyzing spec quality and completeness.
generatedBy: mysd:spec
---

## Rename Notice

This capability was previously named **artifact-analysis** and has been renamed to **spec-health-check**.

All references to the old name `artifact-analysis` MUST be updated to `spec-health-check` in:
- Spec file paths (`openspec/specs/artifact-analysis/` → `openspec/specs/spec-health-check/`)
- Task `spec:` fields referencing the old capability name
- Any cross-references in other spec files

## Requirements

### REQ-01: Spec Coverage Analysis

The system MUST scan all capabilities and report which ones have spec files and which do not.

#### Scenario: Missing Spec Detected

WHEN a capability directory exists under `cmd/` or `internal/`
AND no corresponding `spec.md` exists under `openspec/specs/`
THEN the health check reports it as a coverage gap

### REQ-02: Requirement Ambiguity Check

The system SHOULD flag requirements that use lowercase modal verbs (must, should, may) instead of RFC 2119 UPPERCASE keywords.

#### Scenario: Lowercase Keyword Flagged

WHEN a spec contains "the system must validate input" (lowercase "must")
THEN the health check flags it as an ambiguity warning

### REQ-03: Scenario Completeness

The system MUST verify that every requirement heading has at least one associated scenario.

#### Scenario: Requirement Without Scenario

WHEN a requirement heading `## Requirement: X` has no `### Scenario:` block following it
THEN the health check reports it as a gap finding
