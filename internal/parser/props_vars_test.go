package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestResolveCSSVars_NilVars(t *testing.T) {
	got := resolveCSSVars("color: var(--primary)", nil)
	if got != "color: var(--primary)" {
		t.Errorf("nil vars: got %q, want original", got)
	}
}

func TestResolveCSSVars_EmptyVars(t *testing.T) {
	got := resolveCSSVars("color: var(--primary)", map[string]string{})
	if got != "color: var(--primary)" {
		t.Errorf("empty vars: got %q, want original", got)
	}
}

func TestResolveCSSVars_ExistingVar(t *testing.T) {
	vars := map[string]string{"--primary": "#ff0000"}
	got := resolveCSSVars("color: var(--primary)", vars)
	if got != "color: #ff0000" {
		t.Errorf("existing var: got %q, want %q", got, "color: #ff0000")
	}
}

func TestResolveCSSVars_InexistentVar(t *testing.T) {
	vars := map[string]string{"--secondary": "#00ff00"}
	got := resolveCSSVars("color: var(--primary)", vars)
	if got != "color: var(--primary)" {
		t.Errorf("inexistent var: got %q, want original", got)
	}
}

func TestResolveCSSVars_NoVarInValue(t *testing.T) {
	vars := map[string]string{"--primary": "#ff0000"}
	got := resolveCSSVars("color: red", vars)
	if got != "color: red" {
		t.Errorf("no var in value: got %q, want original", got)
	}
}

func TestResolveCSSVars_MultipleVars(t *testing.T) {
	vars := map[string]string{"--primary": "#ff0000", "--size": "4"}
	got := resolveCSSVars("padding: var(--size); color: var(--primary)", vars)
	want := "padding: 4; color: #ff0000"
	if got != want {
		t.Errorf("multiple vars: got %q, want %q", got, want)
	}
}

func TestApplyProps_CSSVarsInline(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"color": "var(--primary)"},
		}}}},
	}
	vars := map[string]string{"--primary": "red"}
	result := ComputeStyles(doc, nil, vars)
	s := result.Pages[0].Children[0].Style
	if s.Color != tcell.GetColor("red") {
		t.Errorf("color con var: got %v, want red", s.Color)
	}
}

func TestApplyProps_CSSVarsInlineInexistent(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"color": "var(--undefined)"},
		}}}},
	}
	vars := map[string]string{"--primary": "red"}
	result := ComputeStyles(doc, nil, vars)
	s := result.Pages[0].Children[0].Style
	// var(--undefined) no se resuelve -> parseColor falla -> color queda como default (white)
	if s.Color != tcell.ColorWhite {
		t.Errorf("color con var inexistente: got %v, want white (default)", s.Color)
	}
}
