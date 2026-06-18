// Manejo de liberacion de click: auto-copy de seleccion o elemento.
package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

// handleMouseRelease procesa la liberacion del click: copia seleccion o elemento.
func handleMouseRelease(s *AppState, buttons tcell.ButtonMask, el *ast.ElementNode) {
	if buttons == tcell.ButtonNone && s.PrevButtons != tcell.ButtonNone {
		if s.SelActive {
			minX, minY, maxX, maxY := normalizedSelectionRect()
			var selected strings.Builder
			var walk func(n *ast.ElementNode)
			walk = func(n *ast.ElementNode) {
				if n.Tag == "text" {
					screenY := n.BoundBox.Y - s.ScrollY
					for i, r := range n.Text {
						tx := n.BoundBox.X + i
						ty := screenY
						if tx >= minX && tx <= maxX && ty >= minY && ty <= maxY {
							selected.WriteRune(r)
						}
					}
				}
				for i := range n.Children {
					walk(&n.Children[i])
				}
			}
			for pi := range s.Doc.Pages {
				for i := range s.Doc.Pages[pi].Children {
					walk(&s.Doc.Pages[pi].Children[i])
				}
			}
			if selected.Len() > 0 {
				copyToClipboard(selected.String())
			}
			s.SelActive = false
		} else if s.Eng.AutoCopy() && el != nil {
			text := getElementText(el)
			if text != "" {
				copyToClipboard(text)
			}
		}
	}
}
