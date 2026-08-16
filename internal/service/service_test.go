package service

import (
	"errors"
	"testing"

	"onlineexam/internal/config"
	"onlineexam/internal/model"
	"onlineexam/internal/store"
)

func newSvc() (*store.Store, *Service) {
	s := store.New()
	return s, New(s, config.Load())
}

func TestCreateAndQuestions(t *testing.T) {
	_, svc := newSvc()
	if err := svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 5}); err != nil {
		t.Fatal(err)
	}
	if err := svc.AddQuestion(&model.Question{}); err == nil {
		t.Fatal("invalid question accepted")
	}
	if len(svc.ListQuestions("e1")) != 1 {
		t.Fatal("ListQuestions len")
	}
}

func TestSubmitRejectsInvalid(t *testing.T) {
	_, svc := newSvc()
	if _, err := svc.Submit(&model.Submission{}); err == nil {
		t.Fatal("invalid submission accepted")
	}
}

func TestGradeSubmission(t *testing.T) {
	_, svc := newSvc()
	_ = svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"})
	_ = svc.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 5})
	_ = svc.AddQuestion(&model.Question{ID: "q2", ExamID: "e1", Score: 3})
	id, _ := svc.Submit(&model.Submission{ID: "s1", ExamID: "e1", Answers: []model.Answer{
		{QuestionID: "q1", Correct: true},
		{QuestionID: "q2", Correct: false},
	}})
	score, err := svc.GradeSubmission(id)
	if err != nil {
		t.Fatal(err)
	}
	if score != 5 {
		t.Fatalf("score=%v want 5", score)
	}
}

func TestGradeMissingWraps(t *testing.T) {
	_, svc := newSvc()
	if _, err := svc.GradeSubmission("nope"); !errors.Is(err, store.ErrSubmissionNotFound) {
		t.Fatalf("errors.Is=%v err=%v", errors.Is(err, store.ErrSubmissionNotFound), err)
	}
}

func TestSubmissionBatchesSorted(t *testing.T) {
	_, svc := newSvc()
	_ = svc.CreateExam(&model.Exam{ID: "e1", Title: "数学"})
	_, _ = svc.Submit(&model.Submission{ID: "s1", ExamID: "e1"})
	_, _ = svc.Submit(&model.Submission{ID: "s2", ExamID: "e1"})
	bs := svc.SubmissionBatches()
	if len(bs) == 0 || len(bs[0]) == 0 || bs[0][0].ID != "sub-1" {
		t.Fatalf("batches=%v", bs)
	}
}
