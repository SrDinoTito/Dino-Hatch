// Package layout implementa el motor de layout flexbox simplificado.
package layout

import (
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// intrinsicSize calcula el tamano intrinseco de un elemento
// en el eje primario (measureWidth=true → ancho, false → alto).
// Salta hijos con grow porque su tamano depende del espacio disponible.
func intrinsicSize(n *ast.ElementNode, measureWidth bool) int {
	if n.Tag == "text" {
		if measureWidth {
			return len(n.Text)
		}
		return 1
	}

	if n.Tag == "textarea" {
		lines := strings.Split(n.Text, "\n")
		var sz int
		if measureWidth {
			maxW := 0
			for _, line := range lines {
				if len(line) > maxW {
					maxW = len(line)
				}
			}
			sz = maxW
		} else {
			sz = len(lines)
			// Respetar MaxHeight como tope absoluto (migrado a ComputedStyle)
			if n.Style.MaxHeight > 0 && sz > n.Style.MaxHeight {
				sz = n.Style.MaxHeight
			}
		}
		if sz < 1 {
			sz = 1
		}
		if n.Style.Border {
			sz += 2
		}
		if n.Style.Padding > 0 {
			sz += n.Style.Padding * 2
		}
		return sz
	}

	isRow := n.Style.Direction == "row"
	gap := n.Style.Gap

	// Misma direccion: sumar hijos + gaps en el eje primario
	if isRow == measureWidth {
		total := 0
		childCount := 0
		for i := range n.Children {
			child := &n.Children[i]
			if child.Style.Grow > 0 {
				continue
			}
			childCount++
			var sz int
			if measureWidth {
				sz = child.Style.Width
			} else {
				sz = child.Style.Height
			}
			if sz <= 0 {
				sz = intrinsicSize(child, measureWidth)
			}
			total += sz
		}
		if childCount > 1 {
			total += gap * (childCount - 1)
		}
		if total < 1 {
			total = 1
		}
		if n.Style.Border {
			total += 2
		}
		if n.Style.Padding > 0 {
			total += n.Style.Padding * 2
		}
		// Clamping min/max
		if measureWidth {
			if n.Style.MinWidth > 0 && total < n.Style.MinWidth {
				total = n.Style.MinWidth
			}
			if n.Style.MaxWidth > 0 && total > n.Style.MaxWidth {
				total = n.Style.MaxWidth
			}
		} else {
			if n.Style.MinHeight > 0 && total < n.Style.MinHeight {
				total = n.Style.MinHeight
			}
			if n.Style.MaxHeight > 0 && total > n.Style.MaxHeight {
				total = n.Style.MaxHeight
			}
		}
		return total
	}

	// Direccion opuesta: max de hijos en este eje transversal
	maxSz := 0
	for i := range n.Children {
		child := &n.Children[i]
		var sz int
		if measureWidth {
			sz = child.Style.Width
		} else {
			sz = child.Style.Height
		}
		if sz <= 0 {
			sz = intrinsicSize(child, measureWidth)
		}
		if sz > maxSz {
			maxSz = sz
		}
	}
	if maxSz < 1 {
		maxSz = 1
	}
	if n.Style.Border {
		maxSz += 2
	}
	if n.Style.Padding > 0 {
		maxSz += n.Style.Padding * 2
	}
	// Clamping min/max
	if measureWidth {
		if n.Style.MinWidth > 0 && maxSz < n.Style.MinWidth {
			maxSz = n.Style.MinWidth
		}
		if n.Style.MaxWidth > 0 && maxSz > n.Style.MaxWidth {
			maxSz = n.Style.MaxWidth
		}
	} else {
		if n.Style.MinHeight > 0 && maxSz < n.Style.MinHeight {
			maxSz = n.Style.MinHeight
		}
		if n.Style.MaxHeight > 0 && maxSz > n.Style.MaxHeight {
			maxSz = n.Style.MaxHeight
		}
	}
	return maxSz
}
