package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// TestContentHeight verifica ContentHeight: suma alturas de hijos + gaps + padding.
func TestContentHeight(t *testing.T) {
	t.Run("nil retorna 0", func(t *testing.T) {
		got := ContentHeight(nil)
		assertEq(t, got, 0, "ContentHeight(nil)")
	})

	t.Run("sin hijos retorna 0", func(t *testing.T) {
		n := &ast.ElementNode{Children: []ast.ElementNode{}}
		got := ContentHeight(n)
		assertEq(t, got, 0, "sin hijos")
	})

	t.Run("un hijo con BoundBox.H", func(t *testing.T) {
		n := &ast.ElementNode{
			Style:    ast.ComputedStyle{Padding: 1},
			Children: []ast.ElementNode{
				{BoundBox: ast.BoundBox{H: 5}},
			},
		}
		got := ContentHeight(n)
		// 5 (hijo) + 1*2 (padding) = 7
		assertEq(t, got, 7, "un hijo H=5, padding=1")
	})

	t.Run("dos hijos con gap y padding", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Gap: 2, Padding: 1},
			Children: []ast.ElementNode{
				{BoundBox: ast.BoundBox{H: 10}},
				{BoundBox: ast.BoundBox{H: 20}},
			},
		}
		got := ContentHeight(n)
		// child0(10) + gap(2) + child1(20) + gap(2) + padding*2(2) - lastGap(2) = 34
		assertEq(t, got, 34, "dos hijos H=10,20 gap=2 padding=1")
	})

	t.Run("hijo sin BoundBox.H usa Style.Height", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Padding: 1},
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Height: 7}},
			},
		}
		got := ContentHeight(n)
		// 7 (Style.Height) + 1*2 (padding) = 9
		assertEq(t, got, 9, "hijo Style.Height=7, padding=1")
	})

	t.Run("hijo sin BoundBox.H ni Style.Height usa minimo 1", func(t *testing.T) {
		n := &ast.ElementNode{
			Children: []ast.ElementNode{
				{}, // sin BoundBox.H, sin Style.Height
			},
		}
		got := ContentHeight(n)
		// 1 (mínimo) + 0 (padding) = 1
		assertEq(t, got, 1, "hijo vacío, mínimo 1")
	})

	t.Run("sin gap con padding", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Padding: 2},
			Children: []ast.ElementNode{
				{BoundBox: ast.BoundBox{H: 3}},
				{BoundBox: ast.BoundBox{H: 7}},
			},
		}
		got := ContentHeight(n)
		// 3 + 0(gap) + 7 + 0(gap) + 2*2(padding) = 14
		assertEq(t, got, 14, "dos hijos sin gap, padding=2")
	})
}

// TestContentWidth verifica ContentWidth: suma anchos de hijos + gaps + padding.
func TestContentWidth(t *testing.T) {
	t.Run("nil retorna 0", func(t *testing.T) {
		got := ContentWidth(nil)
		assertEq(t, got, 0, "ContentWidth(nil)")
	})

	t.Run("sin hijos retorna 0", func(t *testing.T) {
		n := &ast.ElementNode{Children: []ast.ElementNode{}}
		got := ContentWidth(n)
		assertEq(t, got, 0, "sin hijos")
	})

	t.Run("un hijo con BoundBox.W", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Padding: 2},
			Children: []ast.ElementNode{
				{BoundBox: ast.BoundBox{W: 15}},
			},
		}
		got := ContentWidth(n)
		// 15 (hijo) + 2*2 (padding) = 19
		assertEq(t, got, 19, "un hijo W=15, padding=2")
	})

	t.Run("dos hijos con gap", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Gap: 1},
			Children: []ast.ElementNode{
				{BoundBox: ast.BoundBox{W: 10}},
				{BoundBox: ast.BoundBox{W: 20}},
			},
		}
		got := ContentWidth(n)
		// child0(10) + gap(1) + child1(20) + gap(1) - lastGap(1) = 31
		assertEq(t, got, 31, "dos hijos W=10,20 gap=1")
	})

	t.Run("hijo sin BoundBox.W usa Style.Width", func(t *testing.T) {
		n := &ast.ElementNode{
			Style: ast.ComputedStyle{Padding: 1},
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Width: 8}},
			},
		}
		got := ContentWidth(n)
		// 8 (Style.Width) + 1*2 (padding) = 10
		assertEq(t, got, 10, "hijo Style.Width=8, padding=1")
	})

	t.Run("hijo sin BoundBox.W ni Style.Width usa minimo 1", func(t *testing.T) {
		n := &ast.ElementNode{
			Children: []ast.ElementNode{
				{}, // sin BoundBox.W, sin Style.Width
			},
		}
		got := ContentWidth(n)
		// 1 (mínimo) + 0 (padding) = 1
		assertEq(t, got, 1, "hijo vacío, mínimo 1")
	})
}
