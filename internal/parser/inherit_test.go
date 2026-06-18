package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestComputeStyles_Inheritance(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Children: []ast.ElementNode{{
				Tag: "box",
				Attrs: map[string]string{
					"color": "green",
				},
				Children: []ast.ElementNode{{
					Tag:  "text",
					Text: "hereda color",
				}},
			}},
		}},
	}
	result := ComputeStyles(doc, nil, nil)

	parent := result.Pages[0].Children[0]
	child := parent.Children[0]

	if parent.Style.Color != tcell.GetColor("green") {
		t.Errorf("parent color: got %v, want green", parent.Style.Color)
	}
	if child.Style.Color != tcell.GetColor("green") {
		t.Errorf("child deberia heredar color: got %v, want green", child.Style.Color)
	}
	if child.Style.Width != 0 {
		t.Errorf("width no se hereda: got %d, want 0", child.Style.Width)
	}
}
