package appkit

import (
	"fmt"
	"strings"

	soothe "github.com/mirasoth/soothe-client-go"
)

const maxThinkingStepRunes = 280

// Default thinking-step event allowlist. Overridable via
// ClassifierConfig.ThinkingStepEvents.
var defaultThinkingStepEvents = map[string]bool{
	soothe.EventStrangeLoopStepStarted:   true,
	soothe.EventStrangeLoopStepQueued:    true,
	soothe.EventStrangeLoopStepCompleted: true,
	soothe.EventReplayComplete:           true,
	soothe.EventPlanCreated:              true,
	soothe.EventGoalCreated:              true,
	soothe.EventToolStarted:              true,
}

// extractThinkingStep maps an allowlisted progress event to one structured UI
// line. Free-form streams (tokens, reports, reasoning) are excluded. The
// allowlist is configurable via ClassifierConfig.ThinkingStepEvents.
func (cl *EventClassifier) extractThinkingStep(eventType string, data map[string]interface{}) (string, bool) {
	if eventType == "" || data == nil {
		return "", false
	}
	eventType = strings.TrimSpace(eventType)

	allow := cl.cfg.ThinkingStepEvents
	if allow == nil {
		allow = defaultThinkingStepEvents
	}
	if !allow[eventType] {
		return "", false
	}

	var line string
	switch eventType {
	case soothe.EventStrangeLoopStepStarted:
		line = formatStepLine(data, "")
	case soothe.EventStrangeLoopStepQueued:
		line = formatStepLine(data, "queued")
	case soothe.EventStrangeLoopStepCompleted:
		line = formatStepCompletionLine(data)
	case soothe.EventReplayComplete:
		line = "History replayed"
	case soothe.EventPlanCreated:
		if g := strField(data, "goal"); g != "" {
			line = "Goal: " + g
		}
	case soothe.EventGoalCreated:
		if g := strField(data, "friendly_message", "description"); g != "" {
			line = "Goal: " + g
		}
	case soothe.EventToolStarted:
		if name := strField(data, "tool_name", "name"); name != "" {
			line = "Tool: " + name
		}
	default:
		return "", false
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return "", false
	}
	if len([]rune(line)) > maxThinkingStepRunes {
		runes := []rune(line)
		line = string(runes[:maxThinkingStepRunes]) + "…"
	}
	return line, true
}

func formatStepLine(data map[string]interface{}, suffix string) string {
	stepID := strField(data, "step_id")
	desc := strField(data, "description")
	switch {
	case stepID != "" && suffix != "":
		return fmt.Sprintf("Step %s: %s", stepID, suffix)
	case stepID != "" && desc != "":
		return fmt.Sprintf("Step %s: %s", stepID, desc)
	case stepID != "":
		return fmt.Sprintf("Step %s", stepID)
	case desc != "" && suffix != "":
		return fmt.Sprintf("Step: %s", suffix)
	case desc != "":
		return fmt.Sprintf("Step: %s", desc)
	case suffix != "":
		return "Step: " + suffix
	default:
		return ""
	}
}

func formatStepCompletionLine(data map[string]interface{}) string {
	stepID := strField(data, "step_id")
	summary := strField(data, "summary")
	if failed, ok := data["success"].(bool); ok && !failed {
		switch {
		case stepID != "" && summary != "":
			return fmt.Sprintf("Step %s failed: %s", stepID, summary)
		case stepID != "":
			return fmt.Sprintf("Step %s failed", stepID)
		case summary != "":
			return "Step failed: " + summary
		}
		return ""
	}
	switch {
	case stepID != "" && summary != "":
		return fmt.Sprintf("Step %s done: %s", stepID, summary)
	case stepID != "":
		return fmt.Sprintf("Step %s done", stepID)
	case summary != "":
		return "Step: " + summary
	default:
		return ""
	}
}

func strField(data map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := data[key].(string); ok {
			if s := strings.TrimSpace(v); s != "" {
				return s
			}
		}
	}
	return ""
}

func thinkingStepResult(step string) ChatEventResult {
	return ChatEventResult{
		ThinkingStep: strings.TrimSpace(step),
		Terminal:     ChatEventContinue,
	}
}
