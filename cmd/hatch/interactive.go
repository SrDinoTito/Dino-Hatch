// Funciones de interaccion: hit testing, clipboard.
package main

import (
	"os/exec"
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// hitTest encuentra el elemento mas profundo en las coordenadas dadas (page-space).
func hitTest(doc *ast.Document, px, py int) *ast.ElementNode {
	for i := range doc.Pages {
		if n := hitTestChildren(doc.Pages[i].Children, px, py); n != nil {
			return n
		}
	}
	return nil
}

// hitTestChildren busca recursivamente en los hijos de un elemento.
func hitTestChildren(children []ast.ElementNode, px, py int) *ast.ElementNode {
	for i := range children {
		child := &children[i]
		bb := child.BoundBox
		if px >= bb.X && px < bb.X+bb.W && py >= bb.Y && py < bb.Y+bb.H {
			// Botones: devolver el boton directamente, no sus hijos texto
			if child.Tag == "button" {
				return child
			}
			// Para otros elementos, buscar el mas profundo
			if len(child.Children) > 0 {
				if n := hitTestChildren(child.Children, px, py); n != nil {
					return n
				}
			}
			return child
		}
	}
	return nil
}

// copyToClipboard copia texto al portapapeles del sistema via xclip.
func copyToClipboard(text string) {
	if text == "" {
		return
	}
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	_ = cmd.Run()
}

// getElementText devuelve el texto visible de un elemento.
func getElementText(el *ast.ElementNode) string {
	switch el.Tag {
	case "text", "button", "input", "textarea":
		if el.Text != "" {
			return el.Text
		}
		if v, ok := el.Attrs["value"]; ok {
			return v
		}
	}
	return ""
}
