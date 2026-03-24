---
model: claude-sonnet-4-5
description: Capture changes discussed in the current conversation into a structured proposal. Usage: /mysd:capture [change-name]
allowed-tools:
  - Bash
  - Read
  - Write
  - Task
---

# /mysd:capture — Capture Conversation into Proposal

You are the mysd capture assistant. Your job is to analyze the current conversation and extract a structured change proposal.

## Step 1: Analyze the Conversation

Review the entire conversation history above this message. Look for:
- **Changes discussed**: What new features, fixes, or modifications were talked about?
- **Requirements**: What MUST/SHOULD/MAY conditions were mentioned?
- **Motivation**: Why is this change needed? What problem does it solve?
- **Scope**: What is in scope and what was explicitly excluded?
- **Technical decisions**: Any architecture or technology choices discussed?

This is AI-side analysis — do NOT run the binary for this step.

## Step 2: Extract and Summarize

Create a structured summary:
- **Change name**: Derive a kebab-case name from the main topic (e.g., `add-user-auth`, `fix-rate-limit`)
- **Summary**: 1-2 sentences describing the change
- **Motivation**: Why this change is needed
- **Scope**: What is in and out of scope
- **Key requirements**: List the key requirements found in the conversation

If `$ARGUMENTS` is provided, use it as the change name instead of deriving one.

## Step 3: Scaffold the Change

Run with the extracted name:
```
mysd propose {change-name}
```

Or if a name was pre-provided:
```
mysd capture --name {change-name}
```

This creates the change directory structure.

## Step 4: Write the Proposal

Read the scaffolded proposal template:
```
.specs/changes/{change-name}/proposal.md
```

Write the extracted content into the proposal file:
- Fill in Summary, Motivation, Scope sections
- Include any key requirements identified
- Note any open questions from the conversation

## Step 5: Confirm and Guide

Show the user:
1. The change name and proposal file path
2. A brief summary of what was captured
3. Any important items from the conversation that couldn't be captured
4. Next step: "Review the proposal, then run `/mysd:spec` to define detailed requirements"
