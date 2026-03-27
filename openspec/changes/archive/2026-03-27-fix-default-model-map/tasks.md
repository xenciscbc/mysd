## 1. Profile model resolution uses short names — 更新 DefaultModelMap

- [x] 1.1 更新 `internal/config/config.go` 中 quality profile 的 DefaultModelMap，確保 profile model resolution uses short names：spec-writer、designer、planner、verifier、researcher、advisor、proposal-writer、plan-checker 改為 `"opus"`，executor 和 fast-forward 維持 `"sonnet"`
- [x] 1.2 更新 balanced profile 的 DefaultModelMap：spec-writer、designer、planner、verifier、advisor、plan-checker 改為 `"opus"`，executor、fast-forward、researcher、proposal-writer 維持 `"sonnet"`
- [x] 1.3 更新 budget profile 的 DefaultModelMap：spec-writer 從 `"haiku"` 改為 `"sonnet"`，其餘不變

## 2. 更新測試

- [x] 2.1 更新 `internal/config/config_test.go` 中所有 ResolveModel 測試案例，使預期值符合新的 DefaultModelMap 分配（quality 思考型 role 預期 opus、balanced 判斷型 role 預期 opus、budget spec-writer 預期 sonnet）
- [x] 2.2 執行 `go test ./internal/config/...` 確認所有測試通過

## 3. 驗證

- [x] 3.1 執行 `go test ./...` 確認全專案測試通過
- [x] 3.2 執行 `go build -o mysd.exe .` 確認編譯成功
- [x] 3.3 執行 `./mysd.exe model` 確認 quality profile 顯示 opus 分配、balanced profile 顯示混合分配、budget profile 顯示修正後的分配
