## ADDED Requirements

### Requirement: Custom profile definition in configuration

The system SHALL allow users to define custom model profiles in the `custom_profiles` section of `mysd.yaml`. Each custom profile SHALL have a `base` field specifying an existing built-in profile name and a `models` field containing role-to-model mappings that override the base profile.

#### Scenario: Valid custom profile definition

- **WHEN** the user defines a custom profile in `mysd.yaml` with `base: balanced` and `models: { executor: opus }`
- **THEN** the system SHALL load the custom profile without error

#### Scenario: Custom profile with empty models

- **WHEN** the user defines a custom profile with a valid `base` and an empty `models` map
- **THEN** the system SHALL treat the custom profile as identical to its base profile

### Requirement: Custom profile model resolution

The system SHALL resolve models for a custom profile using the following priority order:
1. `ModelOverrides[role]` (per-role override)
2. `CustomProfiles[profile].Models[role]` (custom profile's model mapping)
3. `DefaultModelMap[base][role]` (base profile's mapping)
4. `"sonnet"` (fallback)

#### Scenario: Role overridden in custom profile

- **WHEN** the active profile is a custom profile with `base: balanced` and `models: { executor: opus }`
- **THEN** the system SHALL resolve the `executor` role to `opus`

#### Scenario: Role not overridden in custom profile

- **WHEN** the active profile is a custom profile with `base: balanced` and `models: { executor: opus }`
- **THEN** the system SHALL resolve the `planner` role to the value defined in the `balanced` built-in profile

#### Scenario: ModelOverrides takes precedence over custom profile

- **WHEN** the active profile is a custom profile with `models: { executor: opus }` and `ModelOverrides` contains `executor: haiku`
- **THEN** the system SHALL resolve the `executor` role to `haiku`

### Requirement: Custom profile selection via CLI

The `mysd model set <name>` command SHALL accept custom profile names in addition to built-in profile names. The lookup order SHALL be: built-in profiles first, then custom profiles.

#### Scenario: Set a custom profile

- **WHEN** the user runs `mysd model set my-team` and `my-team` is defined in `custom_profiles`
- **THEN** the system SHALL set `model_profile` to `my-team` and confirm success

#### Scenario: Set an unknown profile

- **WHEN** the user runs `mysd model set nonexistent` and `nonexistent` is neither a built-in nor a custom profile
- **THEN** the system SHALL return an error listing all available profiles (built-in and custom)

### Requirement: Custom profile display

The `mysd model` command SHALL display the resolved model for each role when the active profile is a custom profile, showing the profile name.

#### Scenario: Display custom profile

- **WHEN** the active profile is `my-team` (a custom profile with `base: balanced`)
- **THEN** the system SHALL display `Profile: my-team` and list all roles with their resolved models

### Requirement: Invalid role name warning

The system SHALL emit a warning when a custom profile's `models` map contains a role name that is not in the known roles list.

#### Scenario: Unknown role in custom profile

- **WHEN** the user defines a custom profile with `models: { excutor: opus }` (typo)
- **THEN** the system SHALL emit a warning indicating `excutor` is not a known role

### Requirement: Invalid base profile warning

The system SHALL emit a warning when a custom profile's `base` field references a profile name that is not a built-in profile.

#### Scenario: Unknown base profile

- **WHEN** the user defines a custom profile with `base: premium`
- **THEN** the system SHALL emit a warning indicating `premium` is not a valid base profile
