package typespeed

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func unicodeTestModel(prompt string) Model {
	return Model{
		PromptStr:        prompt,
		PromptSlice:      strings.Split(prompt, " "),
		PromptUnderlines: strings.Repeat(" ", utf8.RuneCountInString(prompt)+1),
		State: &State{
			SeenIdxSet: make(map[int]int),
		},
	}
}

func TestIsValidCharAcceptsUnicodeRunes(t *testing.T) {
	for _, input := range []string{"你", "好", "é", "🙂", " "} {
		if !isValidChar(input) {
			t.Fatalf("expected %q to be accepted as a printable character", input)
		}
	}

	for _, input := range []string{"", "ab", "\n"} {
		if isValidChar(input) {
			t.Fatalf("expected %q to be rejected", input)
		}
	}
}

func TestTypeCharAdvancesByRuneForChineseInput(t *testing.T) {
	m := unicodeTestModel("你好 世界")

	typeChar(&m, "你")

	if m.PromptIdx != 1 {
		t.Fatalf("PromptIdx = %d, want 1 rune", m.PromptIdx)
	}
	if m.InputLen != 1 {
		t.Fatalf("InputLen = %d, want 1", m.InputLen)
	}
	if m.InputStrPlain != "你" {
		t.Fatalf("InputStrPlain = %q, want %q", m.InputStrPlain, "你")
	}
	if m.State.Hits != 1 || m.State.Errors != 0 {
		t.Fatalf("hits/errors = %d/%d, want 1/0", m.State.Hits, m.State.Errors)
	}
}

func TestBackspaceRemovesWholeUnicodeRune(t *testing.T) {
	m := unicodeTestModel("你好 世界")
	typeChar(&m, "你")

	backspace(&m)

	if m.InputStrPlain != "" {
		t.Fatalf("InputStrPlain = %q, want empty string", m.InputStrPlain)
	}
	if m.InputLen != 0 {
		t.Fatalf("InputLen = %d, want 0", m.InputLen)
	}
	if m.PromptIdx != 0 {
		t.Fatalf("PromptIdx = %d, want 0", m.PromptIdx)
	}
	if len(m.CharLenSlice) != 0 {
		t.Fatalf("CharLenSlice has %d entries, want 0", len(m.CharLenSlice))
	}
}

func TestTypeCharCompletesUnicodeWordByRuneCount(t *testing.T) {
	m := unicodeTestModel("你好 世界")

	typeChar(&m, "你")
	typeChar(&m, "好")

	if m.State.WordCompletions != 1 {
		t.Fatalf("WordCompletions = %d, want 1", m.State.WordCompletions)
	}
	if m.WordIdx != 1 {
		t.Fatalf("WordIdx = %d, want 1", m.WordIdx)
	}
	if m.PromptIdx != 2 {
		t.Fatalf("PromptIdx = %d, want 2 runes", m.PromptIdx)
	}
	if m.InputStrPlain != "" {
		t.Fatalf("InputStrPlain = %q, want reset after word completion", m.InputStrPlain)
	}
}
