package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestLayout_AlignCenter(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Align = "center"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 60, Height: 20, Grow: 1}},
			},
		}},
	}
	Layout(doc, 100, 50)
	// primary: fills 100, cross: height=20 centered in 50 => y=15
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 15, W: 100, H: 20})
}

func TestLayout_JustifyCenter(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Justify = "center"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 20}},
				{Style: ast.ComputedStyle{Width: 20}},
			},
		}},
	}
	Layout(doc, 100, 50)
	// totalUsed = 40, extraSpace = 60, startOffset = 30
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 30, Y: 0, W: 20, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 50, Y: 0, W: 20, H: 50})
}

func TestLayout_NestedBoxes(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Grow: 1},
				Children: []ast.ElementNode{
					{Style: ast.ComputedStyle{Grow: 1}},
				},
			}},
		}},
	}
	Layout(doc, 80, 24)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
	assertBox(t, doc.Pages[0].Children[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 80, H: 24})
}

func TestLayout_Overflow(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Height: 40}},
				{Style: ast.ComputedStyle{Height: 30}},
				{Style: ast.ComputedStyle{Grow: 1}},
			},
		}},
	}
	Layout(doc, 50, 50)
	// fixed = 40+30 = 70 > 50, remaining = -20, grow child gets 0
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 50, H: 40})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 40, W: 50, H: 30})
	assertBox(t, doc.Pages[0].Children[2].BoundBox, ast.BoundBox{X: 0, Y: 70, W: 50, H: 0})
}
