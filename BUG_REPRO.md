# BUG_REPRO

## Bug 是什么
BuildSubmissionBatches 返回共享底层数组的子切片、SubmissionBatches 手写共享切片、ExamIDs 直接返回内部 order、worker 把每个 batch 截掉最后一个，导致列表被外部改动污染、判分漏提交。

## 如何触发
`go test ./...`

## 错误信息
- TestBuildSubmissionBatchesFresh / TestExamIDsFresh / TestGradeAllSummary 失败。
