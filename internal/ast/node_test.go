package ast

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestDefaultStyle verifica que DefaultStyle devuelve valores correctos
func TestDefaultStyle(t *testing.T) {
	s := DefaultStyle()

	if s.Width != 0 {
		t.Errorf("DefaultStyle Width = %d, want 0", s.Width)
	}
	if s.Height != 0 {
		t.Errorf("DefaultStyle Height = %d, want 0", s.Height)
	}
	if s.Grow != 0 {
		t.Errorf("DefaultStyle Grow = %f, want 0", s.Grow)
	}
	if s.Gap != 0 {
		t.Errorf("DefaultStyle Gap = %d, want 0", s.Gap)
	}
	if s.Align != "stretch" {
		t.Errorf("DefaultStyle Align = %s, want stretch", s.Align)
	}
	if s.Justify != "start" {
		t.Errorf("DefaultStyle Justify = %s, want start", s.Justify)
	}
	if s.Direction != "column" {
		t.Errorf("DefaultStyle Direction = %s, want column", s.Direction)
	}
	if s.Color != tcell.ColorWhite {
		t.Errorf("DefaultStyle Color = %v, want ColorWhite", s.Color)
	}
	if s.BgColor != tcell.ColorReset {
		t.Errorf("DefaultStyle BgColor = %v, want ColorReset", s.BgColor)
	}
	if s.Border {
		t.Error("DefaultStyle Border = true, want false")
	}
	if s.Padding != 0 {
		t.Errorf("DefaultStyle Padding = %d, want 0", s.Padding)
	}
	if s.Margin != 0 {
		t.Errorf("DefaultStyle Margin = %d, want 0", s.Margin)
	}
	if s.MinWidth != 0 {
		t.Errorf("DefaultStyle MinWidth = %d, want 0", s.MinWidth)
	}
	if s.MinHeight != 0 {
		t.Errorf("DefaultStyle MinHeight = %d, want 0", s.MinHeight)
	}
	if s.MaxWidth != 0 {
		t.Errorf("DefaultStyle MaxWidth = %d, want 0", s.MaxWidth)
	}
	if s.MaxHeight != 0 {
		t.Errorf("DefaultStyle MaxHeight = %d, want 0", s.MaxHeight)
	}
	if s.Overflow != "visible" {
		t.Errorf("DefaultStyle Overflow = %q, want visible", s.Overflow)
	}
}

// TestBoundBoxZero verifica que BoundBox cero funciona
func TestBoundBoxZero(t *testing.T) {
	bb := BoundBox{}
	if bb.X != 0 || bb.Y != 0 || bb.W != 0 || bb.H != 0 {
		t.Errorf("BoundBox zero = %+v, want {0,0,0,0}", bb)
	}
}

// TestStyleRuleEmpty verifica que StyleRule vacío funciona
func TestStyleRuleEmpty(t *testing.T) {
	sr := StyleRule{}
	if sr.Selector != "" {
		t.Errorf("StyleRule.Selector = %q, want empty", sr.Selector)
	}
	if len(sr.Properties) != 0 {
		t.Errorf("StyleRule.Properties = %v, want empty", sr.Properties)
	}
}

// TestDocumentInit verifica que Document se inicializa vacío
func TestDocumentInit(t *testing.T) {
	doc := &Document{}
	if len(doc.Pages) != 0 {
		t.Errorf("Document.Pages = %v, want empty", doc.Pages)
	}
}

// TestElementNodeDefault verifica que ElementNode se inicializa con valores cero
func TestElementNodeDefault(t *testing.T) {
	en := ElementNode{}
	if en.Tag != "" {
		t.Errorf("ElementNode.Tag = %q, want empty", en.Tag)
	}
	if en.Text != "" {
		t.Errorf("ElementNode.Text = %q, want empty", en.Text)
	}
	if len(en.Children) != 0 {
		t.Errorf("ElementNode.Children = %v, want empty", en.Children)
	}
}
