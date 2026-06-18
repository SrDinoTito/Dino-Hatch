// Helper functions de navegacion y busqueda en el AST.
package main

import "github.com/srdino/dino-hatch/internal/ast"

// ContentHeight retorna la altura total del contenido del documento.
func (s *AppState) ContentHeight() int {
	maxH := 0
	var walk func(n *ast.ElementNode)
	walk = func(n *ast.ElementNode) {
		bottom := n.BoundBox.Y + n.BoundBox.H
		if bottom > maxH {
			maxH = bottom
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
	return maxH
}

// FindElementByID busca un elemento por su ID en el documento principal.
func (s *AppState) FindElementByID(id string) *ast.ElementNode {
	return findElementByID(s.Doc, id)
}

// findElementByID busca recursivamente un elemento por su ID en el documento dado.
func findElementByID(doc *ast.Document, id string) *ast.ElementNode {
	var found *ast.ElementNode
	var walk func(n *ast.ElementNode)
	walk = func(n *ast.ElementNode) {
		if found != nil {
			return
		}
		if n.Attrs["id"] == id {
			found = n
			return
		}
		for i := range n.Children {
			walk(&n.Children[i])
		}
	}
	for pi := range doc.Pages {
		for i := range doc.Pages[pi].Children {
			walk(&doc.Pages[pi].Children[i])
		}
	}
	return found
}
