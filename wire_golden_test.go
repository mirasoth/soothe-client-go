package soothe

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// goldenWireSHA256 pins the host wire fixture (tests/wire_fixtures/wire_golden.json).
// Regenerating the fixture on the host requires updating this pin in the same change.
const goldenWireSHA256 = "ba5a5a408fc2c259b768eae2915289de1820a931da37aad30fdd410d52b16627"

// loadWireGolden loads the pinned fixture and returns its named cases.
func loadWireGolden(t *testing.T) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "wire_golden.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != goldenWireSHA256 {
		t.Fatalf("wire_golden.json sha256 = %s, want %s", got, goldenWireSHA256)
	}
	var golden struct {
		Cases map[string]interface{} `json:"cases"`
	}
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	return golden.Cases
}

// TestWireGolden_AllCasesDecodeWithoutDrops decodes every golden case through
// the client codec (split → decode → expand) and asserts no frame errors or
// silent drops. Legacy flat payload shapes (error / connection_init / response)
// must still decode as envelopes without forcing nested-error semantics.
func TestWireGolden_AllCasesDecodeWithoutDrops(t *testing.T) {
	cases := loadWireGolden(t)
	names := make([]string, 0, len(cases))
	for name := range cases {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			wire, err := json.Marshal(cases[name])
			if err != nil {
				t.Fatalf("marshal case: %v", err)
			}
			frames := SplitSootheWirePayload(wire)
			if len(frames) == 0 {
				t.Fatal("no frames split from case")
			}
			for i, frame := range frames {
				msg, err := DecodeMessage(frame)
				if err != nil {
					t.Fatalf("frame %d: decode: %v", i, err)
				}
				if msg == nil {
					t.Fatalf("frame %d: decoded to nil (dropped)", i)
				}
				expanded := ExpandWireMessages(msg)
				if len(expanded) == 0 {
					t.Fatalf("frame %d: expanded to zero messages", i)
				}
				for j, m := range expanded {
					if m == nil {
						t.Fatalf("frame %d: expanded message %d is nil", i, j)
					}
				}
			}
		})
	}
}

// TestWireGolden_TypedShapes pins the decoded Go types for protocol-1 cases.
func TestWireGolden_TypedShapes(t *testing.T) {
	cases := loadWireGolden(t)

	for _, name := range []string{"connection_init", "response", "error"} {
		wire, err := json.Marshal(cases[name])
		if err != nil {
			t.Fatalf("%s: marshal: %v", name, err)
		}
		msg, err := DecodeMessage(wire)
		if err != nil {
			t.Fatalf("%s: decode: %v", name, err)
		}
		env, ok := msg.(Envelope)
		if !ok {
			t.Fatalf("%s: decoded %T, want Envelope", name, msg)
		}
		if env.Type != name {
			t.Errorf("%s: envelope type = %q", name, env.Type)
		}
		if env.ID == "" {
			t.Errorf("%s: envelope id is empty", name)
		}
	}

	// batch: bare JSON-array frame expands to one decoded envelope per element.
	batchWire, err := json.Marshal(cases["batch"])
	if err != nil {
		t.Fatalf("batch: marshal: %v", err)
	}
	msg, err := DecodeMessage(batchWire)
	if err != nil {
		t.Fatalf("batch: decode: %v", err)
	}
	expanded := ExpandWireMessages(msg)
	if len(expanded) != 2 {
		t.Fatalf("batch: expanded to %d messages, want 2", len(expanded))
	}
	for i, want := range []string{"ping", "pong"} {
		env, ok := expanded[i].(Envelope)
		if !ok {
			t.Fatalf("batch[%d]: decoded %T, want Envelope", i, expanded[i])
		}
		if env.Type != want {
			t.Errorf("batch[%d]: type = %q, want %q", i, env.Type, want)
		}
	}

	// stream_tool_update: custom payload type survives as a generic map.
	toolWire, err := json.Marshal(cases["stream_tool_update"])
	if err != nil {
		t.Fatalf("stream_tool_update: marshal: %v", err)
	}
	toolMsg, err := DecodeMessage(toolWire)
	if err != nil {
		t.Fatalf("stream_tool_update: decode: %v", err)
	}
	toolMap, ok := toolMsg.(map[string]interface{})
	if !ok {
		t.Fatalf("stream_tool_update: decoded %T, want map", toolMsg)
	}
	if toolMap["type"] != "soothe.stream.tool_call.update" {
		t.Errorf("stream_tool_update: type = %v", toolMap["type"])
	}
}
