# Phase 6: Executor Wave Grouping & Worktree Engine - Research

**Researched:** 2026-03-25
**Domain:** Go topological sort, git worktree lifecycle, concurrent execution, Windows path constraints
**Confidence:** HIGH (direct codebase analysis + verified git worktree behavior)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Worktree 管理主體**
- D-01: Go binary 的 `internal/worktree/` 新 package 負責所有 git worktree lifecycle：create（`git worktree add`）、remove（`git worktree remove`）、disk space check（FEXEC-10）、Windows longpaths 設定（`git config core.longpaths true`，FEXEC-11）。SKILL.md 透過呼叫 `mysd worktree` subcommand 委派。符合 v1.0 binary-as-state-manager 架構原則，邏輯可測試、deterministic。
- D-02: Worktree 路徑：`.worktrees/T{id}/`（短路徑，Windows MAX_PATH 相容，沿用 ROADMAP 規格）。Branch 命名：`mysd/{change-name}/T{id}-{slug}`（沿用 ROADMAP 規格）。

**執行模式切換 UX**
- D-03: 只在有實際並行機會時才詢問執行模式：
  - Tasks 都無 `depends` 且都無 `files` 欄位 → 直接 sequential，不詢問
  - Tasks 有 `depends` 或 `files` 欄位（有並行機會）→ 詢問「Sequential（安全穩定）/ Wave parallel（N 個 tasks 並行）」
- D-04: `ffe` 和 `--auto` 跳過詢問：有 `depends`/`files` 時用 wave mode，否則用 sequential。符合 FAUTO-03/FAUTO-04 的 auto 語義。

**Merge 衝突失敗 UX**
- D-05: AI 自動解衝突最多嘗試 3 次（每次：解衝突 → build + test 驗證），失敗後：
  - 保留 worktree 在 `.worktrees/T{id}/`
  - 顯示清楚的錯誤訊息：失敗原因 + worktree 路徑 + branch 名稱 + 建議下一步（`cd .worktrees/T{id}` 手動解衝突）
  - 該 wave 其他已成功的 tasks 照常 merge（continue-on-failure policy，沿用 Phase 5 討論決策）
  - 無自動 resume 機制，人工解決後再次執行
- D-06: Merge 成功的 worktree 自動刪除（FEXEC-08），失敗的保留。

**Wave 執行進度顯示**
- D-07: 使用 lipgloss Printer 輸出 inline status：
  - Wave 開始：`Wave 1/3: T1, T2, T3 並行執行中...`
  - Task 完成：inline 輸出 `T1 ✓`（成功）或 `T1 ✗`（失敗）
  - Wave 結束：摘要行 `Wave 1 complete: 2 succeeded, 1 failed`
  - 格式簡潔，符合現有 lipgloss Printer 模式

### Claude's Discretion
- `internal/worktree/` 的 Go interface 設計細節（struct vs function-based API）
- disk space check 的臨界值（多少 MB 算不足）
- topological sort 演算法選擇（Kahn's vs DFS — Phase 5 D-04 已標示 Claude discretion）
- wave grouping 時 files overlap 的比較邏輯（exact match vs prefix match）
- `mysd worktree` subcommand 的 CLI 介面細節

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FEXEC-01 | Wave grouping 演算法依 `depends` 做 topological sort 分層 | Kahn's algorithm 用 BFS，純函數 `BuildWaveGroups(tasks []TaskItem) [][]TaskItem` in `internal/executor/waves.go` |
| FEXEC-02 | 同層 tasks 檢查 `files` overlap，有 overlap 拆到不同 wave | 在 `BuildWaveGroups` 後執行 file-overlap split pass；exact string match 足夠（決策由 Claude discretion） |
| FEXEC-03 | 每個並行 task spawn executor with `isolation: "worktree"` | SKILL.md `mysd-execute.md` 修改：依 `wave_groups` 並行 Task → mysd-executor，傳入 `worktree_path` 和 `branch` |
| FEXEC-04 | Worktree branch 命名 `mysd/{change-name}/T{id}-{task-slug}` | `internal/worktree/worktree.go` `Create()` 函數計算 slug，`git worktree add -b {branch} {path}` |
| FEXEC-05 | Worktree 建在 `.worktrees/T{id}/`（短路徑，Windows 相容） | `internal/worktree/worktree.go` 使用 `cfg.WorktreeDir`（預設 `.worktrees`）+ `T{id}` 子路徑 |
| FEXEC-06 | 合併依 task ID 順序，`git merge --no-ff` | SKILL.md merge loop：按升序 task ID 執行 `git merge --no-ff {branch}` |
| FEXEC-07 | AI 自動解衝突 → build + test 驗證 → 失敗 AI 修復 → 最多 3 次 → 仍失敗通知使用者 | SKILL.md mysd-execute.md 的 merge loop 中加入 retry logic |
| FEXEC-08 | 成功自動刪除 worktree + branch；失敗保留供檢查 | `internal/worktree/worktree.go` `Remove(id)` = `git worktree remove {path}` + `git branch -d {branch}` |
| FEXEC-09 | Wave 中一個 task 失敗，其他繼續跑完 | SKILL.md：continue-on-failure；collect all results，merge successful ones only |
| FEXEC-10 | Worktree 建立前檢查磁碟空間（disk space guard） | `internal/worktree/worktree.go` 使用 `syscall.Statfs` (Unix) / `GetDiskFreeSpaceEx` via `golang.org/x/sys/windows` |
| FEXEC-11 | Windows worktree 自動設定 `git config core.longpaths true` | `internal/worktree/worktree.go` `Create()` 時偵測 `runtime.GOOS == "windows"` 後執行 |
| FEXEC-12 | Executor 遵守 task 的 `skills` 欄位，執行時使用指定的 slash commands | `mysd-executor.md` 修改：加入「如 `assigned_task.skills` 非空，優先使用指定 skills」指示 |
</phase_requirements>

---

## Summary

Phase 6 建立 mysd 的 worktree 並行執行引擎。核心工作分三個方向：

**Go binary 層（`internal/`）：** 新增 `internal/executor/waves.go` 實作 `BuildWaveGroups()` — 用 Kahn's algorithm 對 `TaskItem.Depends` 做 topological sort 分層，再用 file overlap pass 拆解同層衝突。新增 `internal/worktree/` package 負責完整 git worktree lifecycle：create、remove、disk space check、Windows longpaths 設定。`ExecutionContext` 擴展 `WaveGroups` 欄位供 SKILL.md 消費。

**SKILL.md 層（`plugin/commands/mysd-execute.md`）：** 重寫 execute orchestrator 加入 wave 執行模式選擇邏輯、每 wave 並行 spawn executors、merge loop（升序 task ID、retry 3 次、continue-on-failure）、進度顯示。

**Agent 層（`plugin/agents/mysd-executor.md`）：** 擴展 executor agent 支援 worktree isolation mode — 接收 `worktree_path` 和 `branch` 參數，在指定 worktree 目錄執行、commit，並在完成後回報狀態。

**Primary recommendation:** 先實作並測試 `BuildWaveGroups()` 純函數（無外部依賴，易測試），再實作 `internal/worktree/`（需要 git 存在但可用 test repo 隔離測試），最後修改 SKILL.md 和 agent。

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib `os/exec` | stdlib | 執行 `git worktree add/remove/list` | 所有 git 操作都是子進程呼叫，無需 libgit2 |
| Go stdlib `sync` | stdlib | `sync.WaitGroup` + `sync.Mutex` 追蹤 wave 並發結果 | Wave 並發結果收集的標準工具 |
| Go stdlib `runtime` | stdlib | `runtime.GOOS` 偵測 Windows | OS 偵測不需要外部依賴 |
| `golang.org/x/sys/windows` | 已在間接依賴中 | `GetDiskFreeSpaceExW` 取得磁碟空間 | Windows API 最精確；Unix 用 `syscall.Statfs` |
| `github.com/stretchr/testify` | v1.11.1 | 單元測試 assert/require | 已有，專案標準 |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/charmbracelet/lipgloss` | v1.1.0 | Wave 進度顯示 | `Printer.Printf` + 直接 lipgloss styles（已有 `internal/executor/status.go` 中的 styles） |

**無需新增 go.mod 依賴。** 所有必要工具已在現有 `go.mod` 中。

---

## Architecture Patterns

### Recommended Project Structure（Phase 6 新增）

```
internal/
├── executor/
│   ├── context.go         # 擴展：加 WaveGroups [][]TaskItem, WorktreeDir, AutoMode
│   ├── waves.go           # 新增：BuildWaveGroups(), HasParallelOpportunity()
│   └── waves_test.go      # 新增：topological sort + file overlap 測試
│
└── worktree/
    ├── worktree.go        # 新增：WorktreeManager, Create(), Remove(), CheckDiskSpace()
    └── worktree_test.go   # 新增：使用 t.TempDir() + git init 的整合測試

plugin/
├── commands/
│   └── mysd-execute.md    # 修改：wave mode 選擇 + merge loop
└── agents/
    └── mysd-executor.md   # 修改：worktree isolation + skills 欄位遵守

cmd/
└── execute.go             # 小修：把 WaveGroups 注入 ExecutionContext
```

### Pattern 1: BuildWaveGroups — Kahn's Algorithm

**What:** Kahn's algorithm (BFS-based topological sort) 計算 wave layers，再執行 file overlap split pass。
**When to use:** `executor.BuildContextFromParts()` 呼叫後，填入 `ExecutionContext.WaveGroups`。

**Kahn's algorithm 選擇理由（Claude's discretion 說明）：**
- 天生產出 BFS 層（即 wave layers），DFS 需要額外 post-processing
- 對 cycle detection 更直觀（in-degree 降為 0 時才入隊）
- 純函數，易於測試

```go
// Source: 直接實作，無外部依賴
// internal/executor/waves.go

// BuildWaveGroups computes dependency-ordered wave layers from pending tasks.
// Uses Kahn's algorithm (BFS topological sort) to produce layers,
// then splits any same-layer tasks with file overlap into separate waves.
func BuildWaveGroups(tasks []TaskItem) [][]TaskItem {
    if len(tasks) == 0 {
        return nil
    }

    // Build index and in-degree map
    idToTask := make(map[int]TaskItem, len(tasks))
    inDegree  := make(map[int]int, len(tasks))
    adj       := make(map[int][]int) // id -> ids that depend on it

    for _, t := range tasks {
        idToTask[t.ID] = t
        if _, ok := inDegree[t.ID]; !ok {
            inDegree[t.ID] = 0
        }
        for _, dep := range t.Depends {
            adj[dep] = append(adj[dep], t.ID)
            inDegree[t.ID]++
        }
    }

    // BFS layer extraction
    var layers [][]TaskItem
    queue := []int{}
    for id, deg := range inDegree {
        if deg == 0 {
            queue = append(queue, id)
        }
    }
    // Sort for determinism
    sort.Ints(queue)

    for len(queue) > 0 {
        layerIDs := make([]int, len(queue))
        copy(layerIDs, queue)
        queue = queue[:0]

        var layer []TaskItem
        for _, id := range layerIDs {
            layer = append(layer, idToTask[id])
            for _, next := range adj[id] {
                inDegree[next]--
                if inDegree[next] == 0 {
                    queue = append(queue, next)
                }
            }
        }
        sort.Ints(queue)
        sort.Slice(layer, func(i, j int) bool { return layer[i].ID < layer[j].ID })
        layers = append(layers, layer)
    }

    // File overlap split pass
    return splitByFileOverlap(layers)
}

// HasParallelOpportunity returns true if any task has Depends or Files set.
// Used by SKILL.md decision point (D-03): only show wave mode prompt when true.
func HasParallelOpportunity(tasks []TaskItem) bool {
    for _, t := range tasks {
        if len(t.Depends) > 0 || len(t.Files) > 0 {
            return true
        }
    }
    return false
}
```

**File overlap split pass（Claude's discretion — 選 exact string match）：**

```go
// splitByFileOverlap ensures no two tasks in the same wave touch the same file.
// Uses exact string matching (case-sensitive) — sufficient for file path comparison.
func splitByFileOverlap(layers [][]TaskItem) [][]TaskItem {
    var result [][]TaskItem
    for _, layer := range layers {
        result = append(result, splitLayer(layer)...)
    }
    return result
}

func splitLayer(tasks []TaskItem) [][]TaskItem {
    var sublayers [][]TaskItem
    for _, t := range tasks {
        placed := false
        for i := range sublayers {
            if !hasFileConflict(sublayers[i], t) {
                sublayers[i] = append(sublayers[i], t)
                placed = true
                break
            }
        }
        if !placed {
            sublayers = append(sublayers, []TaskItem{t})
        }
    }
    return sublayers
}

func hasFileConflict(layer []TaskItem, t TaskItem) bool {
    fileSet := make(map[string]struct{})
    for _, existing := range layer {
        for _, f := range existing.Files {
            fileSet[f] = struct{}{}
        }
    }
    for _, f := range t.Files {
        if _, ok := fileSet[f]; ok {
            return true
        }
    }
    return false
}
```

### Pattern 2: WorktreeManager — 薄 os/exec 包裝

**What:** `internal/worktree/` package 是 `os/exec` 的薄包裝，每個方法執行一個 git 子命令。
**When to use:** 由 `cmd/execute.go` 或 SKILL.md 透過 `mysd worktree` subcommand 呼叫。

```go
// Source: 直接實作
// internal/worktree/worktree.go

type WorktreeManager struct {
    RepoRoot    string // abs path to repo root
    WorktreeDir string // relative, e.g. ".worktrees"
    ChangeName  string
}

// Create adds a new worktree at .worktrees/T{id}/ on branch mysd/{change}/T{id}-{slug}.
// On Windows, also sets core.longpaths=true.
func (m *WorktreeManager) Create(id int, taskName string) (path string, branch string, err error) {
    // 1. Set longpaths on Windows
    if runtime.GOOS == "windows" {
        if err := m.setLongPaths(); err != nil {
            return "", "", fmt.Errorf("set longpaths: %w", err)
        }
    }
    // 2. Check disk space (500MB threshold — Claude discretion)
    if err := m.CheckDiskSpace(500 * 1024 * 1024); err != nil {
        return "", "", err
    }
    // 3. Compute paths
    slug := toSlug(taskName)
    worktreePath := filepath.Join(m.RepoRoot, m.WorktreeDir, fmt.Sprintf("T%d", id))
    branchName   := fmt.Sprintf("mysd/%s/T%d-%s", m.ChangeName, id, slug)
    // 4. git worktree add -b {branch} {path}
    cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
    cmd.Dir = m.RepoRoot
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", "", fmt.Errorf("git worktree add: %w\n%s", err, out)
    }
    return worktreePath, branchName, nil
}

// Remove deletes the worktree directory and branch (success cleanup).
func (m *WorktreeManager) Remove(id int, branch string) error {
    worktreePath := filepath.Join(m.RepoRoot, m.WorktreeDir, fmt.Sprintf("T%d", id))
    // git worktree remove --force {path}
    rmCmd := exec.Command("git", "worktree", "remove", "--force", worktreePath)
    rmCmd.Dir = m.RepoRoot
    if out, err := rmCmd.CombinedOutput(); err != nil {
        return fmt.Errorf("git worktree remove: %w\n%s", err, out)
    }
    // git branch -d {branch}
    delCmd := exec.Command("git", "branch", "-d", branch)
    delCmd.Dir = m.RepoRoot
    if out, err := delCmd.CombinedOutput(); err != nil {
        // Non-fatal: worktree already removed, branch delete failure is tolerated
        fmt.Fprintf(os.Stderr, "warning: git branch -d %s: %v\n%s\n", branch, err, out)
    }
    return nil
}
```

**Disk space check — 實作方式（Claude's discretion：500MB 臨界值）：**

```go
// CheckDiskSpace returns error if available bytes < minBytes.
func (m *WorktreeManager) CheckDiskSpace(minBytes uint64) error {
    available, err := getAvailableBytes(m.RepoRoot)
    if err != nil {
        // Non-fatal if we can't determine disk space
        return nil
    }
    if available < minBytes {
        return fmt.Errorf(
            "insufficient disk space: need at least %dMB, have %dMB available at %s",
            minBytes/(1024*1024), available/(1024*1024), m.RepoRoot,
        )
    }
    return nil
}

// getAvailableBytes returns available bytes at the given path.
// Uses syscall.Statfs on Unix, GetDiskFreeSpaceExW on Windows.
```

### Pattern 3: ExecutionContext WaveGroups 擴展

`BuildContextFromParts()` 在 Phase 6 需要計算並填入 `WaveGroups`：

```go
// internal/executor/context.go — 擴展

type ExecutionContext struct {
    // ... 所有現有欄位不變 ...
    WaveGroups  [][]TaskItem `json:"wave_groups,omitempty"`  // NEW Phase 6
    WorktreeDir string       `json:"worktree_dir,omitempty"` // NEW Phase 5 (已加)
    AutoMode    bool         `json:"auto_mode,omitempty"`    // NEW Phase 5 (已加)
}

// BuildContextFromParts 修改：在填完 PendingTasks 後加：
// ctx.WaveGroups = BuildWaveGroups(ctx.PendingTasks)
// ctx.WorktreeDir = cfg.WorktreeDir
// ctx.AutoMode    = cfg.AutoMode
```

注意：`WorktreeDir` 和 `AutoMode` 已在 `ProjectConfig` 中（Phase 5 完成），但 `ExecutionContext` 還沒有這兩個欄位 — Phase 6 需補上。

### Pattern 4: SKILL.md Wave 執行流程

```
// plugin/commands/mysd-execute.md 修改後的邏輯

Step 1: mysd execute --context-only → JSON
  - 解析 wave_groups, worktree_dir, auto_mode, change_name

Step 2: 決定執行模式（D-03）
  - has_parallel_opportunity = wave_groups 非空且有多個 tasks
  - if NOT has_parallel_opportunity → skip ask, use sequential
  - if auto_mode → use wave (if wave_groups non-empty), else sequential
  - else → ask user: "Sequential / Wave parallel ({N} tasks)"

Step 3A: Sequential mode
  - 單一 Task → mysd-executor (全部 pending_tasks)

Step 3B: Wave mode
  For wave_index, wave in wave_groups:
    Print: "Wave {i+1}/{total}: T{ids}... 並行執行"
    Spawn parallel Tasks → mysd-executor (per task, with worktree params)
    Wait for all to complete (continue-on-failure)

    For task_id in ascending order:
      if task succeeded:
        git merge --no-ff {branch}
        if conflict: retry up to 3 times (resolve → go build → go test)
        if resolved: mysd worktree remove {id} {branch}
        if not resolved after 3 tries: preserve worktree, show error

    Print: "Wave {i+1} complete: {succeeded} succeeded, {failed} failed"

Step 4: Post-execution summary
```

### Pattern 5: mysd-executor Worktree Isolation

`mysd-executor.md` 擴展 — worktree isolation 模式下的 input 和行為：

```markdown
// 新增 input 欄位
- `worktree_path`: 若非空，所有操作在此目錄中執行（非 repo root）
- `branch`: 已建立的 branch 名稱
- `isolation`: "worktree" | "none"

// 執行行為
- 若 isolation == "worktree"，切換到 worktree_path 執行所有 Bash 命令
- task-update 仍呼叫 `mysd task-update {id} done`（在 worktree 目錄有效）
- atomic_commits 改為 commit 到 worktree 的 branch
- 完成後輸出結構化 JSON 讓 SKILL.md 判斷成功/失敗
```

### Anti-Patterns to Avoid

- **Worktrees 跨 wave 前啟動：** 下一 wave 的 worktrees 必須等上一 wave 所有 merge 完成後才建立（否則 branch 從舊 HEAD 分叉，產生額外 merge commit）。參見 ARCHITECTURE.md Anti-Pattern 3。
- **Binary 讀取 stdin：** `cmd/execute.go` 不得加 interactive prompt。所有使用者決策（sequential vs wave）由 SKILL.md 問完後以 flag 傳入。
- **WorktreeManager 匯入 executor：** 依賴圖必須單向。`worktree` package 不得 import `executor`；`executor` 的 `waves.go` 是純函數，不 import `worktree`。
- **Cycle detection 跳過：** `BuildWaveGroups` 必須偵測 cycle（Kahn's algorithm：若最後有 node 的 in-degree > 0 即有 cycle），回傳明確 error 而非 silent hang。

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Git worktree 路徑管理 | 自己 mkdir + symlink | `git worktree add -b {branch} {path}` | git 確保 worktree metadata (.git/worktrees/) 正確；手動 mkdir 不會更新 .git |
| Topological sort | 複雜遞迴圖演算法 | Kahn's algorithm（~40 行 Go） | 教科書演算法，無外部依賴；不要用 third-party graph library |
| Disk space check on Windows | 自己呼叫 WMI | `GetDiskFreeSpaceExW` via `golang.org/x/sys/windows` | 唯一可靠方式；`os.Stat` 不提供可用空間 |
| slug 計算 | 正規表達式輪子 | `strings.Map` + `strings.ToLower` | slug 只需 [a-z0-9-]，stdlib 足夠 |
| Concurrent result collection | Channel pipeline 設計 | `sync.WaitGroup` + `[]WaveResult` + `sync.Mutex` | Wave 是 simple fan-out/fan-in，不需要複雜的 channel 設計 |

**Key insight:** 所有複雜的並行執行邏輯（spawn agents、retry、conflict resolution）都在 SKILL.md 中，因為這些是 AI 決策流程。Go binary 只負責 deterministic 的 git 操作和計算，這些可以被完整單元測試。

---

## Runtime State Inventory

> 本 phase 涉及 worktree 目錄的建立和清理，但不包含 rename/refactor。無需完整盤點。

| Category | Items Found | Action Required |
|----------|-------------|------------------|
| Stored data | `ExecutionContext` JSON schema 擴展 (`wave_groups` 欄位新增) | 舊 SKILL.md 消費此 JSON 時不受影響（omitempty，新欄位不存在時為 null） |
| Live service config | 無 | 無 |
| OS-registered state | git worktrees 建立後在 `.git/worktrees/` 有 metadata | `git worktree remove` 或 `git worktree prune` 清理 |
| Secrets/env vars | 無 | 無 |
| Build artifacts | 無 | 無 |

**Windows worktree 路徑上限：** `.worktrees/T{id}/` 最短路徑設計已在 CONTEXT.md D-02 確認。在 Windows 上 `git config core.longpaths true` 由 `internal/worktree/` 自動設定（FEXEC-11）。當前環境 `core.longpaths` 為未設定狀態，驗證 Phase 6 的自動設定邏輯是必要的。

---

## Common Pitfalls

### Pitfall 1: Cycle detection 缺失導致無限迴圈

**What goes wrong:** `BuildWaveGroups` 沒有偵測 cycle。若 tasks.md 有 A depends B, B depends A，Kahn's algorithm 的 queue 會空掉但還有 node 在 inDegree > 0，導致回傳空或部分結果，silently 跳過部分 tasks。
**Why it happens:** 開發時只測試正常 DAG。
**How to avoid:** Kahn's algorithm 完成後，檢查 `len(processed) != len(tasks)`；若不等，回傳 `ErrCyclicDependency`。
**Warning signs:** `wave_groups` 中 task 數量少於 `pending_tasks` 數量。

### Pitfall 2: git worktree remove 失敗 silently 留下 orphan worktrees

**What goes wrong:** merge 成功後呼叫 `git worktree remove`，但 worktree 目錄中有 untracked 或 modified files，git 拒絕 remove（exit 1）。若不處理，`.git/worktrees/` 內會留下 orphan metadata。
**Why it happens:** executor agent 在 worktree 中 build 時產生 `*.test` 或 `vendor/` 等臨時檔案。
**How to avoid:** 使用 `git worktree remove --force`（強制清理）。並在 `Remove()` 之前呼叫 `git worktree prune` 清理 stale metadata。
**Warning signs:** `git worktree list` 顯示已刪除路徑的條目。

### Pitfall 3: Wave 間 HEAD 不一致導致 merge conflict 放大

**What goes wrong:** Wave 0 執行時，同時開始建立 Wave 1 的 worktrees。Wave 0 merge 後，Wave 1 的 branches 從舊 HEAD 分叉，merge 時出現本不必要的衝突。
**Why it happens:** SKILL.md 邏輯錯誤地提前建立下一 wave worktrees。
**How to avoid:** SKILL.md 的 wave loop 必須先完成整個 wave 的 merge 步驟，才進入下一波次的 worktree 建立。明確的 wave boundary contract。
**Warning signs:** Wave N 的 merge 產生大量無關 files 的衝突。

### Pitfall 4: Windows MAX_PATH 在 worktree 內部的嵌套路徑

**What goes wrong:** `.worktrees/T1/` 的根路徑雖短，但 worktree 內部的 Go 套件路徑（例如 `internal/some-long-package-name/`）加上 repo 根路徑本身，仍可能超過 260 字元。
**Why it happens:** repo 本身位於深層路徑（如 `D:\work_data\project\go\mysd\`）。
**How to avoid:** `core.longpaths true` 解決 git 層的問題。Go build 工具在 Windows 上也需要路徑夠短；建議在 CI 環境驗證（STATE.md 已記錄此 concern）。`core.longpaths` 的設定必須是 global（`git config --global`）或 project-level；建議 project-level 以免影響其他 repo。

### Pitfall 5: disk space check 跨平台 API 差異

**What goes wrong:** 用 `os.Stat` 或 `os.Getwd()` 加 filepath 計算磁碟空間，在 Windows 上不可用。
**Why it happens:** 開發者在 Linux 用 `syscall.Statfs` 實作，Windows build 失敗。
**How to avoid:** 用 build tags 分離實作：`worktree_unix.go` (`//go:build !windows`) 用 `syscall.Statfs`；`worktree_windows.go` (`//go:build windows`) 用 `golang.org/x/sys/windows` 的 `GetDiskFreeSpaceExW`。`golang.org/x/sys/windows` 已在 go.mod 間接依賴中。

---

## Code Examples

### WaveGroups JSON 輸出格式（SKILL.md 消費）

```json
// mysd execute --context-only 輸出（Phase 6 後）
{
  "change_name": "my-feature",
  "pending_tasks": [
    {"id": 1, "name": "setup auth", "depends": [], "files": ["auth.go"]},
    {"id": 2, "name": "setup cache", "depends": [], "files": ["cache.go"]},
    {"id": 3, "name": "wire up", "depends": [1, 2], "files": ["main.go"]}
  ],
  "wave_groups": [
    [
      {"id": 1, "name": "setup auth", "files": ["auth.go"]},
      {"id": 2, "name": "setup cache", "files": ["cache.go"]}
    ],
    [
      {"id": 3, "name": "wire up", "depends": [1, 2], "files": ["main.go"]}
    ]
  ],
  "worktree_dir": ".worktrees",
  "auto_mode": false
}
```

### git worktree 操作序列

```bash
# Create worktree (FEXEC-04, FEXEC-05)
# Windows: first set longpaths
git config core.longpaths true

# Create on branch from current HEAD
git worktree add -b mysd/my-feature/T1-setup-auth .worktrees/T1

# Merge after task completes (FEXEC-06)
git merge --no-ff mysd/my-feature/T1-setup-auth

# Cleanup on success (FEXEC-08)
git worktree remove --force .worktrees/T1
git branch -d mysd/my-feature/T1-setup-auth
git worktree prune
```

### HasParallelOpportunity 使用方式

```go
// cmd/execute.go 或 SKILL.md 消費後的決策：
// ExecutionContext.WaveGroups 長度 > 0 且至少一個 wave 有 2+ tasks
// → SKILL.md 詢問 sequential vs wave

// 或用 HasParallelOpportunity() 在 binary 層預計算
// ExecutionContext 中加入 has_parallel_opportunity bool 欄位供 SKILL.md 直接消費
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `wave_groups: [][]int{}` placeholder（plan.go L75） | `[][]TaskItem` 實際計算（Phase 6） | Phase 6 | SKILL.md 終於能讀到真實 wave 資訊 |
| `execution_mode: "wave"` 靠 agent_count 分配 task（舊 execute.md） | Wave groups 依 depends/files 計算（Phase 6） | Phase 6 | 並行是 dependency-aware，不再是 round-robin |
| mysd-executor 不知道 worktree 概念 | mysd-executor 在 `worktree_path` 下工作 | Phase 6 | executor agent 真正隔離執行 |

---

## Open Questions

1. **`mysd worktree` subcommand 是否必要？**
   - What we know: CONTEXT.md D-01 說「SKILL.md 透過呼叫 `mysd worktree` subcommand 委派」，但 SKILL.md 也可以直接 call git 命令
   - What's unclear: worktree subcommand 的 CLI 介面細節（Claude's discretion）
   - Recommendation: 實作最小 `mysd worktree create {id} {task-name}` 和 `mysd worktree remove {id} {branch}` subcommand，讓 SKILL.md 呼叫；disk space check 和 longpaths 在 Create 內部自動處理。這比 SKILL.md 直接 git 更容易測試和 debug。

2. **`wave_groups` 在 `ExecutionContext` 中的計算時機**
   - What we know: `BuildContextFromParts` 是純函數，已接收全部 tasks；`BuildWaveGroups` 也是純函數
   - What's unclear: `plan.go --context-only` 中的 `wave_groups: [][]int{}` placeholder 是否要一併修改
   - Recommendation: Phase 6 同時修改 `execute --context-only` 輸出和 `plan --context-only` 輸出中的 `wave_groups`，確保兩個 context 的 schema 一致（plan context 的 wave_groups 改為 `[][]TaskItem` 並實際計算）。

3. **disk space check 臨界值**
   - What we know: Claude's discretion；需要留給 worktree clone + build artifacts
   - Recommendation: 500MB 作為預設，可透過 `ProjectConfig` 中未來擴展的 `worktree_min_disk_mb` 欄位配置。Phase 6 先 hardcode 500MB，Phase 7+ 再開放配置。

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go 1.23+ | Binary compilation | ✓ | go1.25.5 | — |
| git | git worktree operations | ✓ | 2.43.0.windows.1 | — |
| git worktree | FEXEC-03~08 | ✓ | 2.43.0（git 2.5+ 支援） | — |
| core.longpaths | FEXEC-11 Windows | ✓ (設定後) | 當前未設定（Phase 6 自動設定） | — |
| Disk space（153GB available） | FEXEC-10 | ✓ | 153GB free | — |
| `golang.org/x/sys/windows` | FEXEC-10 Windows | ✓ | 間接依賴已在 go.mod | — |

**Missing dependencies with no fallback:** 無

**Note:** 當前環境已有多個 `.claude/worktrees/agent-*` worktrees（`git worktree list` 顯示 9+ 個），確認 `git worktree` 功能在此環境正常運作。

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go testing + testify v1.11.1 |
| Config file | none — 使用 `go test ./...` |
| Quick run command | `go test ./internal/executor/... ./internal/worktree/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FEXEC-01 | `BuildWaveGroups` 依 depends 正確分層 | unit | `go test ./internal/executor/... -run TestBuildWaveGroups` | ❌ Wave 0 |
| FEXEC-01 | Cycle detection 回傳 error | unit | `go test ./internal/executor/... -run TestBuildWaveGroups_Cycle` | ❌ Wave 0 |
| FEXEC-02 | Files overlap 拆到不同 wave | unit | `go test ./internal/executor/... -run TestBuildWaveGroups_FileOverlap` | ❌ Wave 0 |
| FEXEC-02 | 無 overlap 不增加 wave | unit | `go test ./internal/executor/... -run TestBuildWaveGroups_NoOverlap` | ❌ Wave 0 |
| FEXEC-04 | Branch 命名 `mysd/{change}/T{id}-{slug}` | unit | `go test ./internal/worktree/... -run TestCreate_BranchName` | ❌ Wave 0 |
| FEXEC-05 | Worktree 路徑 `.worktrees/T{id}/` | unit | `go test ./internal/worktree/... -run TestCreate_Path` | ❌ Wave 0 |
| FEXEC-08 | Remove 後 worktree 目錄消失 | integration | `go test ./internal/worktree/... -run TestRemove` | ❌ Wave 0 |
| FEXEC-10 | Disk space check 在可用空間不足時回傳 error | unit | `go test ./internal/worktree/... -run TestCheckDiskSpace` | ❌ Wave 0 |
| FEXEC-11 | Windows 環境 Create 時設定 core.longpaths | unit | `go test ./internal/worktree/... -run TestCreate_WindowsLongPaths` | ❌ Wave 0 |
| FEXEC-03,06,07,09,12 | Wave parallel execute + merge + retry + continue | manual-only | 需要實際 Claude Code + Task tool 執行 | N/A |

**Manual-only 說明（FEXEC-03,06,07,09,12）：** 這些需求的核心邏輯在 SKILL.md orchestrator 中（AI 並行 spawn、AI 解衝突），無法透過 Go 單元測試驗證。驗收標準是執行 `/mysd:execute` 後觀察 wave 行為。

### Sampling Rate
- **Per task commit:** `go test ./internal/executor/... ./internal/worktree/... -count=1`
- **Per wave merge:** `go test ./... -count=1`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/executor/waves.go` — `BuildWaveGroups`, `HasParallelOpportunity`, `splitByFileOverlap`
- [ ] `internal/executor/waves_test.go` — covers FEXEC-01, FEXEC-02
- [ ] `internal/worktree/worktree.go` — `WorktreeManager`, `Create`, `Remove`, `CheckDiskSpace`
- [ ] `internal/worktree/worktree_test.go` — covers FEXEC-04, FEXEC-05, FEXEC-08, FEXEC-10, FEXEC-11
- [ ] Build tag 分離：`internal/worktree/diskspace_unix.go` + `internal/worktree/diskspace_windows.go`

---

## Project Constraints (from CLAUDE.md)

- **Tech stack:** Go — 單一 binary，跨平台編譯（不加新的外部依賴）
- **Plugin 形式:** Claude Code slash commands + agent definitions（SKILL.md 是 orchestrator，不做業務邏輯）
- **Binary-as-state-manager:** cmd/ 只解析參數和輸出；業務邏輯在 internal/ packages
- **Convention over configuration:** WorktreeDir 預設 `.worktrees`；Windows longpaths 自動設定不需要使用者配置
- **Pure function packages:** `waves.go` 的 `BuildWaveGroups` 和 `HasParallelOpportunity` 必須是純函數（no I/O）
- **Sidecar pattern:** 不修改 tasks.md 的 status 以外的欄位；wave 狀態追蹤不需要新 sidecar（SKILL.md 在 memory 中管理 wave 進度）
- **Additive-only extension:** 新欄位（`WaveGroups`, `WorktreeDir`, `AutoMode`）使用 `omitempty`，確保向後相容

---

## Sources

### Primary (HIGH confidence)
- 直接 codebase 分析 — `internal/executor/context.go`（TaskItem 結構、BuildContextFromParts 邏輯）
- 直接 codebase 分析 — `internal/executor/waves.go`（不存在，需建立）
- 直接 codebase 分析 — `internal/config/defaults.go`（WorktreeDir, AutoMode 已存在）
- 直接 codebase 分析 — `go.mod`（確認 `golang.org/x/sys` 間接依賴已在 go.mod）
- `.planning/research/ARCHITECTURE.md` — v1.1 架構研究，`internal/worktree/` 完整設計
- `.planning/phases/06-executor-wave-grouping-worktree-engine/06-CONTEXT.md` — 所有 D-01~D-07 決策
- `git worktree list` 實際執行 — 確認環境中 git worktree 可用（2.43.0.windows.1）

### Secondary (MEDIUM confidence)
- Go stdlib `os/exec` 文件 — git subcommand 呼叫模式
- Kahn's algorithm — 教科書拓樸排序演算法，多個來源確認

### Tertiary (LOW confidence)
- Windows `GetDiskFreeSpaceExW` via `golang.org/x/sys/windows` — 已知 API，但在此專案環境尚未實測

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — go.mod 直接確認，無新外部依賴
- Architecture (waves.go): HIGH — 純函數演算法，有現有 planchecker 作為模式參考
- Architecture (worktree package): HIGH — ARCHITECTURE.md 有完整藍圖，git worktree 命令有文件
- SKILL.md flow: MEDIUM — 邏輯清楚，但 AI orchestration 細節需實際測試
- Windows disk space check: MEDIUM — API 已知，但此環境未實測 build

**Research date:** 2026-03-25
**Valid until:** 2026-04-25（git worktree API 穩定，30 天有效）
