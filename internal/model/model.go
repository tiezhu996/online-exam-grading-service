package model

import "sort"

type Exam struct {
	ID          string
	Title       string
	DurationMin int
	Status      string
}

type Question struct {
	ID     string
	ExamID string
	Text   string
	Score  float64
}

type Answer struct {
	QuestionID string
	Correct    bool
}

type Submission struct {
	ID      string
	ExamID  string
	Answers []Answer
	Score   float64
	Status  string
}

type Summary struct {
	Checked int
	Graded  int
	Failed  int
}

const (
	StatusDraft   = "draft"
	StatusPending = "pending"
	StatusGraded  = "graded"
)

func ValidAnswer(a *Answer) bool {
	return a != nil && a.QuestionID != ""
}

func ValidQuestion(q *Question) bool {
	return q != nil && q.ID != "" && q.ExamID != "" && q.Score >= 0
}

func ValidSubmission(s *Submission) bool {
	return s != nil && s.ID != "" && s.ExamID != ""
}

func SortQuestions(qs []*Question) []*Question {
	sort.SliceStable(qs, func(i, j int) bool { return qs[i].ID < qs[j].ID })
	return qs
}

func SortSubmissions(subs []*Submission) []*Submission {
	sort.SliceStable(subs, func(i, j int) bool { return subs[i].ID < subs[j].ID })
	return subs
}

func BuildSubmissionBatches(subs []*Submission, size int) [][]*Submission {
	if size <= 0 {
		size = 1
	}
	out := make([][]*Submission, 0, (len(subs)+size-1)/size)
	for i := 0; i < len(subs); i += size {
		end := i + size
		if end > len(subs) {
			end = len(subs)
		}
		b := make([]*Submission, end-i)
		copy(b, subs[i:end])
		out = append(out, b)
	}
	return out
}

func MergeSummary(dst, src Summary) Summary {
	dst.Checked += src.Checked
	dst.Graded += src.Graded
	return dst
}

func ScoreFor(qs []*Question, answers []Answer) float64 {
	idx := map[string]float64{}
	for _, q := range qs {
		idx[q.ID] = q.Score
	}
	var total float64
	for _, a := range answers {
		if a.Correct {
			total += idx[a.QuestionID]
		}
	}
	return total
}
