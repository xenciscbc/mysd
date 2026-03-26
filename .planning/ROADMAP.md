# Roadmap: my-ssd

## Milestones

- ✅ **v1.0 MVP** — Phases 1-4 (shipped 2026-03-24) — [archive](milestones/v1.0-ROADMAP.md)
- 🚧 **v1.1 Interactive Discovery & Parallel Execution** — Phases 5-9 (in progress)

## Phases

<details>
<summary>✅ v1.0 MVP (Phases 1-4) — SHIPPED 2026-03-24</summary>

- [x] Phase 1: Foundation (3/3 plans) — completed 2026-03-23
- [x] Phase 2: Execution Engine (6/6 plans) — completed 2026-03-24
- [x] Phase 3: Verification & Feedback Loop (5/5 plans) — completed 2026-03-24
- [x] Phase 4: Plugin Layer & Distribution (4/4 plans) — completed 2026-03-24

</details>

### 🚧 v1.1 Interactive Discovery & Parallel Execution (In Progress)

**Milestone Goal:** 讓 mysd 具備互動式需求探索、model profile 分層、並行執行及修復機制

- [x] **Phase 5: Schema Foundation & Plan-Checker** - 擴展 TaskEntry schema + 新增 plan-checker 基礎設施 + model profile 分層 (completed 2026-03-25)
- [x] **Phase 6: Executor Wave Grouping & Worktree Engine** - Topological sort wave 分層 + git worktree 並行執行引擎 (completed 2026-03-25)
- [ ] **Phase 7: New Binary Commands & Scanner Refactor** - model/lang 新指令 + 通用掃描器 + skills 對應流程
- [ ] **Phase 8: SKILL.md Orchestrators & Agent Definitions** - 4 個新 agent + discuss/fix 指令 + auto mode
- [ ] **Phase 9: Interactive Discovery Integration** - propose/spec/discuss 的互動式探索雙模式整合

## Phase Details

### Phase 5: Schema Foundation & Plan-Checker
**Goal**: TaskEntry schema 已擴展並向後相容，plan-checker 可驗證 MUST 覆蓋率，所有新 agent role 已有 model profile 對應
**Depends on**: Phase 4
**Requirements**: FSCHEMA-01, FSCHEMA-02, FSCHEMA-03, FSCHEMA-04, FSCHEMA-05, FSCHEMA-06, FSCHEMA-07, FAGENT-04, FMODEL-01, FMODEL-02, FMODEL-03
**Success Criteria** (what must be TRUE):
  1. 現有 tasks.md 檔案在新版 binary 下可正常讀寫（零值向後相容），含 depends/files/satisfies/skills 欄位的新 tasks.md 可正確序列化和反序列化
  2. 執行 `mysd plan --check` 後，plan-checker 列出所有未被任何 task satisfies 欄位覆蓋的 MUST items，顯示缺口清單並互動式詢問補齊方式
  3. 執行 `mysd plan --context-only` 後，輸出的 JSON 包含 WaveGroups、WorktreeDir、AutoMode 欄位
  4. `mysd.yaml` 的 model profile（quality/balanced/budget）可正確 resolve 到 researcher、advisor、proposal-writer、plan-checker 四個新 agent role
  5. `openspec/config.yaml` 可被 binary 產生和讀取，包含 project metadata 和 locale 欄位
**Plans**: 2 plans
Plans:
- [x] 05-01-PLAN.md — Schema extension + model profile + openspec config writer
- [x] 05-02-PLAN.md — Plan-checker package + cmd wiring + agent definition

### Phase 6: Executor Wave Grouping & Worktree Engine
**Goal**: 依賴關係正確的 tasks 自動分組為波次，每個波次內的 tasks 可在獨立 git worktree 中並行執行，衝突由 AI 自動解決
**Depends on**: Phase 5
**Requirements**: FEXEC-01, FEXEC-02, FEXEC-03, FEXEC-04, FEXEC-05, FEXEC-06, FEXEC-07, FEXEC-08, FEXEC-09, FEXEC-10, FEXEC-11, FEXEC-12
**Success Criteria** (what must be TRUE):
  1. 含 depends 欄位的 tasks.md 執行後，depends 關係正確決定執行順序，同層無依賴的 tasks 分在同一 wave
  2. 同一 wave 內 files 欄位有 overlap 的 tasks 被自動拆到不同 wave，避免並行衝突
  3. 並行執行時每個 task 在 `.worktrees/T{id}/` 建立獨立 worktree，branch 命名為 `mysd/{change-name}/T{id}-{slug}`
  4. 所有 wave 執行完畢後，tasks 依 ID 順序以 `git merge --no-ff` 合入主線；成功後 worktree 和 branch 被自動清理
  5. 執行前磁碟空間不足時顯示可讀錯誤並中止；Windows 環境下 `git config core.longpaths true` 被自動設定；一個 task 失敗不中止同 wave 其他 tasks
**Plans**: 4 plans
Plans:
- [x] 06-01-PLAN.md — Wave grouping algorithm (Kahn's topological sort + file overlap split) + ExecutionContext extension
- [x] 06-02-PLAN.md — Worktree lifecycle package (create/remove/disk space/longpaths) + CLI subcommand
- [x] 06-03-PLAN.md — cmd layer wiring (execute + plan context-only emit real wave_groups)
- [x] 06-04-PLAN.md — SKILL.md execute orchestrator rewrite + executor agent worktree isolation
**UI hint**: no

### Phase 7: New Binary Commands & Scanner Refactor
**Goal**: 使用者可透過 `/mysd:model` 和 `/mysd:lang` 管理設定，scan 支援任意語言，plan 完成後可確認 skills 對應
**Depends on**: Phase 5
**Requirements**: FCMD-03, FCMD-04, FCMD-05, FSCAN-01, FSCAN-02, FSCAN-03, FSCAN-04, FSCAN-05, SKILL-01, SKILL-02, SKILL-03, SKILL-04
**Success Criteria** (what must be TRUE):
  1. 執行 `/mysd:model` 顯示目前 profile 名稱及所有 agent role 的 resolved model；執行 `/mysd:model set quality` 切換 profile 並寫入 mysd.yaml
  2. 執行 `/mysd:lang` 互動式選擇語言後，`mysd.yaml` 的 response_language 和 `openspec/config.yaml` 的 locale 原子同步更新（兩者同時成功或同時不變）
  3. 在 Go/Node.js/Python 等非 Go 專案執行 `/mysd:scan` 後，正確偵測語言並產生 `openspec/config.yaml` + `openspec/specs/` 下的 spec 文件
  4. 已存在 `openspec/config.yaml` 時 scan 只增量更新 specs，config 保持不變
  5. plan 完成後顯示 task 與推薦 skills 的對應表，使用者可逐一調整或批次同意；ffe 模式下跳過互動直接使用推薦值
**Plans**: TBD

### Phase 8: SKILL.md Orchestrators & Agent Definitions
**Goal**: 4 個新 agent definitions 和對應的 SKILL.md 檔案全部就位，discuss/fix 指令可用，auto mode 跨指令運作
**Depends on**: Phase 6, Phase 7
**Requirements**: FCMD-01, FCMD-02, FAGENT-01, FAGENT-02, FAGENT-03, FAGENT-05, FAGENT-06, FAGENT-07, FAUTO-01, FAUTO-02, FAUTO-03, FAUTO-04
**Success Criteria** (what must be TRUE):
  1. 執行 `/mysd:discuss` 後進入討論流程，支援 4 維度並行 research，結論自動觸發 re-plan + plan-checker；無任何 agent definition 包含 Task tool 呼叫
  2. 執行 `/mysd:fix` 後在 worktree 隔離環境修復程式碼，可選 research 模式，修復完成後 worktree 清理
  3. 所有 9 個 agent definitions（含 4 個新 agent）通過手動 audit：無嵌套 subagent spawning，Task tool 只在 SKILL.md orchestrator 層使用
  4. `ff`/`ffe` 指令隱含 `--auto`，跳過 research，直接依照既有 spec 執行；`--auto` flag 在 propose/spec/discuss/plan 跳過互動
**Plans**: TBD

### Phase 9: Interactive Discovery Integration
**Goal**: propose、spec、discuss 三個階段支援互動式探索雙模式，discovery 狀態持久化，scope guardrail 防止 scope creep
**Depends on**: Phase 8
**Requirements**: DISC-01, DISC-02, DISC-03, DISC-04, DISC-05, DISC-06, DISC-07, DISC-08, DISC-09
**Success Criteria** (what must be TRUE):
  1. propose 階段開始時互動式詢問是否使用 research；選擇後，4 個維度（Codebase/Domain/Architecture/Pitfalls）並行啟動 researcher，SKILL.md orchestrator 並行 spawn advisor agents 分析 gray areas
  2. 雙層探索循環有明確終止條件：每輪最多 3 個 gray areas，每個 area 最多 3 個深挖問題；每個 area 完成後呈現「繼續/完成」的二元選擇並顯示剩餘配額
  3. 超出目前 spec 範圍的建議被 redirect 到 deferred notes，不修改當前 spec 內容（scope guardrail 正常運作）
  4. discuss 結論自動更新對應的 spec/design/tasks 檔案，更新後自動執行 re-plan + plan-checker
  5. discovery 狀態持久化為 `discovery-state.json`，research summary 可被 ff/ffe 重用；`--auto` 完全跳過探索循環直接使用 AI 第一推薦
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 5 → 6 → 7 → 8 → 9

| Phase | Milestone | Plans Complete | Status | Completed |
|-------|-----------|----------------|--------|-----------|
| 1. Foundation | v1.0 | 3/3 | Complete | 2026-03-23 |
| 2. Execution Engine | v1.0 | 6/6 | Complete | 2026-03-24 |
| 3. Verification & Feedback Loop | v1.0 | 5/5 | Complete | 2026-03-24 |
| 4. Plugin Layer & Distribution | v1.0 | 4/4 | Complete | 2026-03-24 |
| 5. Schema Foundation & Plan-Checker | v1.1 | 2/2 | Complete   | 2026-03-25 |
| 6. Executor Wave Grouping & Worktree Engine | v1.1 | 4/4 | Complete   | 2026-03-25 |
| 7. New Binary Commands & Scanner Refactor | v1.1 | 4/5 | In Progress|  |
| 8. SKILL.md Orchestrators & Agent Definitions | v1.1 | 0/TBD | Not started | - |
| 9. Interactive Discovery Integration | v1.1 | 0/TBD | Not started | - |
