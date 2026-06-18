package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// TestIntrinsicSize_Text cubre elementos <text> en intrinsicSize.
func TestIntrinsicSize_Text(t *testing.T) {
	n := &ast.ElementNode{Tag: "text", Text: "hello"}
	assertEq(t, intrinsicSize(n, true), 5, "text width")
	assertEq(t, intrinsicSize(n, false), 1, "text height")
}

// TestIntrinsicSize_Textarea cubre textarea multiline, maxheight, border, padding, min-size.
func TestIntrinsicSize_Textarea(t *testing.T) {
	// multiline: width = max line, height = line count
	ta := &ast.ElementNode{Tag: "textarea", Text: "abc\ndef\nghijk"}
	assertEq(t, intrinsicSize(ta, true), 5, "textarea width")
	assertEq(t, intrinsicSize(ta, false), 3, "textarea height")

	// MaxHeight clamping: 6 lines cap a 4
	tamh := &ast.ElementNode{
		Tag: "textarea", Text: "a\nb\nc\nd\ne\nf",
		Style: ast.ComputedStyle{MaxHeight: 4},
	}
	assertEq(t, intrinsicSize(tamh, false), 4, "textarea maxheight cap")

	// Border + Padding suman a width y height
	tabp := &ast.ElementNode{
		Tag: "textarea", Text: "hi",
		Style: ast.ComputedStyle{Border: true, Padding: 1},
	}
	// width: 2 + border(2) + padding(2) = 6
	assertEq(t, intrinsicSize(tabp, true), 6, "textarea border+padding width")
	// height: 1 + border(2) + padding(2) = 5
	assertEq(t, intrinsicSize(tabp, false), 5, "textarea border+padding height")
}

// TestIntrinsicSize_TextareaVacio cubre textarea vacio: sz<1 → 1.
func TestIntrinsicSize_TextareaVacio(t *testing.T) {
	n := &ast.ElementNode{Tag: "textarea", Text: ""}
	assertEq(t, intrinsicSize(n, true), 1, "textarea empty width")
	assertEq(t, intrinsicSize(n, false), 1, "textarea empty height")
}

// TestIntrinsicSize_ContainerSameDirection suma hijos + gaps + border + padding.
func TestIntrinsicSize_ContainerSameDirection(t *testing.T) {
	// columna: medir altura (misma direccion) con gap, border, padding
	con := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "column", Gap: 1, Border: true, Padding: 1},
		Children: []ast.ElementNode{
			{Style: ast.ComputedStyle{Height: 4}},
			{Style: ast.ComputedStyle{Height: 5}},
		},
	}
	// 4+5+1(gap)=10 + border(2) + padding(2) = 14
	assertEq(t, intrinsicSize(con, false), 14, "col same dir height")

	// fila: medir ancho (misma direccion) con padding
	conRow := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "row", Padding: 1},
		Children: []ast.ElementNode{
			{Style: ast.ComputedStyle{Width: 10}},
			{Style: ast.ComputedStyle{Width: 20}},
		},
	}
	// 10+20=30 + padding(2) = 32
	assertEq(t, intrinsicSize(conRow, true), 32, "row same dir width")
}

// TestIntrinsicSize_ContainerSkipsGrow verifica que hijos con grow se omiten.
func TestIntrinsicSize_ContainerSkipsGrow(t *testing.T) {
	n := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "column"},
		Children: []ast.ElementNode{
			{Style: ast.ComputedStyle{Grow: 1}},
			{Style: ast.ComputedStyle{Height: 7}},
		},
	}
	// Solo cuenta el hijo sin grow
	assertEq(t, intrinsicSize(n, false), 7, "skip grow height")
}

// TestIntrinsicSize_ContainerOppositeDirection calcula max de hijos en cross axis.
func TestIntrinsicSize_ContainerOppositeDirection(t *testing.T) {
	// fila, medir altura (cross axis): max de heights + border + padding
	con := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "row", Border: true, Padding: 1},
		Children: []ast.ElementNode{
			{Style: ast.ComputedStyle{Height: 3}},
			{Style: ast.ComputedStyle{Height: 7}},
		},
	}
	// max=7 + border(2) + padding(2) = 11
	assertEq(t, intrinsicSize(con, false), 11, "row opp dir height")

	// columna, medir ancho (cross axis): max de widths
	con2 := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "column"},
		Children: []ast.ElementNode{
			{Style: ast.ComputedStyle{Width: 15}},
			{Style: ast.ComputedStyle{Width: 7}},
		},
	}
	assertEq(t, intrinsicSize(con2, true), 15, "col opp dir width")
}

// TestIntrinsicSize_ContainerMinMax cubre clamping min/max en intrinsicSize.
func TestIntrinsicSize_ContainerMinMax(t *testing.T) {
	// MinWidth (same direction, row)
	minW := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "row", MinWidth: 50},
		Children: []ast.ElementNode{{Style: ast.ComputedStyle{Width: 5}}},
	}
	assertEq(t, intrinsicSize(minW, true), 50, "same dir minwidth")

	// MaxHeight (same direction, column)
	maxH := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "column", MaxHeight: 10},
		Children: []ast.ElementNode{{Style: ast.ComputedStyle{Height: 50}}},
	}
	assertEq(t, intrinsicSize(maxH, false), 10, "same dir maxheight")

	// MinHeight (opposite direction, row=mide height)
	oppMin := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "row", MinHeight: 20},
		Children: []ast.ElementNode{{Style: ast.ComputedStyle{Height: 5}}},
	}
	assertEq(t, intrinsicSize(oppMin, false), 20, "opp dir minheight")

	// MaxWidth (opposite direction, col=mide width)
	oppMax := &ast.ElementNode{
		Style: ast.ComputedStyle{Direction: "column", MaxWidth: 8},
		Children: []ast.ElementNode{{Style: ast.ComputedStyle{Width: 15}}},
	}
	assertEq(t, intrinsicSize(oppMax, true), 8, "opp dir maxwidth")
}
