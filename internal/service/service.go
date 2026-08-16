package service

import (
	"errors"
	"fmt"

	"onlineexam/internal/config"
	"onlineexam/internal/model"
	"onlineexam/internal/store"
)

type Service struct {
	store     *store.Store
	batchSize int
}

func New(s *store.Store, cfg *config.Config) *Service {
	b := cfg.BatchSize
	if b <= 0 {
		b = 1
	}
	return &Service{store: s, batchSize: b}
}

func (svc *Service) CreateExam(e *model.Exam) error {
	if e == nil || e.ID == "" || e.Title == "" {
		return errors.New("invalid exam")
	}
	if err := svc.store.PutExam(e); err != nil {
		return fmt.Errorf("create exam %s: %w", e.ID, err)
	}
	return nil
}

func (svc *Service) AddQuestion(q *model.Question) error {
	if !model.ValidQuestion(q) {
		return errors.New("invalid question")
	}
	if err := svc.store.AddQuestion(q); err != nil {
		return fmt.Errorf("add question %s: %w", q.ID, err)
	}
	return nil
}

func (svc *Service) Submit(sub *model.Submission) (string, error) {
	if !model.ValidSubmission(sub) {
		return "", errors.New("invalid submission")
	}
	id, err := svc.store.AddSubmission(sub)
	if err != nil {
		return "", fmt.Errorf("submit %s: %w", sub.ID, err)
	}
	return id, nil
}

func (svc *Service) ListQuestions(examID string) []*model.Question {
	return svc.store.ListQuestions(examID)
}

func (svc *Service) ListSubmissions(examID string) []*model.Submission {
	return svc.store.ListSubmissions(examID)
}

func (svc *Service) GradeSubmission(id string) (float64, error) {
	sub, err := svc.store.GetSubmission(id)
	if err != nil {
		return 0, fmt.Errorf("grade %s: %w", id, err)
	}
	qs := svc.store.ListQuestions(sub.ExamID)
	sub.Score = model.ScoreFor(qs, sub.Answers)
	if err := svc.store.MarkGraded(id); err != nil {
		return 0, fmt.Errorf("mark graded %s: %w", id, err)
	}
	return sub.Score, nil
}

func (svc *Service) SubmissionBatches() [][]*model.Submission {
	subs := svc.store.AllSubmissions()
	pending := make([]*model.Submission, 0, len(subs))
	for _, s := range subs {
		if s.Status != model.StatusGraded {
			pending = append(pending, s)
		}
	}
	model.SortSubmissions(pending)
	out := make([][]*model.Submission, 0)
	for i := 0; i < len(pending); i += svc.batchSize {
		end := i + svc.batchSize
		if end > len(pending) {
			end = len(pending)
		}
		out = append(out, pending[i:end])
	}
	return out
}
