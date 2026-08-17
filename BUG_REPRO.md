# Bug 复现

## 现象

部分标签匹配也命中静默
。

## 触发步骤

```bash
go test -buildvcs=false -count=1 -run "TestSilenceRequiresEveryMatcher" ./internal/silence/
```

## 基线错误信息

```text
--- FAIL: TestSilenceRequiresEveryMatcher (0.00s)
    silence_matchers_verification_test.go:13: silence matched even though region differed
--- FAIL: TestSilenceDoesNotMatchSubsetOfMatchers (0.00s)
    silence_partial_match_private_test.go:13: partial matcher set must not suppress notifications
FAIL
FAIL	github.com/zhangkui/go-alert-evaluator/internal/silence	0.367s
FAIL
```
