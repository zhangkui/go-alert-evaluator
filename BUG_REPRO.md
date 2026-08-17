# Bug 复现

## 现象

重复时间戳写入返回两个点
。

## 触发步骤

```bash
go test -buildvcs=false -count=1 -run "TestDuplicateTimestampReplacesExistingSample" ./internal/store/
```

## 基线错误信息

```text
--- FAIL: TestDuplicatePointDoesNotInflateWindowCardinality (0.00s)
    duplicate_sample_count_private_test.go:17: unexpected duplicate handling: []model.Sample{model.Sample{Timestamp:time.Date(2026, time.August, 17, 8, 50, 0, 0, time.UTC), Value:5}, model.Sample{Timestamp:time.Date(2026, time.August, 17, 8, 50, 0, 0, time.UTC), Value:8}}
--- FAIL: TestDuplicateTimestampReplacesExistingSample (0.00s)
    duplicate_sample_verification_test.go:17: duplicate timestamp produced 2 points, want 1
FAIL
FAIL	github.com/zhangkui/go-alert-evaluator/internal/store	0.333s
FAIL
```
