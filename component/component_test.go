package component

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTextComponent(t *testing.T) {
	component := Text("Ciao").Color(Green).Bold().Append(Text("!"))

	assertJSON(t, component, `{"text":"Ciao","extra":[{"text":"!"}],"color":"green","bold":true}`)
}

func TestComponentVariants(t *testing.T) {
	assertJSON(t, Translatable("chat.type.text", Text("Alex")), `{"translate":"chat.type.text","with":[{"text":"Alex"}]}`)
	assertJSON(t, Score("Alex", "points"), `{"score":{"name":"Alex","objective":"points"}}`)
	assertJSON(t, Selector("@p").Font("minecraft:default"), `{"selector":"@p","font":"minecraft:default"}`)
	assertJSON(t, Keybind("key.jump").Click(RunCommand, "/jump"), `{"keybind":"key.jump","clickEvent":{"action":"run_command","value":"/jump"}}`)
}

func TestHoverEvent(t *testing.T) {
	component := Text("hover").Hover(ShowText, Text("details"))

	assertJSON(t, component, `{"text":"hover","hoverEvent":{"action":"show_text","contents":{"text":"details"}}}`)
}

func TestExtendedComponents(t *testing.T) {
	assertJSON(t, Translatable("chat.type.text").Fallback("fallback"), `{"translate":"chat.type.text","fallback":"fallback"}`)
	assertJSON(t, ScoreValue("Alex", "points", "10"), `{"score":{"name":"Alex","objective":"points","value":"10"}}`)
	assertJSON(t, Selector("@a").Separator(Text(", ")), `{"selector":"@a","separator":{"text":", "}}`)
	assertJSON(t, EntityNBT("CustomName", "@s", true), `{"nbt":"CustomName","entity":"@s","interpret":true}`)
	assertJSON(t, BlockNBT("Items", "~ ~ ~", false), `{"nbt":"Items","block":"~ ~ ~","interpret":false}`)
}

func TestAudienceAndFormats(t *testing.T) {
	audience := recordingAudience{}
	if err := SendMessage([]Audience{audience}, Text("hello")); err != nil {
		t.Fatalf("send message: %v", err)
	}
	if Plain(ParseLegacy("§aHello §lworld")) != "Hello world" {
		t.Fatal("legacy text was not parsed")
	}
	parsed := ParseMiniMessage("<red>Hello <bold>world</bold></red>")
	if Plain(parsed) != "Hello world" {
		t.Fatalf("unexpected MiniMessage text: %q", Plain(parsed))
	}
	if got := Plain(ParsePlain("plain")); got != "plain" {
		t.Fatalf("unexpected plain text: %q", got)
	}
}

func TestMiniMessageEventsRGBAndGradient(t *testing.T) {
	component := ParseMiniMessage("<click:run_command:/spawn><hover:show_text:'Info'><#ff0000>Go</#ff0000></hover></click>")
	encoded, err := Marshal(component)
	if err != nil {
		t.Fatalf("marshal MiniMessage: %v", err)
	}
	if !strings.Contains(string(encoded), `"clickEvent"`) || !strings.Contains(string(encoded), `"hoverEvent"`) {
		t.Fatalf("events missing from JSON: %s", encoded)
	}
	gradient := ParseMiniMessage("<gradient:#ff0000:#0000ff>AB</gradient>")
	gradientJSON, err := Marshal(gradient)
	if err != nil {
		t.Fatalf("marshal gradient: %v", err)
	}
	if !strings.Contains(string(gradientJSON), "#ff0000") || !strings.Contains(string(gradientJSON), "#0000ff") {
		t.Fatalf("gradient colors missing from JSON: %s", gradientJSON)
	}
}

type recordingAudience struct{}

func (recordingAudience) SendMessage(Component) error { return nil }

func assertJSON(t *testing.T, component Component, expected string) {
	t.Helper()

	actual, err := json.Marshal(component)
	if err != nil {
		t.Fatalf("marshal component: %v", err)
	}

	var actualValue any
	var expectedValue any
	if err := json.Unmarshal(actual, &actualValue); err != nil {
		t.Fatalf("decode actual JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("decode expected JSON: %v", err)
	}
	if !equalJSON(actualValue, expectedValue) {
		t.Fatalf("JSON mismatch\nactual:   %s\nexpected: %s", actual, expected)
	}
}

func equalJSON(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) == string(rightJSON)
}
