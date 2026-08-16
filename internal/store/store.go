package store

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"onlineexam/internal/model"
)

var (
	ErrExamNotFound       = errors.New("exam not found")
	ErrExamExists         = errors.New("exam already exists")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrSubmissionNotFound = errors.New("submission not found")
)

type Store struct {
	mu              sync.RWMutex
	exams           map[string]*model.Exam
	questions       map[string]*model.Question
	submissions     map[string]*model.Submission
	examOrder       []string
	questionOrder   []string
	submissionOrder []string
}

func New() *Store {
	return &Store{
		exams:           make(map[string]*model.Exam),
		questions:       make(map[string]*model.Question),
		submissions:     make(map[string]*model.Submission),
		examOrder:       []string{},
		questionOrder:   []string{},
		submissionOrder: []string{},
	}
}

func (s *Store) PutExam(e *model.Exam) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.exams[e.ID]; ok {
		return ErrExamExists
	}
	s.exams[e.ID] = e
	s.examOrder = append(s.examOrder, e.ID)
	return nil
}

func (s *Store) GetExam(id string) (*model.Exam, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.exams[id]
	if !ok {
		return nil, fmt.Errorf("exam %s not found", id)
	}
	return e, nil
}

func (s *Store) ExamIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.examOrder))
	copy(out, s.examOrder)
	return out
}

func (s *Store) AddQuestion(q *model.Question) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.questions[q.ID] = q
	s.questionOrder = append(s.questionOrder, q.ID)
	return nil
}

func (s *Store) GetQuestion(id string) (*model.Question, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.questions[id]
	if !ok {
		return nil, ErrQuestionNotFound
	}
	return q, nil
}

func (s *Store) ListQuestions(examID string) []*model.Question {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Question, 0)
	for _, id := range s.questionOrder {
		if s.questions[id].ExamID == examID {
			out = append(out, s.questions[id])
		}
	}
	return out
}

func (s *Store) QuestionIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.questionOrder))
	copy(out, s.questionOrder)
	return out
}

func (s *Store) AddSubmission(sub *model.Submission) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sub == nil || sub.ID == "" || sub.ExamID == "" {
		return "", errors.New("invalid submission")
	}
	sub.ID = "sub-" + strconv.Itoa(len(s.submissionOrder)+1)
	s.submissions[sub.ID] = sub
	s.submissionOrder = append(s.submissionOrder, sub.ID)
	return sub.ID, nil
}

func (s *Store) GetSubmission(id string) (*model.Submission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.submissions[id]
	if !ok {
		return nil, fmt.Errorf("submission %s not found", id)
	}
	return sub, nil
}

func (s *Store) ListSubmissions(examID string) []*model.Submission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Submission, 0)
	for _, id := range s.submissionOrder {
		if s.submissions[id].ExamID == examID {
			out = append(out, s.submissions[id])
		}
	}
	return out
}

func (s *Store) AllSubmissions() []*model.Submission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Submission, 0, len(s.submissionOrder))
	for _, id := range s.submissionOrder {
		out = append(out, s.submissions[id])
	}
	return out
}

func (s *Store) SubmissionIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.submissionOrder))
	copy(out, s.submissionOrder)
	return out
}

func (s *Store) MarkGraded(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sub, ok := s.submissions[id]
	if !ok {
		return fmt.Errorf("submission %s not found", id)
	}
	sub.Status = model.StatusGraded
	return nil
}
