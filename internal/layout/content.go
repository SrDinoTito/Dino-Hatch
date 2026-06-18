// Package layout implementa el motor de layout flexbox simplificado.
package layout

import "github.com/srdino/dino-hatch/internal/ast"

// ContentHeight calcula la altura total del contenido (intrínseco) de un elemento,
// sumando alturas de hijos + gaps + padding. Usado para scroll containers.
func ContentHeight(n *ast.ElementNode) int {
	if n == nil || len(n.Children) == 0 {
		return 0
	}
	total := 0
	for _, child := range n.Children {
		if child.BoundBox.H > 0 {
			total += child.BoundBox.H
		} else if child.Style.Height > 0 {
			total += child.Style.Height
		} else {
			total += 1 // mínimo 1 línea
		}
		// gap entre hijos
		total += n.Style.Gap
	}
	// padding
	total += n.Style.Padding * 2
	// menos el último gap
	if len(n.Children) > 0 && n.Style.Gap > 0 {
		total -= n.Style.Gap
	}
	return total
}

// ContentWidth calcula el ancho total del contenido (intrínseco) de un elemento,
// sumando anchos de hijos + gaps + padding. Usado para scroll containers.
func ContentWidth(n *ast.ElementNode) int {
	if n == nil || len(n.Children) == 0 {
		return 0
	}
	total := 0
	for _, child := range n.Children {
		if child.BoundBox.W > 0 {
			total += child.BoundBox.W
		} else if child.Style.Width > 0 {
			total += child.Style.Width
		} else {
			total += 1
		}
		total += n.Style.Gap
	}
	total += n.Style.Padding * 2
	if len(n.Children) > 0 && n.Style.Gap > 0 {
		total -= n.Style.Gap
	}
	return total
}
