// Package layout implementa el motor de layout flexbox simplificado.
// Calcula BoundBox para cada nodo del AST segun reglas flexbox.
package layout

import "github.com/srdino/dino-hatch/internal/ast"

// Layout calcula BoundBox para cada nodo del AST.
// containerW, containerH: tamano disponible (usualmente terminal size).
func Layout(doc *ast.Document, containerW, containerH int) {
	for i := range doc.Pages {
		layoutPage(&doc.Pages[i], containerW, containerH)
	}
}

// layoutPage posiciona los hijos directos de una pagina.
// Si hay un solo hijo sin grow, recibe el tamano completo de la terminal
// (para evitar overflow). En cualquier otro caso, se usa layoutChildren
// que distribuye el espacio segun reglas flexbox.
func layoutPage(p *ast.Page, cw, ch int) {
	w, h := cw, ch
	if p.Width > 0 {
		w = p.Width
	}
	if p.Height > 0 {
		h = p.Height
	}
	if len(p.Children) == 1 && p.Children[0].Style.Grow <= 0 {
		layoutNode(&p.Children[0], 0, 0, w, h)
	} else {
		layoutChildren(&p.Children, 0, 0, w, h, &p.Style)
	}
}

// layoutNode establece el BoundBox de un nodo y procesa sus hijos.
func layoutNode(n *ast.ElementNode, x, y, w, h int) {
	n.BoundBox = ast.BoundBox{X: x, Y: y, W: w, H: h}

	// Si el elemento tiene border, restar 1 de cada lado para el contenido interno
	cx, cy, cw, ch := x, y, w, h
	if n.Style.Border && w > 2 && h > 2 {
		cx++
		cy++
		cw -= 2
		ch -= 2
	}
	// Padding interno: reduce el area de contenido para los hijos
	pad := n.Style.Padding
	if pad > 0 && cw > pad*2 && ch > pad*2 {
		cx += pad
		cy += pad
		cw -= pad * 2
		ch -= pad * 2
	}
	layoutChildren(&n.Children, cx, cy, cw, ch, &n.Style)

	// Post-expansion: si el nodo no tiene tamano explicito,
	// expandir BoundBox para cubrir a los hijos
	hasExplicitW := n.Style.Width > 0
	hasExplicitH := n.Style.Height > 0

	// Sin hijos AST: solo textarea necesita post-expansion de ancho
	// (la altura NO se expande para evitar crecimiento infinito al escribir)
	if len(n.Children) == 0 {
		if n.Tag == "textarea" {
			iw := intrinsicSize(n, true)
			if !hasExplicitW && iw > n.BoundBox.W {
				n.BoundBox.W = iw
			}
		}
		return
	}

	if hasExplicitW && hasExplicitH {
		return
	}
	maxW, maxH := 0, 0
	isRow := n.Style.Direction == "row"
	for i := range n.Children {
		child := &n.Children[i]
		var iw, ih int
		if isRow {
			// Row: primary=width (BoundBox, no stretch), cross=height (intrinsic, evita stretch)
			iw = child.BoundBox.W
			ih = intrinsicSize(child, false)
			if child.Tag == "textarea" {
				ih = child.BoundBox.H // textarea: BoundBox fijo y correcto
			}
		} else {
			// Column: primary=height (BoundBox, no stretch), cross=width (intrinsic, evita stretch)
			iw = intrinsicSize(child, true)
			ih = child.BoundBox.H
		}
		childEndX := child.BoundBox.X + iw - x
		childEndY := child.BoundBox.Y + ih - y
		if childEndX > maxW {
			maxW = childEndX
		}
		if childEndY > maxH {
			maxH = childEndY
		}
	}
	if n.Style.Border {
		maxW += 1 // solo borde derecho (izquierdo ya offseteado por X+1)
		maxH += 1 // solo borde inferior (superior ya offseteado por Y+1)
	}
	if !hasExplicitW && maxW > n.BoundBox.W {
		n.BoundBox.W = maxW
	}
	// Recalcular altura solo si los hijos necesitan MAS espacio del asignado.
	// No reducir para evitar inestabilidad: el textarea se controla via
	// estabilizacion selectiva (MaxHeight) y no via post-expansion.
	if !hasExplicitH && maxH > n.BoundBox.H {
		n.BoundBox.H = maxH
	}

	// Clamping min/max sobre el BoundBox final
	if n.Style.MinWidth > 0 && n.BoundBox.W < n.Style.MinWidth {
		n.BoundBox.W = n.Style.MinWidth
	}
	if n.Style.MinHeight > 0 && n.BoundBox.H < n.Style.MinHeight {
		n.BoundBox.H = n.Style.MinHeight
	}
	if n.Style.MaxWidth > 0 && n.BoundBox.W > n.Style.MaxWidth {
		n.BoundBox.W = n.Style.MaxWidth
	}
	if n.Style.MaxHeight > 0 && n.BoundBox.H > n.Style.MaxHeight {
		n.BoundBox.H = n.Style.MaxHeight
	}
}
