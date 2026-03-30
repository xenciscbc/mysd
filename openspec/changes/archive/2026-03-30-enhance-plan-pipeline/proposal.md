---
spec-version: "1"
change: enhance-plan-pipeline
type: feature
status: proposed
---

## Why

mysd plan 目前的 artifact 產出缺乏結構化指引（agent 自由發揮），品質把關只靠後端 reviewer + analyze-fix，沒有即時自檢。此外，planning 和 execution 都是以整個 change 為單位，無法針對單一 spec 獨立進行，也無法接入外部 planning 工具（如 gstack plan-eng-review）的產出。

## What Changes

- 新增 `mysd instructions <artifact-id> --change <name> --json` CLI 指令，為每個 artifact 提供結構化的 template、rules、instruction、dependencies
- 在 plan SKILL.md 的 orchestrator 層加入 Inline Self-Review 步驟（4 項檢查：placeholders、consistency、scope、ambiguity），在 planner 完成後、reviewer 之前執行
- 在 `mysd instructions` 的 rules 中內嵌 self-review checklist，讓 agent 寫作時就遵守
- TasksFrontmatterV2 新增 `spec` 欄位，標記每個 task 所屬的 spec
- `mysd plan` 和 `mysd apply` 支援 `--spec X` 直接指定 spec，或不指定時 interactive 列出可選 specs
- `mysd plan --from <path>` 支援接受外部 plan/tasks 作為 planner 的 context 輸入（如 gstack 產出）

## Capabilities

### New Capabilities

- `artifact-instructions`: CLI 指令 `mysd instructions`，為 designer/planner agent 提供結構化的 template、rules、instruction、dependencies JSON 輸出
- `inline-self-review`: plan orchestrator 在 planner 完成後、reviewer 之前執行 4 項品質自檢（placeholders、consistency、scope、ambiguity）並直接修復

### Modified Capabilities

- `planning`: 支援 per-spec planning（`--spec X`）、external input（`--from <path>`）、interactive spec selection、以及使用 `mysd instructions` 輸出指引 agent
- `execution`: 支援 per-spec execution（`--spec X`）、interactive spec selection；TasksFrontmatterV2 新增 `spec` 欄位

## Impact

- 新增 Go code: `cmd/instructions.go`（新 CLI 指令）
- 修改 Go code: `cmd/plan.go`（`--spec`、`--from` flags）、`cmd/execute.go`（`--spec` flag、context 過濾）
- 修改 Go code: `internal/tasks/`（TasksFrontmatterV2 加 `spec` 欄位解析）
- 修改 skill: `mysd/skills/plan/SKILL.md`（instructions 呼叫、self-review 步驟、spec selection）
- 修改 skill: `mysd/skills/apply/SKILL.md`（spec selection）
- 修改 agent: `mysd/agents/mysd-planner.md`（接收 instructions 結構）
- 修改 agent: `mysd/agents/mysd-designer.md`（接收 instructions 結構）
- 影響 spec: `openspec/specs/planning/spec.md`、`openspec/specs/execution/spec.md`
