package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// TestLayout_PageCustomSize: page con Width/Height explicitos.
func TestLayout_PageCustomSize(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Width: 60, Height: 30,
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 80, 24)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 60, H: 30})
}

// TestLayout_PageSingleChildNoGrow: layoutPage con 1 hijo sin grow.
func TestLayout_PageSingleChildNoGrow(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 10, Height: 5}},
			},
		}},
	}
	Layout(doc, 80, 24)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
}

// TestLayout_BorderPadding: layoutNode con border y padding.
func TestLayout_BorderPadding(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Grow: 1, Border: true, Padding: 2, Direction: "row"},
				Children: []ast.ElementNode{
					{Style: ast.ComputedStyle{Grow: 1}},
					{Style: ast.ComputedStyle{Grow: 1}},
				},
			}},
		}},
	}
	Layout(doc, 80, 24)
	parent := &doc.Pages[0].Children[0]
	assertBox(t, parent.BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
	// border(1) + padding(2) → content area: x=3, y=3, w=74, h=18
	// row grow children: each 74/2=37
	assertBox(t, parent.Children[0].BoundBox, ast.BoundBox{X: 3, Y: 3, W: 37, H: 18})
	assertBox(t, parent.Children[1].BoundBox, ast.BoundBox{X: 40, Y: 3, W: 37, H: 18})
}

// TestLayout_TextareaPostExpansion: textarea sin width explicito se expande.
func TestLayout_TextareaPostExpansion(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Tag:  "textarea",
				Text: "this is a very long text that should expand the box",
			}},
		}},
	}
	Layout(doc, 20, 24)
	// intrinsic width = 51 > 20 → post-expansion: BoundBox.W = 51
	assertEq(t, doc.Pages[0].Children[0].BoundBox.W, 51, "textarea post-expand width")
}

// TestLayout_ExplicitWidthHeight: nodo con W+H explicitos → early return.
func TestLayout_ExplicitWidthHeight(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Width: 50, Height: 10, Direction: "row"},
				Children: []ast.ElementNode{
					{Style: ast.ComputedStyle{Grow: 1}},
				},
			}},
		}},
	}
	Layout(doc, 80, 24)
	// Early return: BoundBox stays as assigned by layoutNode, not modified by post-expansion
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
}
