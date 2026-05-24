package main

import (
	"fmt"
	"strings"
)

// Result is the bounded output of a single research-with-eval-gate run.
//
// The shape is deliberately small. A real implementation would carry trace
// IDs, retrieval context, tool calls, and evaluator metadata; the production
// concern this example demonstrates is the boundary, not the integration.
type Result struct {
	Task      string
	Synthesis string
	Score     int
	Escalate  bool
}

// requiredTerms is the eval rubric. A synthesis must mention each of these to
// pass the gate. In production this would be a real evaluator (rubric-based,
// LLM-as-judge, golden-set similarity, etc.); the shape is the same.
var requiredTerms = []string{"evaluation", "observability", "rollback", "boundaries"}

// research is the agent's "do work" step. It returns content shaped by the
// task input, so the eval gate can produce visibly different outcomes.
//
// Deterministic by design: this example exists to demonstrate the control
// boundary, not model behavior. Production code would call an LLM here.
func research(task string) string {
	lower := strings.ToLower(task)
	switch {
	case strings.Contains(lower, "production-grade"):
		return "Production AI requires evaluation, observability, rollback paths, and clear autonomy boundaries."
	case strings.Contains(lower, "agent"):
		// Intentionally weak synthesis: misses two required terms.
		return "Agents should be evaluated and should have observability."
	default:
		// Generic minimal synthesis: misses most required terms.
		return "AI features should be tested before launch."
	}
}

// evaluate scores a synthesis against the rubric. One point per required term.
func evaluate(synthesis string) int {
	score := 0
	lower := strings.ToLower(synthesis)
	for _, term := range requiredTerms {
		if strings.Contains(lower, term) {
			score++
		}
	}
	return score
}

// run is the bounded workflow:
//
//  1. accept a task
//  2. produce a synthesis
//  3. score it
//  4. decide: PASS if score >= threshold, otherwise ESCALATE
func run(task string) Result {
	const passThreshold = 3
	s := research(task)
	score := evaluate(s)
	return Result{
		Task:      task,
		Synthesis: s,
		Score:     score,
		Escalate:  score < passThreshold,
	}
}

// print emits a single line per result. Format is stable so it is grep-able
// from CI logs: PASS or ESCALATE leads, then the score, then the synthesis.
func print(r Result) {
	verdict := "PASS"
	if r.Escalate {
		verdict = "ESCALATE"
	}
	fmt.Printf("%s: score=%d task=%q synthesis=%q\n",
		verdict, r.Score, r.Task, r.Synthesis)
}

func main() {
	// Two scenarios run on every invocation so the boundary is visible without
	// requiring CLI arguments. A strong task passes the gate; a weak task
	// triggers escalation.
	scenarios := []string{
		"What makes an AI feature production-grade?",
		"How should we ship an agent?",
	}
	for _, task := range scenarios {
		print(run(task))
	}
}
