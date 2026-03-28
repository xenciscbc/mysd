#!/usr/bin/env node
// mysd-statusline: Claude Code statusline for mysd projects
// Shows: model | change | directory | context usage

const fs = require('fs');
const path = require('path');
const os = require('os');

// Read JSON from stdin
let input = '';
// Timeout guard: if stdin doesn't close within 3s (e.g. pipe issues on
// Windows/Git Bash), exit silently instead of hanging. See gsd #775.
const stdinTimeout = setTimeout(() => process.exit(0), 3000);
process.stdin.setEncoding('utf8');
process.stdin.on('data', chunk => input += chunk);
process.stdin.on('end', () => {
  clearTimeout(stdinTimeout);
  try {
    const data = JSON.parse(input);

    // Extract core fields
    const workspace = data.workspace?.current_dir || process.cwd();
    const session = data.session_id || '';
    const remaining = data.context_window?.remaining_percentage;

    // Read statusline_enabled from .claude/mysd.yaml (D-12)
    function readStatuslineEnabled(workspaceDir) {
      try {
        const yamlPath = path.join(workspaceDir, '.claude', 'mysd.yaml');
        const content = fs.readFileSync(yamlPath, 'utf8');
        const match = content.match(/^statusline_enabled:\s*(.+)$/m);
        if (match) {
          const val = match[1].trim().toLowerCase();
          return val !== 'false';
        }
        return true; // not found = enabled
      } catch (e) { return true; }
    }

    // Model shortname extraction (D-02)
    function extractModelShortname(data) {
      const name = (data.model?.display_name || data.model?.id || '').toLowerCase();
      if (name.includes('opus'))   return 'opus';
      if (name.includes('sonnet')) return 'sonnet';
      if (name.includes('haiku'))  return 'haiku';
      const displayName = data.model?.display_name || '';
      const firstWord = displayName.split(' ')[0];
      return firstWord || 'claude';
    }

    // Read change_name from .specs/state.yaml (D-10)
    function readChangeName(workspaceDir) {
      try {
        const stateYamlPath = path.join(workspaceDir, '.specs', 'state.yaml');
        const content = fs.readFileSync(stateYamlPath, 'utf8');
        const match = content.match(/^change_name:\s*["']?([^"'\n\r]+)["']?\s*$/m);
        return match ? match[1].trim() : null;
      } catch (e) { return null; }
    }

    // Detect GSD coexistence by checking for gsd-context-monitor.js (D-04)
    function detectGsdCoexistence(workspaceDir) {
      const homedir = os.homedir();
      const claudeConfigDir = process.env.CLAUDE_CONFIG_DIR || path.join(homedir, '.claude');
      const checkPaths = [
        path.join(workspaceDir, '.claude', 'hooks', 'gsd-context-monitor.js'),
        path.join(claudeConfigDir, 'hooks', 'gsd-context-monitor.js')
      ];
      return checkPaths.some(p => {
        try { return fs.existsSync(p); } catch (e) { return false; }
      });
    }

    // Context window display (shows USED percentage scaled to usable context)
    // Claude Code reserves ~16.5% for autocompact buffer, so usable context
    // is 83.5% of the total window. We normalize to show 100% at that point.
    const AUTO_COMPACT_BUFFER_PCT = 16.5;
    let ctx = '';
    if (remaining != null) {
      // Normalize: subtract buffer from remaining, scale to usable range
      const usableRemaining = Math.max(0, ((remaining - AUTO_COMPACT_BUFFER_PCT) / (100 - AUTO_COMPACT_BUFFER_PCT)) * 100);
      const used = Math.max(0, Math.min(100, Math.round(100 - usableRemaining)));

      // Bridge file: only when GSD coexists AND session exists (D-04)
      const gsdCoexists = detectGsdCoexistence(workspace);
      if (session && gsdCoexists) {
        try {
          const bridgePath = path.join(os.tmpdir(), `claude-ctx-${session}.json`);
          fs.writeFileSync(bridgePath, JSON.stringify({
            session_id: session,
            remaining_percentage: remaining,
            used_pct: used,
            timestamp: Math.floor(Date.now() / 1000)
          }));
        } catch (e) { /* silent */ }
      }

      // Build progress bar (10 segments)
      const filled = Math.floor(used / 10);
      const bar = '\u2588'.repeat(filled) + '\u2591'.repeat(10 - filled);

      // Color thresholds (D-05): hot face emoji instead of skull for >=80%
      if (used < 50) {
        ctx = ` \x1b[32m${bar} ${used}%\x1b[0m`;
      } else if (used < 65) {
        ctx = ` \x1b[33m${bar} ${used}%\x1b[0m`;
      } else if (used < 80) {
        ctx = ` \x1b[38;5;208m${bar} ${used}%\x1b[0m`;
      } else {
        ctx = ` \x1b[5;31m\uD83E\uDD75 ${bar} ${used}%\x1b[0m`;  // hot face emoji (not skull)
      }
    }

    // Check statusline_enabled AFTER bridge file write (D-12)
    const statuslineEnabled = readStatuslineEnabled(workspace);
    if (!statuslineEnabled) {
      process.stdout.write('');
      return;
    }

    // Build and output statusline string (D-01)
    const model = extractModelShortname(data);
    const dirname = path.basename(workspace);
    const changeName = readChangeName(workspace);

    // Format: {model} | {change} | {dir} | {bar} {pct}% (change omitted when absent)
    if (changeName) {
      process.stdout.write(`\x1b[2m${model}\x1b[0m \u2502 \x1b[1m${changeName}\x1b[0m \u2502 \x1b[2m${dirname}\x1b[0m${ctx}`);
    } else {
      process.stdout.write(`\x1b[2m${model}\x1b[0m \u2502 \x1b[2m${dirname}\x1b[0m${ctx}`);
    }
  } catch (e) {
    // Silent fail -- don't break statusline on parse errors
  }
});
