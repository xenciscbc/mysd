# mysd

> **測試中** — 本專案正在積極開發中，API 和工作流程可能會變動。

**以規格驅動開發的 AI 程式設計工具**

mysd 是一個 Go CLI 工具 + Claude Code plugin，將 [OpenSpec](https://github.com/openspec-dev/openspec) 的 Spec-Driven Development（SDD）方法論與規劃/執行/驗證引擎整合為一個無縫系統。

它讓獨立開發者（1 人 + N 個 AI agent）能以結構化規格驅動 AI 編程 — 確保 AI 在寫程式前先對齊需求，並在執行後自動驗證成果。

## 為什麼選 mysd？

- **規格是唯一事實來源** — 不只是文件，而是直接驅動 AI 執行和驗證的依據
- **對齊閘門（Alignment Gate）** — AI 必須先讀取並確認規格才能寫任何程式碼（不可繞過）
- **目標回推驗證** — 獨立的 AI agent 逐項檢查每個 MUST 項目是否有檔案系統證據支持
- **慣例優於配置** — 開箱即用，只在需要時才配置
- **Wave 平行執行** — 無相依性的任務在平行 git worktree 中同時執行，大幅縮短執行時間
- **互動式探索** — 透過 advisor agent 進行 4 維度研究（Codebase / Domain / Architecture / Pitfalls）
- **自動更新** — `mysd update` 從 GitHub Releases 檢查更新、下載對應平台的 binary、自動同步 plugin 檔案
- **延遲筆記** — 範圍守衛，在不中斷當前工作的情況下捕捉超出範圍的想法

## 安裝

```bash
go install github.com/xenciscbc/mysd@latest
```

或從 [GitHub Releases](https://github.com/xenciscbc/mysd/releases) 下載預編譯的 binary。

### Claude Code Plugin

```bash
claude plugin add --marketplace https://github.com/xenciscbc/mysd
```

或手動複製 `plugin/` 目錄：

```bash
cp -r plugin/ ~/.claude/plugins/mysd/
```

下次 Claude Code session 即可使用所有 `/mysd:*` 指令。

## 使用模式

mysd 支援四種使用模式，依你的情況選擇。

### 逐步模式（互動式）

完全掌控每個階段。適合複雜功能或需要在步驟間檢視的情況。

```bash
/mysd:init                        # 一次性專案初始化
/mysd:propose add-user-auth       # 建立變更提案 + 產出物
/mysd:discuss                     # （可選）4 維度研究
/mysd:plan                        # 將設計拆解為任務
/mysd:apply                       # 執行任務（含對齊閘門 + 驗證）
/mysd:archive                     # 封存已完成的變更
```

### 快轉模式（半自動）

跳過互動式確認。需要先有一個 active change（先執行 `/mysd:propose`）。

```bash
/mysd:propose my-feature   # 先建立變更
/mysd:ff my-feature        # plan → apply → archive（無研究、自動模式）
```

### 完全快轉模式（全自動）

與快轉模式相同，但在規劃前多一個研究階段。

```bash
/mysd:propose my-feature   # 先建立變更
/mysd:ffe my-feature       # research → plan → apply → archive（自動模式）
```

### 純執行模式

已經有來自其他工具的 specs、設計文件，或對話中的想法？用 `/mysd:plan` 將它們轉換為可執行的任務，然後執行：

```bash
# 從外部文件轉換為任務
/mysd:plan --from design.md          # 載入外部檔案作為 planner context
/mysd:plan --spec auth --from notes  # 對指定 spec 搭配外部輸入進行規劃

# 或者直接在對話中描述需求，然後：
/mysd:plan                           # planner 會擷取對話中的 context

# 如果已經有 tasks.md，直接跳到執行
/mysd:apply             # 執行 tasks.md 中的待辦任務
/mysd:verify            # 獨立驗證 MUST 項目
/mysd:archive           # 完成後封存
```

## 指令一覽

### 核心工作流程

| 指令 | 說明 | 參數 |
|------|------|------|
| `/mysd:propose` | 建立變更提案，自動產生 spec、design、tasks | `[change-name\|file\|dir] [--auto]` |
| `/mysd:discuss` | 臨時研究，4 維度探索與 advisor agent | `[topic\|change-name\|file\|dir] [--auto]` |
| `/mysd:plan` | 將設計拆解為可執行任務，含 MUST 覆蓋率檢查 | `[--research] [--check] [--spec <name>] [--from <file>] [--auto]` |
| `/mysd:apply` | 執行任務，含 spec 對齊閘門；支援 single/wave/spec 模式 | `[--auto]` |
| `/mysd:verify` | 由獨立 verifier agent 對所有 MUST 項目進行目標回推驗證 | |
| `/mysd:archive` | 封存已完成的變更到 `openspec/changes/archive/`，同步 delta spec | `[--auto]` |

### 快轉

| 指令 | 說明 | 參數 |
|------|------|------|
| `/mysd:ff` | 快轉：plan → apply → archive（假設 spec 已就緒、無研究、自動模式） | `[change-name]` |
| `/mysd:ffe` | 完全快轉：research → plan → apply → archive（含研究、自動模式） | `[change-name]` |

### 文件管理

| 指令 | 說明 | 參數 |
|------|------|------|
| `/mysd:docs` | 管理 `docs_to_update` 清單（封存後自動更新的檔案） | `[add <path> \| remove <path>]` |
| `/mysd:docs-update` | 獨立觸發文件更新 — 支援多種範圍 | `[--change <name> \| --last N \| --full \| "text"]` |

`/mysd:docs-update` 範圍說明：
- **無參數** — 從最近一次封存的變更更新
- `--change <name>` — 從指定的封存變更更新
- `--last N` — 從最近 N 次封存的變更更新
- `--full` — 掃描 codebase，依專案實際狀態更新文件
- `"自由文字"` — 以提供的描述作為更新 context

### 工具指令

| 指令 | 說明 | 參數 |
|------|------|------|
| `/mysd:status` | 顯示工作流程狀態、任務進度、下一步建議 | |
| `/mysd:scan` | 掃描現有 codebase 並產生 OpenSpec 格式的 spec | |
| `/mysd:fix` | 在 worktree 隔離環境中修復失敗的任務，可選研究模式 | `[change-name] [T{id}]` |
| `/mysd:note` | 管理延遲筆記 — 捕捉超出範圍的想法 | `[add {content} \| delete {id}]` |
| `/mysd:model` | 檢視或設定 model profile（quality / balanced / budget） | |
| `/mysd:lang` | 設定回應語言和 OpenSpec locale | |
| `/mysd:update` | 檢查更新並安裝新 binary + 同步 plugin 檔案 | `[--check] [--force]` |
| `/mysd:init` | 初始化專案配置和 openspec 結構 | |
| `/mysd:uat` | 互動式使用者驗收測試 | |
| `/mysd:statusline` | 切換狀態列顯示 | `[on\|off]` |

## 運作原理

### 生命週期

每個程式碼變更遵循結構化的生命週期：

```
propose → [discuss] → plan → apply（含驗證）→ archive
```

1. **Propose** — 建立變更提案，自動產生所有產出物（proposal.md、specs/、design.md、tasks.md）於 `openspec/changes/<name>/`
2. **Discuss** *（可選）* — 執行 4 維度互動式研究（Codebase/Domain/Architecture/Pitfalls），在確定需求前探索未知
3. **Plan** — 將設計拆解為有序、可執行的任務；plan-checker 驗證每個 MUST 項目都有對應的任務
4. **Apply** — AI 先讀取 spec（對齊閘門），然後實作每個任務。驗證在執行後自動進行。
5. **Archive** — 將已完成的變更移至 `openspec/changes/archive/YYYY-MM-DD-<name>/`，同步 delta spec 回主要 specs

### 架構

- **Go binary** 處理狀態管理、spec 解析、配置和結構化 JSON 輸出
- **SKILL.md 檔案** 編排 AI 工作流程（呼叫 binary、呈現結果、委派給 agent）
- **Agent 定義**（13 個 agent）執行實際的 AI 工作（規格撰寫、執行、驗證、研究等）
- **反向呼叫模式** — Claude Code 呼叫 binary，而非反過來。不需要 MCP server。

### Agent 角色

| Agent | 角色 |
|-------|------|
| mysd-proposal-writer | 從使用者描述撰寫變更提案 |
| mysd-spec-writer | 使用 RFC 2119 關鍵字撰寫需求規格 |
| mysd-designer | 架構和技術設計 |
| mysd-researcher | 4 維度 codebase 研究 |
| mysd-advisor | 灰色地帶的取捨分析 |
| mysd-planner | 任務拆解和相依性分析 |
| mysd-plan-checker | 驗證計畫涵蓋所有 MUST 項目 |
| mysd-executor | 依計畫實作任務 |
| mysd-verifier | MUST 項目的目標回推驗證 |
| mysd-reviewer | 程式碼審查 |
| mysd-scanner | 掃描 codebase 以產生 spec |
| mysd-uat-guide | 使用者驗收測試引導 |
| mysd-fast-forward | 編排加速工作流程 |

### 帶 Context 的規劃

`/mysd:plan` 支援多種方式將 context 輸入 planner：

- **Per-spec 規劃**（`--spec <name>`）— 限制規劃範圍為單一 spec capability。適合多 spec 的變更需要增量規劃時使用。
- **外部輸入**（`--from <file>`）— 載入檔案（如設計文件、會議記錄、或其他工具產出的計畫）作為 planner 的額外 context。planner 會將此與 spec 產出物一起使用來產生任務。
- **互動式 spec 選擇** — 當存在多個 spec 且未指定 `--spec` 時，互動式選擇器讓你選擇要規劃哪些 spec。
- **研究階段**（`--research`）— 在規劃前執行一輪聚焦的架構研究，適合複雜或不熟悉的領域。
- **Plan checker**（`--check`）— 規劃後由獨立 agent 驗證 spec 中每個 MUST 項目都有對應的任務。

- **對話 context** — 在對話中描述需求或想法，然後執行 `/mysd:plan`。planner 會擷取討論 context 並搭配現有的 spec 產出物一起使用。

規劃管線還包含自動 self-review（placeholder 偵測、一致性檢查、範圍警告、模糊修正）和 reviewer agent 審查。

### Wave 平行執行

當 spec 的 tasks.md 包含 `depends` 欄位時，mysd 會進行相依性分析：

- 任務按拓撲排序分組為 wave — 每個 wave 包含無相互依賴的任務
- 同一 wave 中無檔案重疊的任務在平行 git worktree 中同時執行，縮短執行時間
- 所有 worktree 完成後，branch 按任務 ID 順序合併（`--no-ff`）
- AI 衝突解決（3 次重試，`go build` + test）；失敗的 worktree 保留供 `/mysd:fix` 使用
- 模式選擇：`sequential`（安全，預設）或 `wave`（平行）。`auto_mode` 跳過模式提示。

### 互動式探索

`/mysd:discuss` 在你確定需求前執行雙迴圈研究：

- **4 個研究維度**：Codebase（現有模式）、Domain（需求和限制）、Architecture（方案選項）、Pitfalls（已知失敗模式）
- Advisor agent 依維度浮現未知項目；你在有足夠 context 時終止每個迴圈
- 產出直接輸入 spec 產出物，不會遺失任何資訊

## 自動更新

mysd 可從 GitHub Releases 自動更新：

```bash
/mysd:update           # 互動式檢查更新並安裝
mysd update --check    # 僅檢查（JSON 輸出，不安裝）
mysd update --force    # 不確認直接更新
```

更新包含：

- 帶 SHA256 checksum 驗證的 binary 替換
- 失敗時自動回滾
- 透過 manifest diff 同步 plugin 檔案（commands + agents）— 只寫入有變更的檔案

## OpenSpec 相容性

mysd 原生讀寫 [OpenSpec](https://github.com/openspec-dev/openspec) 格式：

- 支援 `.specs/` 和 `openspec/` 兩種目錄結構
- 解析帶有 `spec-version` 欄位的 YAML frontmatter
- 處理 RFC 2119 關鍵字（MUST、SHOULD、MAY）和 Delta Specs（ADDED、MODIFIED、REMOVED）
- 將 mysd 指向現有的 OpenSpec 專案，直接執行 `/mysd:apply` 或 `/mysd:verify` 無需遷移

## Model Profiles

mysd 使用 profile 系統控制每個 agent 角色使用的 AI model。需要深度思考的角色（規格撰寫、規劃、驗證）使用較強的 model；執行角色使用效率較高的 model。

```bash
/mysd:model              # 檢視目前的 profile 和角色對應表
/mysd:model set quality  # 切換 profile
```

三種 profile 可用：

### quality（8 opus / 2 sonnet）

最大能力。所有思考角色使用 opus。

| 角色 | Model | 用途 |
|------|-------|------|
| spec-writer | opus | 撰寫需求規格 |
| designer | opus | 架構和技術設計 |
| planner | opus | 任務拆解和相依性分析 |
| executor | sonnet | 依計畫實作任務 |
| verifier | opus | 驗證 spec 滿足度 |
| fast-forward | sonnet | 編排加速工作流程 |
| researcher | opus | 4 維度 codebase 研究 |
| advisor | opus | 灰色地帶的取捨分析 |
| proposal-writer | opus | 撰寫變更提案 |
| plan-checker | opus | 驗證計畫涵蓋所有 MUST 項目 |

### balanced（6 opus / 4 sonnet）— 預設

判斷/設計/閘門角色使用 opus，執行和研究使用 sonnet。

| 角色 | Model |
|------|-------|
| spec-writer | opus |
| designer | opus |
| planner | opus |
| executor | sonnet |
| verifier | opus |
| fast-forward | sonnet |
| researcher | sonnet |
| advisor | opus |
| proposal-writer | sonnet |
| plan-checker | opus |

### budget（7 sonnet / 3 haiku）

最小化成本。spec-writer 以 sonnet 為品質底線。

| 角色 | Model |
|------|-------|
| spec-writer | sonnet |
| designer | haiku |
| planner | sonnet |
| executor | haiku |
| verifier | sonnet |
| fast-forward | haiku |
| researcher | sonnet |
| advisor | sonnet |
| proposal-writer | sonnet |
| plan-checker | sonnet |

獨立指令使用固定 model，不受 profile 影響：`init`、`scan`、`fix` 固定使用 opus；`status`、`lang`、`model`、`note`、`docs`、`update` 固定使用 sonnet。

## 配置

專案配置位於 `.claude/mysd.yaml`：

```yaml
execution_mode: single      # single | wave
agent_count: 1
atomic_commits: false
tdd_mode: false
model_profile: balanced     # quality | balanced | budget
response_language: en       # BCP 47 語言標籤，如 zh-TW、ja、fr
```

- `execution_mode: wave` 啟用平行 worktree 執行（針對無相依性重疊的任務）
- `model_profile` 控制所有 agent 角色的 AI model 選擇 — 詳見 [Model Profiles](#model-profiles) 完整對應表
- `response_language` 設定所有 agent 回應和 OpenSpec locale 的語言

所有選項都可透過 flags 在各指令中覆蓋。

## 技術棧

- **Go 1.23+** — 單一 binary，零執行時依賴
- **Cobra** — CLI 框架
- **Viper** — 配置管理
- **lipgloss** — 終端輸出樣式
- **yaml.v3** — OpenSpec frontmatter 的 YAML 解析
- **GoReleaser** — 跨平台 binary 發佈

## 授權

MIT
