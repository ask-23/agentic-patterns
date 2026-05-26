package main

import (
	"fmt"
	"strings"
)

// Budget is the small operating boundary for an autonomous loop.
//
// MaxSteps caps total work. MaxStaleSteps caps consecutive iterations that do
// not produce new evidence.
type Budget struct {
	MaxSteps      int
	MaxStaleSteps int
}

// Step is one agent iteration.
//
// Evidence is the signal that the loop is still making useful progress. A
// production version would carry tool traces, citations, or state diffs.
type Step struct {
	Summary  string
	Evidence string
}

// Result is the bounded output of a budgeted autonomy run.
type Result struct {
	Task          string
	Steps         int
	EvidenceCount int
	Complete      bool
	Escalate      bool
	Reason        string
}

// work is the agent's "do the next useful thing" step.
//
// Deterministic by design: this example demonstrates the loop boundary, not
// model behavior or real tool calls.
func work(task string, step int) Step {
	lower := strings.ToLower(task)
	switch {
	case strings.Contains(lower, "release"):
		steps := []Step{
			{Summary: "run checks", Evidence: "tests passed"},
			{Summary: "check rollback", Evidence: "rollback path noted"},
			{Summary: "prepare handoff", Evidence: "handoff ready"},
		}
		if step < len(steps) {
			return steps[step]
		}
		return Step{Summary: "release handoff already covered"}
	case strings.Contains(lower, "stale"):
		return Step{Summary: "rechecked the same dashboard"}
	default:
		if step == 0 {
			return Step{Summary: "capture initial context", Evidence: "initial context"}
		}
		return Step{Summary: "re-read task context"}
	}
}

// run is the bounded workflow:
//
//  1. accept a task and loop budget
//  2. take autonomous steps while budget remains
//  3. count new evidence as progress
//  4. decide: COMPLETE once enough evidence exists, otherwise ESCALATE when
//     stale-progress or total-step budgets are exhausted
func run(task string, budget Budget) Result {
	const requiredEvidence = 3

	result := Result{Task: task}
	staleSteps := 0
	maxStaleSteps := budget.MaxStaleSteps
	if maxStaleSteps <= 0 {
		maxStaleSteps = 1
	}

	for result.Steps < budget.MaxSteps {
		step := work(task, result.Steps)
		result.Steps++

		if strings.TrimSpace(step.Evidence) == "" {
			staleSteps++
		} else {
			result.EvidenceCount++
			staleSteps = 0
		}

		if result.EvidenceCount >= requiredEvidence {
			result.Complete = true
			result.Reason = "required evidence gathered"
			return result
		}

		if staleSteps >= maxStaleSteps {
			result.Escalate = true
			result.Reason = "stale progress budget exhausted"
			return result
		}
	}

	result.Escalate = true
	result.Reason = "step budget exhausted"
	return result
}

// print emits a single line per result. Format is stable so it is grep-able
// from CI logs: COMPLETE or ESCALATE leads, then budget counters and reason.
func print(r Result) {
	verdict := "COMPLETE"
	if r.Escalate {
		verdict = "ESCALATE"
	}
	fmt.Printf("%s: steps=%d evidence=%d task=%q reason=%q\n",
		verdict, r.Steps, r.EvidenceCount, r.Task, r.Reason)
}

func main() {
	budget := Budget{
		MaxSteps:      5,
		MaxStaleSteps: 2,
	}
	scenarios := []string{
		"Prepare release handoff",
		"Investigate stale incident",
	}
	for _, task := range scenarios {
		print(run(task, budget))
	}
}
