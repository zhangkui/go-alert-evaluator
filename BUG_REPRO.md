# Bug 复现

## 现象

首次健康评估立即恢复
。

## 触发步骤

```bash
go test -buildvcs=false -count=1 -run "TestFiringWaitsForRecoveryDuration" ./internal/engine/
```

## 基线错误信息

```text
--- FAIL: TestFiringWaitsForRecoveryDuration (0.00s)
    recovery_duration_verification_test.go:24: state = normal, want firing until recover_for elapses
--- FAIL: TestRecoveryTimerResetsWhenThresholdBreachesAgain (0.00s)
    recovery_reset_private_test.go:22: early recovery state normal
FAIL
FAIL	github.com/zhangkui/go-alert-evaluator/internal/engine	0.345s
FAIL
```
