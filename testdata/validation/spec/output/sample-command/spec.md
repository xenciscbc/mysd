---
spec-version: "1.0"
capability: sample-command
delta: ADDED
status: pending
name: Sample Command
description: Command behavior for the sample action, inferred from cmd/sample.go.
generatedBy: mysd:spec
---

## Requirements

### REQ-01: Input Validation <!-- inferred -->

The `SampleCommand` function MUST reject empty string input by returning an error.

The error message MUST indicate that input must not be empty.

#### Scenario: Empty Input Rejected

WHEN `SampleCommand` is called with an empty string ""
THEN it returns an error with message "input must not be empty"
AND the result string is empty

### REQ-02: Input Processing and Output Format <!-- inferred -->

The `SampleCommand` function MUST return a result string prefixed with "processed: " followed by the input value.

#### Scenario: Valid Input Processed

WHEN `SampleCommand` is called with input "hello"
THEN it returns ("processed: hello", nil)
AND no error is returned

### REQ-03: Configuration Structure <!-- inferred -->

The `SampleConfig` struct MUST expose a `Verbose` boolean field tagged as `yaml:"verbose"`.

The `SampleConfig` struct MUST expose an `Output` string field tagged as `yaml:"output"`.

#### Scenario: Config YAML Deserialization

WHEN a YAML document contains `verbose: true` and `output: "/tmp/out"`
THEN deserializing into `SampleConfig` sets `Verbose` to true and `Output` to "/tmp/out"
