// D3: Hit testing para scroll containers. Encuentra el contenedor scrollable
// mas especifico (mas profundo) bajo las coordenadas del mouse.
package main

import (
	"github.com/srdino/dino-hatch/internal/ast"
)

// findScrollContainer busca recursivamente el contenedor scrollable
// mas especifico bajo las coordenadas (mx, my) del mouse en espacio de pantalla.
// pageScrollY es el scroll global de la pagina (s.ScrollY).
// Retorna nil si no hay ningun contenedor scrollable en esa posicion.
func findScrollContainer(page *ast.Page, mx, my, pageScrollY int) *ast.ElementNode {
	var best *ast.ElementNode
	var walk func(n *ast.ElementNode)
	walk = func(n *ast.ElementNode) {
		// Solo interesan contenedores con overflow:scroll
		if n.Style.Overflow == "scroll" {
			screenY := n.BoundBox.Y - pageScrollY
			if mx >= n.BoundBox.X && mx < n.BoundBox.X+n.BoundBox.W &&
				my >= screenY && my < screenY+n.BoundBox.H {
				best = n
			}
		}
		// Recurrir a hijos (DFS) para encontrar el mas profundo
		for i := range n.Children {
			walk(&n.Children[i])
		}
	}
	for i := range page.Children {
		walk(&page.Children[i])
	}
	return best
}
