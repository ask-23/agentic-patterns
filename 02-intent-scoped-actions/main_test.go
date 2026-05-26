package main

import "testing"

func TestRun_InScopeSafeActionAllows(t *testing.T) {
	intent := Intent{
		Service:      "cloudwatch",
		AllowedVerbs: []string{"describe", "list"},
	}

	d := run(intent, "Check the CloudWatch alarm state")
	if d.Escalate {
		t.Fatalf("expected ALLOW, got ESCALATE: %+v", d)
	}
	if !d.Allow {
		t.Fatalf("expected Allow=true: %+v", d)
	}
}

func TestRun_DestructiveActionEscalates(t *testing.T) {
	intent := Intent{
		Service:      "cloudwatch",
		AllowedVerbs: []string{"describe", "list"},
	}

	d := run(intent, "Clean up IAM access")
	if !d.Escalate {
		t.Fatalf("expected ESCALATE, got ALLOW: %+v", d)
	}
	if d.Allow {
		t.Fatalf("expected Allow=false: %+v", d)
	}
}

func TestAuthorize_RejectsUnsafeVerbEvenWhenIntentListsIt(t *testing.T) {
	intent := Intent{
		Service:      "cloudwatch",
		AllowedVerbs: []string{"delete"},
	}
	action := Action{
		Verb:    "delete",
		Service: "cloudwatch",
		Target:  "alarm history",
	}

	d := authorize(intent, action)
	if !d.Escalate {
		t.Fatalf("expected ESCALATE for unsafe verb, got ALLOW: %+v", d)
	}
}

func TestAuthorize_RejectsOutOfScopeService(t *testing.T) {
	intent := Intent{
		Service:      "cloudwatch",
		AllowedVerbs: []string{"describe"},
	}
	action := Action{
		Verb:    "describe",
		Service: "iam",
		Target:  "admin role",
	}

	d := authorize(intent, action)
	if !d.Escalate {
		t.Fatalf("expected ESCALATE for out-of-scope service, got ALLOW: %+v", d)
	}
}

func TestAuthorize_MatchesCaseInsensitively(t *testing.T) {
	intent := Intent{
		Service:      "CloudWatch",
		AllowedVerbs: []string{"Describe"},
	}
	action := Action{
		Verb:    "describe",
		Service: "cloudwatch",
		Target:  "alarm history",
	}

	d := authorize(intent, action)
	if d.Escalate {
		t.Fatalf("expected ALLOW for case-insensitive match, got ESCALATE: %+v", d)
	}
}
