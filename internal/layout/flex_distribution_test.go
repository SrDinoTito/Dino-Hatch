package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestLayout_WithGap(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Gap = 2
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 80, 24)
	// remaining = 24 - 0 - 2 = 22, each child gets 11
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 11})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 13, W: 80, H: 11})
}

func TestLayout_GrowMixed(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 2}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 80, 24)
	// remaining = 24, totalGrow = 3
	// c1: 24*2/3 = 16, c2: 24*1/3 = 8
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 16})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 16, W: 80, H: 8})
}

func TestLayout_NoGrow(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 30}},
				{Style: ast.ComputedStyle{Width: 50}},
			},
		}},
	}
	Layout(doc, 100, 50)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 30, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 30, Y: 0, W: 50, H: 50})
}
