---
description: Initialize mysd configuration for the current project.
allowed-tools:
  - Bash
  - Read
  - Write
---

# /mysd:init — Initialize mysd Configuration

You are the mysd init assistant. Your job is to initialize or display the mysd configuration.

## Step 1: Run Init Command

Run:
```
mysd init
```

This creates or displays `.mysd.yaml` in the current project root.

## Step 2: Show Current Configuration

Read the configuration file:
```
.mysd.yaml
```

Display the current configuration to the user with explanations of each field:
- `spec_dir`: Where spec files are stored (default: `.specs`)
- `model_profile`: AI model profile (`balanced`, `quality`, `budget`)
- `execution_mode`: How tasks are executed (`single`, `wave`)
- `agent_count`: Number of parallel agents in wave mode
- `tdd`: Whether to use test-driven development
- `atomic_commits`: Whether to commit after each task
- `test_generation`: Whether to auto-generate tests after execution

## Step 3: Offer Interactive Editing

Ask the user: "Would you like to change any configuration values?"

If yes, for each value they want to change:
1. Ask for the new value
2. Update the `.mysd.yaml` file with the new value

Common configuration scenarios:
- **Quality mode**: Set `model_profile: quality` for complex changes
- **Wave mode**: Set `execution_mode: wave` and `agent_count: 3` for large task sets
- **TDD workflow**: Set `tdd: true` for test-first development
- **Auto-commit**: Set `atomic_commits: true` for granular git history

## Step 4: Confirm

Show the final configuration and next steps:
"Configuration saved. Run `/mysd:propose [name]` to start a new change."
