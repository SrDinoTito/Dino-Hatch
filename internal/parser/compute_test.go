package parser

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func checkStyle(t *testing.T, got, want ast.ComputedStyle) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ComputedStyle mismatch:\ngot:  %+v\nwant: %+v", got, want)
	}
}

func TestComputeStyles_Defaults(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{
				{Tag: "box"},
			},
		}},
	}
	result := ComputeStyles(doc, nil, nil)
	checkStyle(t, result.Pages[0].Children[0].Style, ast.DefaultStyle())
}

func TestComputeStyles_InlineAttrs(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{
				Tag: "box",
				Attrs: map[string]string{
					"direction": "row",
					"gap":       "3",
					"color":     "red",
				},
			}},
		}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	want := ast.DefaultStyle()
	want.Direction = "row"
	want.Gap = 3
	want.Color = tcell.GetColor("red")
	checkStyle(t, s, want)
}

func TestComputeStyles_HSSRules(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{Tag: "box"}},
		}},
	}
	rules := []ast.StyleRule{{
		Selector: "box",
		Properties: map[string]string{
			"direction": "row",
			"gap":       "2",
			"color":     "blue",
		},
	}}
	result := ComputeStyles(doc, rules, nil)
	s := result.Pages[0].Children[0].Style
	want := ast.DefaultStyle()
	want.Direction = "row"
	want.Gap = 2
	want.Color = tcell.GetColor("blue")
	checkStyle(t, s, want)
}

func TestComputeStyles_InlineOverridesHSS(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{
				Tag: "text",
				Attrs: map[string]string{
					"color": "green",
				},
			}},
		}},
	}
	rules := []ast.StyleRule{{
		Selector: "text",
		Properties: map[string]string{
			"color": "red",
		},
	}}
	result := ComputeStyles(doc, rules, nil)
	s := result.Pages[0].Children[0].Style
	if s.Color != tcell.GetColor("green") {
		t.Errorf("inline deberia pisar HSS: got %v, want green", s.Color)
	}
}

func TestComputeStyles_FullStack(t *testing.T) {
	data := readTestData(t, "with_style.hml")
	doc, rawStyles, err := ParseHML(data)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	var rules []ast.StyleRule
	for _, s := range rawStyles {
		r, err := ParseHSS(s)
		if err != nil {
			t.Fatalf("ParseHSS fallo: %v", err)
		}
		rules = append(rules, r...)
	}

	result := ComputeStyles(doc, rules, nil)

	page := result.Pages[0]
	box := page.Children[0]
	text := box.Children[0]

	if box.Style.Direction != "column" {
		t.Errorf("box direction: got %q, want column", box.Style.Direction)
	}
	if box.Style.Gap != 1 {
		t.Errorf("box gap: got %d, want 1", box.Style.Gap)
	}
	if text.Style.Color != tcell.GetColor("green") {
		t.Errorf("text color: got %v, want green", text.Style.Color)
	}
	if box.Style.BgColor != tcell.ColorReset {
		t.Errorf("box bg: got %v, want ColorReset", box.Style.BgColor)
	}
}

func TestComputeStyles_OverflowDefault(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{
				{Tag: "box"},
			},
		}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Overflow != "visible" {
		t.Errorf("Overflow por defecto = %q, want visible", s.Overflow)
	}
}

func TestComputeStyles_OverflowInline(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{
				Tag: "box",
				Attrs: map[string]string{
					"overflow": "scroll",
				},
			}},
		}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Overflow != "scroll" {
		t.Errorf("Overflow inline = %q, want scroll", s.Overflow)
	}
}

func TestComputeStyles_OverflowHSS(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{Tag: "box"}},
		}},
	}
	rules := []ast.StyleRule{{
		Selector: "box",
		Properties: map[string]string{
			"overflow": "hidden",
		},
	}}
	result := ComputeStyles(doc, rules, nil)
	s := result.Pages[0].Children[0].Style
	if s.Overflow != "hidden" {
		t.Errorf("Overflow HSS = %q, want hidden", s.Overflow)
	}
}
