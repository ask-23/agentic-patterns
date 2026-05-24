# Research With Eval Gate

This example demonstrates a bounded agentic workflow.

The workflow:

1. Accepts a research task.
2. Produces a synthesis shaped by the task.
3. Scores the synthesis against a four-term eval rubric.
4. Returns the answer if the score passes (>= 3 of 4 terms).
5. Escalates if the score fails.

The implementation uses mocked behavior. That is intentional. The point is not model integration. The point is the control boundary.

## Run

```bash
go run ./main.go
```

## Expected output

`main` runs two scenarios so both branches are visible:

```text
PASS: score=4 task="What makes an AI feature production-grade?" synthesis="Production AI requires evaluation, observability, rollback paths, and clear autonomy boundaries."
ESCALATE: score=1 task="How should we ship an agent?" synthesis="Agents should be evaluated and should have observability."
```

The first task produces a synthesis that mentions all four required terms (`evaluation`, `observability`, `rollback`, `boundaries`); the gate passes. The second task produces a deliberately weaker synthesis that misses two of the four; the gate escalates.

## Test

```bash
go test ./...
```

Covers both branches plus a generic-task escalation case.

## Production concern demonstrated

An autonomous step should not publish output merely because it completed. It should pass a quality gate or escalate. The mechanism here is intentionally tiny — substring scoring against a fixed rubric — so the boundary stays visible. A real version would substitute a rubric-based or LLM-judge evaluator behind the same gate.
