// Indicador visual de scroll interno para elementos textarea.
// Wave 5: Barra vertical proporcional en el margen derecho del area de contenido.
package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/render"
)

// renderTextareaScrollbar dibuja un indicador de scroll vertical dentro del
// textarea. Aparece solo si el contenido excede el area visible. Usa el
// caracter ▐ (medio bloque) con color gris, proporcional al contenido.
func renderTextareaScrollbar(cb *render.CellBuffer, node *ast.ElementNode, st *inputState, screenX, screenY, w, h int) {
	// Obtener texto efectivo (node.Text o fallback a atributo "value")
	text := node.Text
	if text == "" {
		if v, ok := node.Attrs["value"]; ok {
			text = v
		}
	}

	lines := strings.Split(text, "\n")
	contentLines := len(lines)

	// Altura visible del area de contenido (inner)
	tmaxH := h
	if node.Style.Border {
		tmaxH = h - 2
	}
	visibleLines := tmaxH

	if st.ScrollY > 0 || contentLines > visibleLines {
		maxScroll := max(0, contentLines-visibleLines)
		if maxScroll > 0 {
			barH := max(1, visibleLines*visibleLines/(visibleLines+maxScroll))

			contentTop := screenY
			if node.Style.Border {
				contentTop = screenY + 1
			}
			barY := contentTop + (st.ScrollY*(visibleLines-barH))/maxScroll

			scrollStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
			scrollX := screenX + w - 1
			if node.Style.Border {
				scrollX -= 1
			}

			for i := 0; i < barH && (barY+i) < contentTop+visibleLines; i++ {
				_ = cb.Set(scrollX, barY+i, '▐', scrollStyle)
			}
		}
	}
}
