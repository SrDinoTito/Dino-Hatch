// C2: Orden de navegacion por Tab (elementos focusables en orden DFS).
package main

import "github.com/srdino/dino-hatch/internal/ast"

// buildFocusOrder construye la lista de elementos focusables (input, textarea, button)
// recorriendo el AST en orden DFS.
func (s *AppState) buildFocusOrder() {
	s.FocusOrder = nil
	var walk func(n *ast.ElementNode)
	walk = func(n *ast.ElementNode) {
		if n.Tag == "input" || n.Tag == "textarea" || n.Tag == "button" {
			s.FocusOrder = append(s.FocusOrder, n)
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
}
