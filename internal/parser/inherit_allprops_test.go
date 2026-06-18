package parser

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// checkInheritedStyle verifica propiedades heredadas en un test.
func checkInheritedStyle(t *testing.T, got, parent ast.ComputedStyle, props []string) {
	t.Helper()
	for _, p := range props {
		switch p {
		case "color":
			if got.Color != parent.Color {
				t.Errorf("color deberia heredar: got %v, want %v", got.Color, parent.Color)
			}
		case "bg":
			if got.BgColor != parent.BgColor {
				t.Errorf("bg deberia heredar: got %v, want %v", got.BgColor, parent.BgColor)
			}
		case "overflow":
			if got.Overflow != parent.Overflow {
				t.Errorf("overflow deberia heredar: got %q, want %q", got.Overflow, parent.Overflow)
			}
		case "direction":
			if got.Direction != parent.Direction {
				t.Errorf("direction deberia heredar: got %q, want %q", got.Direction, parent.Direction)
			}
		case "align":
			if got.Align != parent.Align {
				t.Errorf("align deberia heredar: got %q, want %q", got.Align, parent.Align)
			}
		case "justify":
			if got.Justify != parent.Justify {
				t.Errorf("justify deberia heredar: got %q, want %q", got.Justify, parent.Justify)
			}
		case "text-align":
			if got.TextAlign != parent.TextAlign {
				t.Errorf("text-align deberia heredar: got %q, want %q", got.TextAlign, parent.TextAlign)
			}
		case "valign":
			if got.VAlign != parent.VAlign {
				t.Errorf("valign deberia heredar: got %q, want %q", got.VAlign, parent.VAlign)
			}
		}
	}
}

func TestInherit_AllInheritableProps(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "parent",
			Attrs: map[string]string{
				"direction":  "row",
				"align":      "center",
				"justify":    "end",
				"color":      "green",
				"bg":         "black",
				"text-align": "left",
				"valign":     "top",
				"overflow":   "scroll",
			},
			Children: []ast.ElementNode{{
				Tag: "child",
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	p := result.Pages[0].Children[0]
	c := p.Children[0]

	checkInheritedStyle(t, c.Style, p.Style, []string{
		"color", "bg", "direction", "align", "justify",
		"text-align", "valign", "overflow",
	})
}
