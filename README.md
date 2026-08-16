# online-exam

一个用 Go 写的内存在线考试与自动判分服务，演示考试、题目、提交、判分与并发判分 worker。

## 功能
- 考试/题目/提交管理
- 按答题对错自动判分
- 并发判分 worker 池，支持 context 取消

## 目录结构
```
cmd/onlineexam/      程序入口
internal/config/     环境配置
internal/model/      模型与纯工具函数
internal/store/      内存存储（考试/题目/提交 + 锁）
internal/service/    业务逻辑
internal/worker/     判分 worker 池
```

## 运行与测试
```bash
go build ./...
go test ./...
```
