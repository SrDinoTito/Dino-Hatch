package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestComputeStyles_WithVarsResolved(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"color": "var(--text-color)"},
		}}}},
	}
	vars := map[string]string{"--text-color": "yellow"}
	result := ComputeStyles(doc, nil, vars)
	s := result.Pages[0].Children[0].Style
	if s.Color != tcell.GetColor("yellow") {
		t.Errorf("color via vars: got %v, want yellow", s.Color)
	}
}

func TestComputeStyles_VarsInHSS(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{Tag: "box"}}}},
	}
	rules := []ast.StyleRule{{
		Selector:   "box",
		Properties: map[string]string{"color": "var(--text-color)"},
	}}
	vars := map[string]string{"--text-color": "green"}
	result := ComputeStyles(doc, rules, vars)
	s := result.Pages[0].Children[0].Style
	if s.Color != tcell.GetColor("green") {
		t.Errorf("color via vars en HSS: got %v, want green", s.Color)
	}
}

func TestComputeStyles_MultiplePages(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{
			{
				Name: "page1",
				Children: []ast.ElementNode{{
					Tag:   "box",
					Attrs: map[string]string{"color": "red"},
				}},
			},
			{
				Name: "page2",
				Children: []ast.ElementNode{{
					Tag:   "box",
					Attrs: map[string]string{"color": "blue"},
				}},
			},
		},
	}
	result := ComputeStyles(doc, nil, nil)
	if len(result.Pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(result.Pages))
	}
	if result.Pages[0].Children[0].Style.Color != tcell.GetColor("red") {
		t.Errorf("page1 color: got %v, want red", result.Pages[0].Children[0].Style.Color)
	}
	if result.Pages[1].Children[0].Style.Color != tcell.GetColor("blue") {
		t.Errorf("page2 color: got %v, want blue", result.Pages[1].Children[0].Style.Color)
	}
}

func TestComputeStyles_NestedChildren(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "outer",
			Children: []ast.ElementNode{{
				Tag: "inner",
				Children: []ast.ElementNode{{
					Tag: "leaf",
				}},
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	outer := result.Pages[0].Children[0]
	inner := outer.Children[0]
	leaf := inner.Children[0]

	if outer.Style.Direction != "column" {
		t.Errorf("outer direction: got %q", outer.Style.Direction)
	}
	if inner.Style.Direction != "column" {
		t.Errorf("inner direction: got %q", inner.Style.Direction)
	}
	if leaf.Style.Direction != "column" {
		t.Errorf("leaf direction: got %q", leaf.Style.Direction)
	}
}

func TestComputeStyles_InlineOverridesHSSWithVars(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"color": "var(--inline-color)"},
		}}}},
	}
	rules := []ast.StyleRule{{
		Selector:   "box",
		Properties: map[string]string{"color": "var(--hss-color)"},
	}}
	vars := map[string]string{
		"--inline-color": "purple",
		"--hss-color":    "orange",
	}
	result := ComputeStyles(doc, rules, vars)
	s := result.Pages[0].Children[0].Style
	// inline deberia pisar HSS
	if s.Color != tcell.GetColor("purple") {
		t.Errorf("inline deberia pisar HSS: got %v, want purple", s.Color)
	}
}
