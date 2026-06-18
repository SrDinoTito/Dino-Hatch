package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestApplyProps_InvalidGrow(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"grow": "not-a-number"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Grow != 0 {
		t.Errorf("grow invalido: got %f, want 0", s.Grow)
	}
}

func TestApplyProps_InvalidWidth(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"width": "abc"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Width != 0 {
		t.Errorf("width invalido: got %d, want 0", s.Width)
	}
}

func TestApplyProps_InvalidHeight(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"height": "xxx"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Height != 0 {
		t.Errorf("height invalido: got %d, want 0", s.Height)
	}
}

func TestApplyProps_InvalidColor(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"color": "not-a-color-name"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Color != tcell.ColorWhite {
		t.Errorf("color invalido: got %v, want white (default)", s.Color)
	}
}

func TestApplyProps_InvalidBgColor(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"bg": "not-a-color"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.BgColor != tcell.ColorReset {
		t.Errorf("bg invalido: got %v, want ColorReset", s.BgColor)
	}
}

func TestApplyProps_InvalidPadding(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"padding": "abc"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Padding != 0 {
		t.Errorf("padding invalido: got %d, want 0", s.Padding)
	}
}

func TestApplyProps_InvalidMargin(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"margin": "abc"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Margin != 0 {
		t.Errorf("margin invalido: got %d, want 0", s.Margin)
	}
}

func TestApplyProps_BorderInvalid(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "xyz"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Border {
		t.Errorf("border=xyz: got true, want false (default)")
	}
}

func TestApplyProps_InvalidGap(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"gap": "abc"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Gap != 0 {
		t.Errorf("gap invalido: got %d, want 0", s.Gap)
	}
}

func TestApplyProps_InvalidMinWidth(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"min-width": "abc"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.MinWidth != 0 {
		t.Errorf("min-width invalido: got %d, want 0", s.MinWidth)
	}
}
