# BUG_REPRO

## Bug 是什么
GradeSubmission 对 GetSubmission 的错误做了吞掉处理：出错时返回 (0, nil) 而不是向上传播，导致调用方拿不到错误、只看到 0 分。

## 如何触发
`go test ./...`

## 错误信息
- TestGradeMissingWraps 失败（期望 ErrSubmissionNotFound，实际得到 nil error）。
