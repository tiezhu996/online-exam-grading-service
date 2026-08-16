package model

import "testing"

func TestValidAnswer(t *testing.T) {
	if !ValidAnswer(&Answer{QuestionID: "q1"}) {
		t.Fatal("valid answer rejected")
	}
	if ValidAnswer(nil) || ValidAnswer(&Answer{}) {
		t.Fatal("invalid answer accepted")
	}
}

func TestValidQuestion(t *testing.T) {
	if !ValidQuestion(&Question{ID: "q1", ExamID: "e1", Score: 5}) {
		t.Fatal("valid question rejected")
	}
	if ValidQuestion(nil) || ValidQuestion(&Question{ID: "q1"}) || ValidQuestion(&Question{ID: "", ExamID: "e1", Score: 5}) {
		t.Fatal("invalid question accepted")
	}
}

func TestValidSubmission(t *testing.T) {
	if !ValidSubmission(&Submission{ID: "s1", ExamID: "e1"}) {
		t.Fatal("valid submission rejected")
	}
	if ValidSubmission(nil) || ValidSubmission(&Submission{ID: "s1"}) || ValidSubmission(&Submission{ID: "", ExamID: "e1"}) {
		t.Fatal("invalid submission accepted")
	}
}

func TestScoreFor(t *testing.T) {
	qs := []*Question{{ID: "q1", Score: 2}, {ID: "q2", Score: 3}}
	ans := []Answer{{QuestionID: "q1", Correct: true}, {QuestionID: "q2", Correct: false}}
	if got := ScoreFor(qs, ans); got != 2 {
		t.Fatalf("score=%v want 2", got)
	}
}

func TestBuildSubmissionBatchesFresh(t *testing.T) {
	subs := []*Submission{{ID: "s1"}, {ID: "s2"}, {ID: "s3"}}
	bs := BuildSubmissionBatches(subs, 2)
	if len(bs) != 2 {
		t.Fatalf("batches=%d", len(bs))
	}
	bs[0][0] = &Submission{ID: "zz"}
	if subs[0].ID != "s1" {
		t.Fatal("mutating batch corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Checked: 1, Failed: 1}, Summary{Checked: 2, Graded: 3, Failed: 4})
	if got.Checked != 3 || got.Graded != 3 || got.Failed != 5 {
		t.Fatalf("merge=%+v", got)
	}
}
