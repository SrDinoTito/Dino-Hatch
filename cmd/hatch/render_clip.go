// D3: Funciones auxiliares de clipping para overflow scroll/hidden en contenedores.
// Chequean si un hijo es visible dentro del viewport del padre y limitan su altura.
package main

import "github.com/srdino/dino-hatch/internal/ast"

// isChildVisible chequea si un hijo es visualmente visible dentro del viewport
// del padre, dados el scroll offset del padre y el BoundBox del padre.
// parentScrollY: node.ScrollY (desplazamiento interno del contenedor).
// parentHeight: node.BoundBox.H (altura del area visible del contenedor).
// parentY: node.BoundBox.Y (posicion Y absoluta del contenedor).
func isChildVisible(child *ast.ElementNode, parentScrollY, parentHeight, parentY int) bool {
	// La posicion visual del hijo = child.BoundBox.Y - parentScrollY
	// Es visible si se superpone con el area del padre
	childVisualY := child.BoundBox.Y - parentScrollY
	childBottom := childVisualY + child.BoundBox.H
	parentBottom := parentY + parentHeight
	return childBottom > parentY && childVisualY < parentBottom
}

// clipHeight retorna la altura visible de un hijo dentro del viewport del padre,
// limitando al area visible. Clip en top y bottom segun superposicion.
func clipHeight(child *ast.ElementNode, parentScrollY, parentHeight, parentY int) int {
	childVisualY := child.BoundBox.Y - parentScrollY
	childBottom := childVisualY + child.BoundBox.H
	parentBottom := parentY + parentHeight

	topClip := 0
	if childVisualY < parentY {
		topClip = parentY - childVisualY
	}
	bottomClip := 0
	if childBottom > parentBottom {
		bottomClip = childBottom - parentBottom
	}
	return child.BoundBox.H - topClip - bottomClip
}
