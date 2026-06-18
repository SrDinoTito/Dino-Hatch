// Punto de entrada de renderizado: orquesta el dibujo del AST sobre CellBuffer.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/render"
)

// renderDoc dibuja todas las paginas del documento en el cell buffer.
// scrollY desplaza verticalmente todo el contenido.
func renderDoc(cb *render.CellBuffer, doc *ast.Document, scrollY int) {
	for i := range doc.Pages {
		renderPage(cb, &doc.Pages[i], scrollY)
	}
}

// renderPage dibuja los hijos directos de una pagina.
func renderPage(cb *render.CellBuffer, page *ast.Page, scrollY int) {
	for i := range page.Children {
		renderNode(cb, &page.Children[i], tcell.ColorReset, scrollY)
	}
}

// renderNode selecciona la funcion de render segun el tag del elemento.
// inheritBg es el color de fondo del padre; los textos lo heredan.
func renderNode(cb *render.CellBuffer, node *ast.ElementNode, inheritBg tcell.Color, scrollY int) {
	bb := node.BoundBox

	// Heredar fondo del padre
	bg := node.Style.BgColor
	if bg == tcell.ColorReset {
		bg = inheritBg
	}
	style := tcell.StyleDefault.
		Foreground(node.Style.Color).
		Background(bg)

	// Sobrescribir con color aleatorio si el modo esta activo
	if state.RandomColorsMode {
		if c, ok := state.BoxColors[node]; ok {
			style = style.Background(c)
			bg = c
		}
	}

	// Aplicar scroll vertical (desplazar Y hacia arriba)
	screenY := bb.Y - scrollY

	switch node.Tag {
	case "box":
		renderBoxContent(cb, node, style, bb.X, screenY, bb.W, bb.H)
	case "button":
		renderButton(cb, node, style, bg, bb.X, screenY, bb.W, bb.H, scrollY)
	case "text":
		renderText(cb, node, style, bb.X, screenY, bb.W, bb.H)
	case "input":
		renderInput(cb, node, style, bb.X, screenY, bb.W, bb.H)
	case "textarea":
		renderTextarea(cb, node, style, bb.X, screenY, bb.W, bb.H)
	}

	// Recursion a hijos: text, input, textarea y button se manejan internamente
	if node.Tag != "text" && node.Tag != "input" && node.Tag != "textarea" && node.Tag != "button" {
		needsClip := node.Style.Overflow == "scroll" || node.Style.Overflow == "hidden"
		for i := range node.Children {
			child := &node.Children[i]
			// D3: overflow scroll/hidden — ajustar Y por scroll del contenedor
			if needsClip {
				if !isChildVisible(child, node.ScrollY, node.BoundBox.H, node.BoundBox.Y) {
					continue // skip child fuera de viewport
				}
				childScrollY := scrollY + node.ScrollY
				renderNode(cb, child, bg, childScrollY)
			} else {
				renderNode(cb, child, bg, scrollY)
			}
		}
	}
}
