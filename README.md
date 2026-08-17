# go-alert-evaluator

## 项目说明

go-alert-evaluator 是一个基于内存存储的时间序列告警规则评估服务。它支持写入带标签的指标样本、定义窗口聚合规则、按指定时刻评估告警、配置静默区间，以及查询包含通知决策的状态历史。

服务使用可注入时钟和通知发送器，适合本地开发、监控计算模拟和确定性测试。

## 标准命令

```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd
```

## 使用方式

启动后默认监听 `:8080`。

- `POST /samples`：写入指标样本。
- `POST /rules`：创建告警规则。
- `POST /evaluate`：评估单条规则。
- `POST /evaluate/batch`：批量评估规则。
- `POST /silences`：创建静默区间。
- `GET /alerts/{rule_id}/history`：查询规则状态历史。

时间字段使用 RFC3339 格式，规则窗口和持续时间使用 Go duration 字符串，例如 `5m`、`30s`。
