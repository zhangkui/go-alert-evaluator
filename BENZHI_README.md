# BENZHI 评测说明

## 1. 项目说明

项目：`zhangkui/go-alert-evaluator`。该项目是内存时间序列告警规则评估 Web API，Go 工具链为 `golang:1.22`，前端工具链：无。

## 2. 标准构建、运行和测试命令

```bash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## 3. Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-alert-evaluator-bug5-candidate linux/amd64
./build_benzhi_docker.sh go-alert-evaluator-bug5-candidate linux/arm64
docker run --rm -it --platform linux/amd64 go-alert-evaluator-bug5-candidate bash
docker run --rm -it --platform linux/arm64 go-alert-evaluator-bug5-candidate bash
```

## 4. 题目验证命令

该 diagnosis candidate 保持原始源码以供复现；构建预期退出码为 0，目标验证、私有验证和全量测试预期为非 0。

``bash
GOTOOLCHAIN=local go build -buildvcs=false ./...
GOTOOLCHAIN=local go test -buildvcs=false -count=1 -run "Test
GOTOOLCHAIN=local go test -buildvcs=false -count=20 -run private ./...
GOTOOLCHAIN=local go test -buildvcs=false -count=1 ./...
``

## 5. Bug 复现

复现步骤和基线错误信息见 `BUG_REPRO.md`。
