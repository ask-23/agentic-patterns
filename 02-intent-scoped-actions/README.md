# Intent Scoped Actions

This example demonstrates a bounded agentic workflow.

The workflow:

1. Accepts an operator intent.
2. Lets the agent propose a tool action.
3. Checks the action against the declared service scope.
4. Checks the action against safe verbs and intent-approved verbs.
5. Allows the action if it fits.
6. Escalates if the action is unsafe or outside the declared intent.

The implementation uses mocked behavior sans integrations for clarity.

## Run

```bash
go run ./main.go
```

## Expected output

`main` runs two scenarios so both branches are visible:

```text
ALLOW: verb="describe" service="cloudwatch" target="alarm history" reason="matches intent"
ESCALATE: verb="delete" service="iam" target="old admin role" reason="unsafe verb \"delete\""
```

The first task proposes a read-only CloudWatch action that matches the declared
intent.
The second task proposes a destructive IAM action and escalates before execution.

## Test

```bash
go test ./...
```

Covers the allow branch, destructive-action escalation, out-of-scope service
escalation, and case-insensitive intent matching.

## Production concern demonstrated

An autonomous step should not execute a proposed tool call merely because it
can. It should show that the action still matches the operator's intent. The
mechanism here is intentionally tiny — service matching plus safe verb checks —
so the boundary stays visible. A real version would check structured tool
schemas, resource patterns, approval tiers, and audit trails behind the same
decision point.
