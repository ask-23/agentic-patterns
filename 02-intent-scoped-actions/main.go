package main

import (
	"fmt"
	"strings"
)

// Intent is the operator boundary for a short agent run.
//
// The agent can still choose an action, but the action must stay inside this
// declared service scope and verb set before execution is allowed.
type Intent struct {
	Service      string
	AllowedVerbs []string
}

// Action is the agent's proposed tool call.
//
// A production version would carry tool names, parameters, trace IDs, and
// approval metadata. This example keeps only the fields needed to show the
// boundary.
type Action struct {
	Verb    string
	Service string
	Target  string
	Reason  string
}

// Decision is the bounded result of checking a proposed action.
type Decision struct {
	Task     string
	Action   Action
	Allow    bool
	Escalate bool
	Reason   string
}

var safeVerbs = map[string]bool{
	"describe": true,
	"get":      true,
	"inspect":  true,
	"list":     true,
	"read":     true,
}

// propose is the agent's "choose a tool action" step.
//
// Deterministic by design: the point is the intent boundary, not model
// behavior or an actual cloud integration.
func propose(task string) Action {
	lower := strings.ToLower(task)
	switch {
	case strings.Contains(lower, "cloudwatch") || strings.Contains(lower, "alarm"):
		return Action{
			Verb:    "describe",
			Service: "cloudwatch",
			Target:  "alarm history",
			Reason:  "inspect requested alarm state",
		}
	case strings.Contains(lower, "iam") || strings.Contains(lower, "clean up"):
		return Action{
			Verb:    "delete",
			Service: "iam",
			Target:  "old admin role",
			Reason:  "remove access that looks stale",
		}
	default:
		return Action{
			Verb:    "describe",
			Service: "unknown",
			Target:  "task context",
			Reason:  "gather basic context",
		}
	}
}

// authorize decides whether the proposed action may run.
func authorize(intent Intent, action Action) Decision {
	verb := normalize(action.Verb)
	service := normalize(action.Service)
	intentService := normalize(intent.Service)

	switch {
	case !safeVerbs[verb]:
		return Decision{
			Action:   action,
			Escalate: true,
			Reason:   fmt.Sprintf("unsafe verb %q", action.Verb),
		}
	case service != intentService:
		return Decision{
			Action:   action,
			Escalate: true,
			Reason:   fmt.Sprintf("service %q outside intent service %q", action.Service, intent.Service),
		}
	case !verbAllowed(verb, intent.AllowedVerbs):
		return Decision{
			Action:   action,
			Escalate: true,
			Reason:   fmt.Sprintf("verb %q outside intent verbs", action.Verb),
		}
	default:
		return Decision{
			Action: action,
			Allow:  true,
			Reason: "matches intent",
		}
	}
}

// run is the bounded workflow:
//
//  1. accept operator intent and a task
//  2. let the agent propose an action
//  3. authorize the action against intent, service scope, and safe verbs
//  4. decide: ALLOW or ESCALATE
func run(intent Intent, task string) Decision {
	d := authorize(intent, propose(task))
	d.Task = task
	return d
}

func verbAllowed(verb string, allowed []string) bool {
	for _, candidate := range allowed {
		if verb == normalize(candidate) {
			return true
		}
	}
	return false
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// print emits a single line per decision. Format is stable so it is grep-able
// from CI logs: ALLOW or ESCALATE leads, then the proposed action and reason.
func print(d Decision) {
	verdict := "ALLOW"
	if d.Escalate {
		verdict = "ESCALATE"
	}
	fmt.Printf("%s: verb=%q service=%q target=%q reason=%q\n",
		verdict, d.Action.Verb, d.Action.Service, d.Action.Target, d.Reason)
}

func main() {
	intent := Intent{
		Service:      "cloudwatch",
		AllowedVerbs: []string{"describe", "list"},
	}
	scenarios := []string{
		"Check the CloudWatch alarm state",
		"Clean up IAM access",
	}
	for _, task := range scenarios {
		print(run(intent, task))
	}
}
