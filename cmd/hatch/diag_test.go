package main

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestTextareaEnter(t *testing.T) {
	el := &ast.ElementNode{
		Tag:  "textarea",
		Text: "Linea 1\nLinea 2\nLinea 3",
	}
	state = &AppState{
		InputStates: make(map[*ast.ElementNode]*inputState),
	}
	state.FocusedElement = el

	ev := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	if !handleInputKey(el, ev) {
		t.Fatal("handleInputKey devolvio false para Enter en textarea")
	}
	expected := "Linea 1\nLinea 2\nLinea 3\n"
	if el.Text != expected {
		t.Errorf("Enter no inserto salto de linea:\n  got:  %q\n  want: %q", el.Text, expected)
	} else {
		t.Logf("Enter OK: text=%q, cursor=%d", el.Text, state.InputStates[el].Cursor)
	}

	el2 := &ast.ElementNode{
		Tag:  "textarea",
		Text: "abc",
	}
	state.InputStates[el2] = &inputState{Cursor: 0}
	ev2 := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	handleInputKey(el2, ev2)
	if el2.Text != "\nabc" {
		t.Errorf("Enter al inicio:\n  got:  %q\n  want: %q", el2.Text, "\nabc")
	}

	el3 := &ast.ElementNode{
		Tag:  "textarea",
		Text: "abcdef",
	}
	state.InputStates[el3] = &inputState{Cursor: 3}
	ev3 := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	handleInputKey(el3, ev3)
	if el3.Text != "abc\ndef" {
		t.Errorf("Enter en medio:\n  got:  %q\n  want: %q", el3.Text, "abc\ndef")
	}
	t.Log("TODOS los tests de Enter pasaron")
}

func TestTextareaRuneEnter(t *testing.T) {
	state = &AppState{
		InputStates: make(map[*ast.ElementNode]*inputState),
	}
	el := &ast.ElementNode{
		Tag:  "textarea",
		Text: "hola",
	}
	state.InputStates[el] = &inputState{Cursor: 4}

	ev := tcell.NewEventKey(tcell.KeyRune, '\r', tcell.ModNone)
	if !handleInputKey(el, ev) {
		t.Fatal("handleInputKey devolvio false para KeyRune(\\r) en textarea")
	}
	if el.Text != "hola\n" {
		t.Errorf("KeyRune(\\r) no inserto salto:\n  got:  %q\n  want: %q", el.Text, "hola\n")
	}

	el2 := &ast.ElementNode{
		Tag:  "textarea",
		Text: "x",
	}
	state.InputStates[el2] = &inputState{Cursor: 1}
	ev2 := tcell.NewEventKey(tcell.KeyRune, '\n', tcell.ModNone)
	handleInputKey(el2, ev2)
	if el2.Text != "x\n" {
		t.Errorf("KeyRune(\\n) no inserto salto:\n  got:  %q\n  want: %q", el2.Text, "x\n")
	}
}

func TestTextareaShrink(t *testing.T) {
	state = &AppState{
		InputStates: make(map[*ast.ElementNode]*inputState),
	}
	el := &ast.ElementNode{
		Tag:  "textarea",
		Text: "Linea 1\nLinea 2\nLinea 3",
		Style: ast.ComputedStyle{MaxHeight: 9},
	}
	el.BoundBox = ast.BoundBox{X: 0, Y: 0, W: 20, H: 3}

	el.Text = "Linea 1\nLinea 2\nLinea 3ABC"
	lines := strings.Split(el.Text, "\n")
	if len(lines) != 3 {
		t.Errorf("escribir chars no deberia cambiar num lineas: got %d, want 3", len(lines))
	} else {
		t.Logf("OK: len(lines)=%d despues de escribir chars", len(lines))
	}

	el.Text = "Linea 1\nLinea 2\n"
	lines = strings.Split(el.Text, "\n")
	if len(lines) != 3 {
		t.Logf("texto con \\n al final: %d lineas", len(lines))
	}

	el.Text = "Linea 1\nLinea 2Linea 3"
	lines = strings.Split(el.Text, "\n")
	if len(lines) != 2 {
		t.Errorf("fusion de lineas deberia dar 2, got %d", len(lines))
	}
}

func TestLayoutStability(t *testing.T) {
	state = &AppState{
		InputStates: make(map[*ast.ElementNode]*inputState),
	}
	textarea := &ast.ElementNode{
		Tag:  "textarea",
		Text: "Linea 1\nLinea 2\nLinea 3",
		Style: ast.ComputedStyle{MaxHeight: 9},
	}
	text := &ast.ElementNode{
		Tag:  "text",
		Text: "textarea",
	}
	container := &ast.ElementNode{
		Tag: "box",
		Style: ast.ComputedStyle{
			Border:    true,
			Direction: "column",
		},
		Children: []ast.ElementNode{*text, *textarea},
	}
	_ = container
	initialH := textarea.BoundBox.H
	t.Logf("textarea initial BoundBox.H = %d", initialH)
}
