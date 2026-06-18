package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestApplyProps_PaddingMarginPositive(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "box",
			Attrs: map[string]string{
				"padding": "5",
				"margin":  "3",
			},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Padding != 5 {
		t.Errorf("padding: got %d, want 5", s.Padding)
	}
	if s.Margin != 3 {
		t.Errorf("margin: got %d, want 3", s.Margin)
	}
}

func TestApplyProps_PaddingMarginZero(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "box",
			Attrs: map[string]string{
				"padding": "0",
				"margin":  "0",
			},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Padding != 0 {
		t.Errorf("padding 0: got %d, want 0", s.Padding)
	}
	if s.Margin != 0 {
		t.Errorf("margin 0: got %d, want 0", s.Margin)
	}
}

func TestApplyProps_Clamping(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "box",
			Attrs: map[string]string{
				"min-width":  "10",
				"min-height": "5",
				"max-width":  "200",
				"max-height": "100",
			},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.MinWidth != 10 {
		t.Errorf("min-width: got %d, want 10", s.MinWidth)
	}
	if s.MinHeight != 5 {
		t.Errorf("min-height: got %d, want 5", s.MinHeight)
	}
	if s.MaxWidth != 200 {
		t.Errorf("max-width: got %d, want 200", s.MaxWidth)
	}
	if s.MaxHeight != 100 {
		t.Errorf("max-height: got %d, want 100", s.MaxHeight)
	}
}

func TestApplyProps_TextAlignVAlign(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "box",
			Attrs: map[string]string{
				"text-align": "right",
				"valign":     "bottom",
			},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.TextAlign != "right" {
		t.Errorf("text-align: got %q, want right", s.TextAlign)
	}
	if s.VAlign != "bottom" {
		t.Errorf("valign: got %q, want bottom", s.VAlign)
	}
}

func TestApplyProps_DirectionAlignJustifyGap(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "box",
			Attrs: map[string]string{
				"direction": "row",
				"align":     "center",
				"justify":   "space-between",
				"gap":       "4",
			},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Direction != "row" {
		t.Errorf("direction: got %q", s.Direction)
	}
	if s.Align != "center" {
		t.Errorf("align: got %q", s.Align)
	}
	if s.Justify != "space-between" {
		t.Errorf("justify: got %q", s.Justify)
	}
	if s.Gap != 4 {
		t.Errorf("gap: got %d", s.Gap)
	}
}

func TestApplyProps_BgColor(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"bg": "blue"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.BgColor != tcell.GetColor("blue") {
		t.Errorf("bg: got %v, want blue", s.BgColor)
	}
}

func TestApplyProps_Height(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"height": "10"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Height != 10 {
		t.Errorf("height: got %d, want 10", s.Height)
	}
}
