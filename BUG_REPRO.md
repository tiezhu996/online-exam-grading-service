# BUG_REPRO

## Bug 是什么
ValidAnswer 校验反向、Submit 丢掉了入参校验、MarkGraded 去掉锁、worker 把 wg.Add 放进 goroutine 且汇总无锁，导致无效答案被接受、提交重复判分、统计错误并存在 data race。

## 如何触发
`go test -race ./...`

## 错误信息
- TestValidAnswer / TestSubmitRejectsInvalid / TestGradeAllSummary 失败。
- `go test -race` 报 data race。
