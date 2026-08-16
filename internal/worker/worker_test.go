package worker

import (
	"context"
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
