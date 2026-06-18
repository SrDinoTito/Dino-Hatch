// Indicador visual de scroll en el margen derecho de la pantalla.
// C1: Aparece solo cuando el contenido excede la altura de la terminal.
// D3: drawContainerScrollbar para scrollbars en contenedores con overflow:scroll.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/layout"
	"github.com/srdino/dino-hatch/internal/render"
)

// drawScrollbar dibuja un indicador de scroll vertical en la última columna.
// Barra proporcional al contenido visible, con un marcador de posición.
func drawScrollbar(cb *render.CellBuffer, h, scrollY, maxScroll int) {
	if maxScroll <= 0 {
		return
	}
	barH := max(1, h*h/(h+maxScroll))
	barY := scrollY * (h - barH) / maxScroll
	x := cb.Width() - 1
	if x < 0 {
		return
	}
	barBg := tcell.StyleDefault.Background(tcell.ColorGray)
	lineFg := tcell.StyleDefault.Foreground(tcell.ColorGray)
	for y := 0; y < h; y++ {
		if y >= barY && y < barY+barH {
			cb.Set(x, y, ' ', barBg)
		} else {
			cb.Set(x, y, '│', lineFg)
		}
	}
}

// drawContainerScrollbar dibuja una scrollbar en el margen derecho de un contenedor
// con overflow:scroll. La barra es proporcional al contenido oculto.
// pageScrollY es el scroll global de la pagina (s.ScrollY) para ajuste vertical.
func drawContainerScrollbar(cb *render.CellBuffer, el *ast.ElementNode, pageScrollY int) {
	if el.Style.Overflow != "scroll" {
		return
	}
	contentH := layout.ContentHeight(el)
	maxScroll := contentH - el.BoundBox.H
	if maxScroll <= 0 {
		return
	}
	barH := max(1, el.BoundBox.H*el.BoundBox.H/(el.BoundBox.H+maxScroll))
	screenY := el.BoundBox.Y - pageScrollY
	barY := screenY + (el.ScrollY*(el.BoundBox.H-barH))/maxScroll
	if barY < screenY {
		barY = screenY
	}

	barStyle := tcell.StyleDefault.Background(tcell.ColorGray)
	for i := 0; i < barH && (barY+i) < screenY+el.BoundBox.H; i++ {
		cb.Set(el.BoundBox.X+el.BoundBox.W-1, barY+i, ' ', barStyle)
	}
}

// drawAllContainerScrollbars recorre el AST y dibuja scrollbars para todos
// los contenedores con overflow:scroll.
func drawAllContainerScrollbars(cb *render.CellBuffer, doc *ast.Document, pageScrollY int) {
	var walk func(n *ast.ElementNode)
	walk = func(n *ast.ElementNode) {
		drawContainerScrollbar(cb, n, pageScrollY)
		for i := range n.Children {
			walk(&n.Children[i])
		}
	}
	for i := range doc.Pages {
		for j := range doc.Pages[i].Children {
			walk(&doc.Pages[i].Children[j])
		}
	}
}
