package verifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateScenarioFormat_Complete(t *testing.T) {
	body := `#### Scenario: User logs in

- **GIVEN** a registered user
- **WHEN** they submit valid credentials
- **THEN** the system returns a session token
`
	warnings := ValidateScenarioFormat(body)
	assert.Empty(t, warnings)
}

func TestValidateScenarioFormat_MissingGiven(t *testing.T) {
	body := `#### Scenario: User logs in

- **WHEN** they submit valid credentials
- **THEN** the system returns a session token
`
	warnings := ValidateScenarioFormat(body)
	assert.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "GIVEN")
	assert.Contains(t, warnings[0], "User logs in")
}

func TestValidateScenarioFormat_MissingWhenAndThen(t *testing.T) {
	body := `#### Scenario: Incomplete

- **GIVEN** something exists
`
	warnings := ValidateScenarioFormat(body)
	assert.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "WHEN")
	assert.Contains(t, warnings[0], "THEN")
}

func TestValidateScenarioFormat_MultipleScenarios(t *testing.T) {
	body := `#### Scenario: Good one

- **GIVEN** a thing
- **WHEN** action happens
- **THEN** result occurs

#### Scenario: Bad one

- **WHEN** action happens
- **THEN** result occurs
`
	warnings := ValidateScenarioFormat(body)
	assert.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "Bad one")
	assert.Contains(t, warnings[0], "GIVEN")
}

func TestValidateScenarioFormat_NoScenarios(t *testing.T) {
	body := `# Spec

Some content without scenarios.
`
	warnings := ValidateScenarioFormat(body)
	assert.Empty(t, warnings)
}
