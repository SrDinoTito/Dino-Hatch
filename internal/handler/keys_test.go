package handler

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestFormatKey_Rune(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	got := FormatKey(ev)
	if got != "a" {
		t.Errorf("esperaba 'a', got %q", got)
	}
}

func TestFormatKey_ShiftRune(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyRune, 'A', tcell.ModShift)
	got := FormatKey(ev)
	if got != "A" {
		t.Errorf("esperaba 'A', got %q", got)
	}
}

func TestFormatKey_CtrlC(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
	got := FormatKey(ev)
	if got != "Ctrl+C" {
		t.Errorf("esperaba 'Ctrl+C', got %q", got)
	}
}

func TestFormatKey_Escape(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	got := FormatKey(ev)
	if got != "Esc" {
		t.Errorf("esperaba 'Esc', got %q", got)
	}
}
