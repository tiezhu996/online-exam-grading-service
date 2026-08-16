package worker

import (
	"context"
	"fmt"
	"testing"

	"onlineexam/internal/config"
	"onlineexam/internal/model"
	"onlineexam/internal/service"
	"onlineexam/internal/store"
)

func newPool() (*service.Service, *Pool) {
	s := store.New()
	svc := service.New(s, config.Load())
	return svc, New(svc, 4)
}

func TestGradeAllSummary(t *testing.T) {
	svc, p := newPool()
	_ = svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"})
	_ = svc.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 1})
	for i := 0; i < 10; i++ {
		_, _ = svc.Submit(&model.Submission{ID: string(rune('a' + i)), ExamID: "e1", Answers: []model.Answer{{QuestionID: "q1", Correct: true}}})
	}
	sum := p.GradeAll(context.Background())
	if sum.Checked != 10 || sum.Graded != 10 || sum.Failed != 0 {
		t.Fatalf("summary=%+v", sum)
	}
}

func TestGradeAllCancel(t *testing.T) {
	svc, p := newPool()
	_ = svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"})
	_, _ = svc.Submit(&model.Submission{ID: "s1", ExamID: "e1"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sum := p.GradeAll(ctx)
	if sum.Graded != 0 {
		t.Fatalf("graded=%d want 0", sum.Graded)
	}
}

// TestGradeAllConcurrent verifies that under high concurrency every submission
// is graded exactly once and the summary numbers are accurate. Guards against
// the lost-update race and early-return bugs that caused double grading and
// mismatched counts.
func TestGradeAllConcurrent(t *testing.T) {
	svc, p := newPool()
	_ = svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"})
	_ = svc.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 1})
	const n = 200
	for i := 0; i < n; i++ {
		_, _ = svc.Submit(&model.Submission{ID: fmt.Sprintf("s%d", i), ExamID: "e1", Answers: []model.Answer{{QuestionID: "q1", Correct: true}}})
	}

	sum := p.GradeAll(context.Background())
	if sum.Checked != n || sum.Graded != n || sum.Failed != 0 {
		t.Fatalf("summary=%+v want checked=%d graded=%d failed=0", sum, n, n)
	}

	graded := 0
	for _, id := range svc.ListSubmissions("e1") {
		if id.Status == model.StatusGraded {
			graded++
		}
	}
	if graded != n {
		t.Fatalf("graded submissions=%d want %d (double grading)", graded, n)
	}

	// A second pass must be a no-op: nothing pending left to grade.
	sum2 := p.GradeAll(context.Background())
	if sum2.Checked != 0 || sum2.Graded != 0 {
		t.Fatalf("second pass summary=%+v want zeros", sum2)
	}
}
