# Bug 复现

## 现象

批处理中途取消仍返回成功
。

## 触发步骤

```bash
go test -buildvcs=false -count=1 -run "TestEvaluateBatchStopsAfterContextCancellation" ./internal/service/
```

## 基线错误信息

```text
--- FAIL: TestEvaluateBatchStopsAfterContextCancellation (0.00s)
    batch_context_verification_test.go:26: error = <nil>, want context canceled; results=2
--- FAIL: TestBatchReturnsCancellationBeforeEvaluatingRemainingRules (0.00s)
    batch_deadline_private_test.go:25: error = <nil>, want context canceled
FAIL
FAIL	github.com/zhangkui/go-alert-evaluator/internal/service	0.381s
FAIL
```
