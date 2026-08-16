# BUG_REPRO

## Bug 是什么
Store.New 没有初始化 questions map，首次 AddQuestion 往 nil map 写入。

## 如何触发
`go test ./...`

## 错误信息
- panic: assignment to entry in nil map
