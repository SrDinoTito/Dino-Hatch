package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

// TestLayout_MinWidthClamp: MinWidth via layoutNode (cross axis en columna).
// El hijo tiene MinWidth=50 pero layout le da W=40 (cross axis de columna 40x24).
// Como tiene hijos propios, layoutNode NO retorna temprano y luego clampa W a MinWidth.
func TestLayout_MinWidthClamp(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Grow: 1, MinWidth: 50},
				// Sub-hijo para que layoutNode no haga early return
				// y llegue al clamping en lineas 117-118
				Children: []ast.ElementNode{{
					Style: ast.ComputedStyle{Width: 30, Height: 5},
				}},
			}},
		}},
	}
	// Layout(doc, 40, 24): columna, cross axis = 40 < MinWidth=50
	Layout(doc, 40, 24)
	// post-expansion no cambia W (hijo dentro de 40), MinWidth clampa 40→50
	assertEq(t, doc.Pages[0].Children[0].BoundBox.W, 50, "minwidth clamp")
}

// TestLayout_MinHeightClamp: MinHeight via layoutChildren (overflow en eje primario).
// 2 hijos fijos Height=8 suman 16 > container 15 → overflow.
// El grow child recibe 0 por overflow, MinHeight=10 lo clampa a 10.
func TestLayout_MinHeightClamp(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{
				{Style: ast.ComputedStyle{Height: 8}},
				{Style: ast.ComputedStyle{Height: 8}},
				{Style: ast.ComputedStyle{Grow: 1, MinHeight: 10}},
			},
		}},
	}
	Layout(doc, 30, 15)
	// overflow: remaining = 15-16 = -1, grow child recibe 0, MinHeight lo clampa a 10
	assertEq(t, doc.Pages[0].Children[2].BoundBox.H, 10, "minheight overflow clamp")
}

// TestLayout_MaxWidthClamp: MaxWidth via layoutChildren (eje primario en row).
// 1 hijo Grow=1 en row de 100 recibe W=100, MaxWidth=30 lo clampa a 30.
func TestLayout_MaxWidthClamp(t *testing.T) {
	ps := ast.DefaultStyle()
	ps.Direction = "row"
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ps,
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Grow: 1, MaxWidth: 30},
			}},
		}},
	}
	Layout(doc, 100, 50)
	// layoutChildren row: primary=100, clamp MaxWidth=30
	assertEq(t, doc.Pages[0].Children[0].BoundBox.W, 30, "maxwidth clamp")
}

// TestLayout_MaxHeightClamp: MaxHeight via layoutNode (post-expansion + clamp).
// El hijo tiene un sub-hijo alto (Height=200) que expande BoundBox.H via
// post-expansion (linea 113), luego MaxHeight=50 clampa la altura final
// en las lineas 126-128 de flex.go.
func TestLayout_MaxHeightClamp(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{MaxHeight: 50},
				Children: []ast.ElementNode{{
					// Sub-hijo alto: post-expansion expande H sobre el container
					Style: ast.ComputedStyle{Height: 200},
				}},
			}},
		}},
	}
	Layout(doc, 80, 30)
	// layoutNode: H=30 → post-expansion H=200 → MaxHeight clampa a 50
	assertEq(t, doc.Pages[0].Children[0].BoundBox.H, 50, "maxheight post-expansion clamp")
}

// TestLayout_BorderTooSmall: Border con w<=2 no se aplica porque
// la condicion w>2 && h>2 es falsa.
func TestLayout_BorderTooSmall(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Border: true},
			}},
		}},
	}
	Layout(doc, 2, 2)
	// layoutNode: w=2 no supera 2 → border no se resta
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 2, H: 2})
}

// TestLayout_PaddingTooSmall: Padding con espacio insuficiente.
// padding=3 en cw=6 → cw(6) no es > pad*2(6), no se aplica.
func TestLayout_PaddingTooSmall(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{
			Style: ast.DefaultStyle(),
			Children: []ast.ElementNode{{
				Style: ast.ComputedStyle{Padding: 3},
			}},
		}},
	}
	Layout(doc, 6, 6)
	// layoutNode: pad=3, cw=6, 6 no > 6 → padding no se aplica
	assertBox(t, doc.Pages[0].Children[0].BoundBox, ast.BoundBox{X: 0, Y: 0, W: 6, H: 6})
}
