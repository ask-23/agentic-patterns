package main

import "testing"

func TestRun_StrongTaskPassesGate(t *testing.T) {
	r := run("What makes an AI feature production-grade?")
	if r.Escalate {
		t.Fatalf("expected PASS, got ESCALATE: %+v", r)
	}
	if r.Score < 3 {
		t.Errorf("expected score >= 3, got %d", r.Score)
	}
}

func TestRun_WeakTaskEscalates(t *testing.T) {
	r := run("How should we ship an agent?")
	if !r.Escalate {
		t.Fatalf("expected ESCALATE, got PASS: %+v", r)
	}
	if r.Score >= 3 {
		t.Errorf("expected score < 3, got %d", r.Score)
	}
}

func TestRun_GenericTaskEscalates(t *testing.T) {
	r := run("something else entirely")
	if !r.Escalate {
		t.Fatalf("expected ESCALATE for generic task, got PASS: %+v", r)
	}
}

func TestRun_TaskIsCarriedOnResult(t *testing.T) {
	task := "What makes an AI feature production-grade?"
	r := run(task)
	if r.Task != task {
		t.Fatalf("Task field not propagated: got %q, want %q", r.Task, task)
	}
}

func TestEvaluate_CountsRequiredTerms(t *testing.T) {
	s := "evaluation observability rollback boundaries"
	if got := evaluate(s); got != 4 {
		t.Fatalf("score = %d, want 4", got)
	}
}
