package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestLayout_OneChild(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 80, 24)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
}

func TestLayout_TwoChildren_Row(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
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
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 40, H: 24})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 40, Y: 0, W: 40, H: 24})
}

func TestLayout_TwoChildren_Column(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 80, 24)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 12})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 12, W: 80, H: 12})
}

func TestLayout_DirectionRow(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 100, 50)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 50, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 50, Y: 0, W: 50, H: 50})
}

func TestLayout_DirectionColumn(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "column"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Grow: 1}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 100, 50)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 100, H: 25})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 25, W: 100, H: 25})
}
