# Bug 复现

## 现象

窗口起点样本被排除
。

## 触发步骤

```bash
go test -buildvcs=false -count=1 -run "TestAggregateIncludesSampleAtWindowStart" ./internal/evaluator/
```

## 基线错误信息

```text
--- FAIL: TestCountIncludesOnlyInclusiveWindowSamples (0.00s)
    window_boundary_count_private_test.go:14: count = 1, want 2
--- FAIL: TestAggregateIncludesSampleAtWindowStart (0.00s)
    window_boundary_verification_test.go:14: average = 8, want 5 with inclusive window start
FAIL
FAIL	github.com/zhangkui/go-alert-evaluator/internal/evaluator	0.407s
FAIL
```
