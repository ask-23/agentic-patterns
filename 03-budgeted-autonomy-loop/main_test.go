package main

import "testing"

func TestRun_CompletesWhenEvidenceArrivesWithinBudget(t *testing.T) {
	r := run("Prepare release handoff", Budget{
		MaxSteps:      5,
		MaxStaleSteps: 2,
	})

	if r.Escalate {
		t.Fatalf("expected COMPLETE, got ESCALATE: %+v", r)
	}
	if !r.Complete {
		t.Fatalf("expected Complete=true: %+v", r)
	}
	if r.Steps != 3 {
		t.Fatalf("steps = %d, want 3", r.Steps)
	}
	if r.EvidenceCount != 3 {
		t.Fatalf("evidence count = %d, want 3", r.EvidenceCount)
	}
}

func TestRun_EscalatesWhenProgressGoesStale(t *testing.T) {
	r := run("Investigate stale incident", Budget{
		MaxSteps:      5,
		MaxStaleSteps: 2,
	})

	if !r.Escalate {
		t.Fatalf("expected ESCALATE, got COMPLETE: %+v", r)
	}
	if r.Reason != "stale progress budget exhausted" {
		t.Fatalf("reason = %q, want stale progress budget exhausted", r.Reason)
	}
	if r.Steps != 2 {
		t.Fatalf("steps = %d, want 2", r.Steps)
	}
}

func TestRun_EscalatesWhenStepBudgetRunsOut(t *testing.T) {
	r := run("Prepare release handoff", Budget{
		MaxSteps:      2,
		MaxStaleSteps: 2,
	})

	if !r.Escalate {
		t.Fatalf("expected ESCALATE, got COMPLETE: %+v", r)
	}
	if r.Reason != "step budget exhausted" {
		t.Fatalf("reason = %q, want step budget exhausted", r.Reason)
	}
	if r.EvidenceCount != 2 {
		t.Fatalf("evidence count = %d, want 2", r.EvidenceCount)
	}
}

func TestRun_TaskIsCarriedOnResult(t *testing.T) {
	task := "Prepare release handoff"
	r := run(task, Budget{
		MaxSteps:      5,
		MaxStaleSteps: 2,
	})

	if r.Task != task {
		t.Fatalf("Task field not propagated: got %q, want %q", r.Task, task)
	}
}
