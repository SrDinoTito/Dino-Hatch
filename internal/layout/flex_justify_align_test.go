package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// TestLayout_JustifyEnd: justify="end" alinea al final del eje principal.
// Container row 100x50, 2 hijos Width=20 sin grow → hijo0.X=60, hijo1.X=80.
func TestLayout_JustifyEnd(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Justify = "end"
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
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 60, Y: 0, W: 20, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 80, Y: 0, W: 20, H: 50})
}

// TestLayout_JustifyStart: justify="start" (default) alinea al inicio sin offset extra.
// Container row 100x50, 2 hijos Width=15, gap=2 → hijo0.X=0, hijo1.X=17.
func TestLayout_JustifyStart(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Gap = 2
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 15}},
				{Style: ast.ComputedStyle{Width: 15}},
			},
		}},
	}
	Layout(doc, 100, 50)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 15, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 17, Y: 0, W: 15, H: 50})
}

// TestLayout_JustifySpaceBetween: justify="space-between" distribuye espacio extra entre hijos.
// Container row 100x50, 3 hijos Width=10 → extraSpace=70, betweenGap=35.
// hijo0.X=0, hijo1.X=45, hijo2.X=90.
func TestLayout_JustifySpaceBetween(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Justify = "space-between"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 10}},
				{Style: ast.ComputedStyle{Width: 10}},
				{Style: ast.ComputedStyle{Width: 10}},
			},
		}},
	}
	Layout(doc, 100, 50)
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 10, H: 50})
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 45, Y: 0, W: 10, H: 50})
	assertBox(t, doc.Pages[0].Children[2].BoundBox, ast.BoundBox{X: 90, Y: 0, W: 10, H: 50})
}

// TestLayout_AlignEnd: align="end" alinea al final del eje transversal (cross axis).
// Container row 100x50, 2 hijos (el primero Width=60 Height=20, segundo dummy para
// forzar layoutChildren), align="end" → primer hijo.Y=30 (50-20).
func TestLayout_AlignEnd(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Align = "end"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 60, Height: 20}},
				{Style: ast.ComputedStyle{Width: 10, Height: 20}},
			},
		}},
	}
	Layout(doc, 100, 50)
	// cross axis (height): hijo0 tiene Height=20, contenedor tiene 50 → Y=30
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 30, W: 60, H: 20})
}

// TestLayout_Margin: margin como espacio externo entre hermanos en columna.
// Container column 80x50, 2 hijos Height=10 con Margin=2.
// Cada hijo contribuye margin*2 en eje primario→ totalMargin=8.
// hijo0.Y=2 (start margin), hijo1.Y=16 (10+2+2+2).
func TestLayout_Margin(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(), // direction=column, justify=start
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Height: 10, Margin: 2}},
				{Style: ast.ComputedStyle{Height: 10, Margin: 2}},
			},
		}},
	}
	Layout(doc, 80, 50)
	// hijo0: margin start=2, BoundBox.Y=2
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 2, W: 80, H: 10})
	// hijo1: primaryPos tras hijo0 = 2 + 10 + 2 = 14, + margin start=2 → Y=16
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 0, Y: 16, W: 80, H: 10})
}

// TestLayout_MarginRow: margin en row con gap.
// Container row 100x50, 2 hijos Width=10 con Margin=3, gap=2.
// hijo0.X=3 (start margin), hijo1.X=21 (3+10+3+2+3).
func TestLayout_MarginRow(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	ps.Gap = 2
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 10, Margin: 3}},
				{Style: ast.ComputedStyle{Width: 10, Margin: 3}},
			},
		}},
	}
	Layout(doc, 100, 50)
	// hijo0: margin start=3, BoundBox.X=3
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 3, Y: 0, W: 10, H: 50})
	// hijo1: primaryPos tras hijo0 = 3+10+3+2=18, + margin start=3 → X=21
	assertBox(t, doc.Pages[0].Children[1].BoundBox, ast.BoundBox{X: 21, Y: 0, W: 10, H: 50})
}
