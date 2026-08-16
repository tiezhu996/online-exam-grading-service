package store

import (
	"testing"

	"onlineexam/internal/model"
)

func TestPutGetExam(t *testing.T) {
	s := New()
	if err := s.PutExam(&model.Exam{ID: "e1", Title: "数学"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutExam(&model.Exam{ID: "e1", Title: "x"}); err != ErrExamExists {
		t.Fatalf("dup err=%v", err)
	}
	if _, err := s.GetExam("e1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetExam("nope"); err != ErrExamNotFound {
		t.Fatalf("missing err=%v", err)
	}
}

func TestExamIDsFresh(t *testing.T) {
	s := New()
	_ = s.PutExam(&model.Exam{ID: "e2", Title: "x"})
	_ = s.PutExam(&model.Exam{ID: "e1", Title: "x"})
	ids := s.ExamIDs()
	ids[0] = "zz"
	if s.ExamIDs()[0] != "e2" {
		t.Fatal("ExamIDs aliased")
	}
}

func TestAddListQuestions(t *testing.T) {
	s := New()
	_ = s.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 5})
	if q, err := s.GetQuestion("q1"); err != nil || q.Score != 5 {
		t.Fatalf("q=%v err=%v", q, err)
	}
	if _, err := s.GetQuestion("nope"); err != ErrQuestionNotFound {
		t.Fatalf("missing err=%v", err)
	}
	if len(s.ListQuestions("e1")) != 1 {
		t.Fatal("ListQuestions len")
	}
}

func TestQuestionIDsFresh(t *testing.T) {
	s := New()
	_ = s.AddQuestion(&model.Question{ID: "q2", ExamID: "e1", Score: 1})
	_ = s.AddQuestion(&model.Question{ID: "q1", ExamID: "e1", Score: 1})
	ids := s.QuestionIDs()
	ids[0] = "zz"
	if s.QuestionIDs()[0] != "q2" {
		t.Fatal("QuestionIDs aliased")
	}
}

func TestAddGetSubmission(t *testing.T) {
	s := New()
	id, err := s.AddSubmission(&model.Submission{ID: "s1", ExamID: "e1", Status: model.StatusPending})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty submission id")
	}
	if _, err := s.GetSubmission(id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetSubmission("nope"); err != ErrSubmissionNotFound {
		t.Fatalf("missing err=%v", err)
	}
}

func TestSubmissionIDsFresh(t *testing.T) {
	s := New()
	_, _ = s.AddSubmission(&model.Submission{ID: "s1", ExamID: "e1", Status: model.StatusPending})
	ids := s.SubmissionIDs()
	ids[0] = "zz"
	if s.SubmissionIDs()[0] != "sub-1" {
		t.Fatal("SubmissionIDs aliased")
	}
}

func TestMarkGraded(t *testing.T) {
	s := New()
	id, _ := s.AddSubmission(&model.Submission{ID: "s1", ExamID: "e1", Status: model.StatusPending})
	if err := s.MarkGraded(id); err != nil {
		t.Fatal(err)
	}
	sub, _ := s.GetSubmission(id)
	if sub.Status != model.StatusGraded {
		t.Fatalf("status=%v", sub.Status)
	}
	if err := s.MarkGraded("nope"); err != ErrSubmissionNotFound {
		t.Fatalf("missing err=%v", err)
	}
}
