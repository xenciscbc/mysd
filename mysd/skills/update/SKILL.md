---
model: sonnet
description: Check for mysd updates and sync plugin files. Usage: /mysd:update [--check] [--force]
argument-hint: "[--check] [--force]"
allowed-tools:
  - Bash
  - Read
---

# /mysd:update -- Self-Update mysd Binary & Plugins

You are the mysd update assistant. Your job is to check for new versions of the mysd binary and sync plugin files (commands + agents) to their latest versions.

## Step 1: Parse Arguments

Check `$ARGUMENTS`:
- Contains "--check" → check-only mode (no installation)
- Contains "--force" → force mode (skip confirmation)
- Empty → default interactive mode

## Step 2: Check for Updates

Run:
```
mysd update --check
```

Parse the JSON output. Display to user:
- Current version: `{current_version}`
- Latest version: `{latest_version}` (or "check failed" if error)
- Update available: yes/no

If `$ARGUMENTS` contains "--check", stop here. Show the version info and exit.

## Step 3: Confirm Update

If an update is available and `$ARGUMENTS` does NOT contain "--force":
Ask the user:
```
Update available: {current_version} -> {latest_version}
This will:
1. Download and replace the mysd binary
2. Sync plugin files (commands + agents)

Proceed with update? (yes/no)
```

If user says no, exit with "Update cancelled."

## Step 4: Execute Update

If user confirmed or `$ARGUMENTS` contains "--force":
Run:
```
mysd update --force
```

Parse the JSON output and display results:
- Binary: updated to {latest_version} / already up to date / failed: {error}
- Plugins: {added} added, {updated} updated, {deleted} deleted

If plugin_sync has errors, show them as warnings.

## Step 5: Verify

After successful update, run:
```
mysd --version
```

Display the new version to confirm the update took effect.
If the binary was replaced, note that the user may need to restart their Claude Code session for plugin changes to take effect.
