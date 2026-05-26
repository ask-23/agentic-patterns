# Budgeted Autonomy Loop

This example demonstrates a bounded agentic workflow.

The workflow:

1. Accepts a task and a loop budget.
2. Lets the agent take autonomous steps.
3. Counts new evidence as progress.
4. Resets the stale counter when evidence appears.
5. Completes when enough evidence is gathered.
6. Escalates when the step budget or stale-progress budget is exhausted.

## Run

```bash
go run ./main.go
```

## Expected output

`main` runs two scenarios so both branches are visible:

```text
COMPLETE: steps=3 evidence=3 task="Prepare release handoff" reason="required evidence gathered"
ESCALATE: steps=2 evidence=0 task="Investigate stale incident" reason="stale progress budget exhausted"
```

The first task keeps producing evidence and completes inside the budget. The
second task keeps circling without new evidence and escalates after two stale
steps.

## Test

```bash
go test ./...
```

Covers completion inside budget, stale-progress escalation, total-step
escalation, and task propagation.

## Production concern demonstrated

An autonomous loop should have room to continue while it is making real
progress. It should not keep spending turns just because it can. The mechanism
here is intentionally tiny — a total step budget and a consecutive
stale-progress budget — so the boundary stays visible. A real version would
define evidence from tool results, state diffs, eval scores, or human-readable
artifacts behind the same loop gate.
