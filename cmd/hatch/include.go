// Resolucion de includes <include src="..."> en el pipeline.
// Los includes se resuelven despues del parseo HML y antes de ComputeStyles.
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/parser"
)

// resolveIncludes recorre el AST y reemplaza nodos <include> por el contenido
// parseado del archivo src. Se llama antes de ComputeStyles.
func (s *AppState) resolveIncludes(doc *ast.Document) error {
	baseDir := s.BaseDir
	for i := range doc.Pages {
		if err := walkAndResolve(&doc.Pages[i].Children, s, doc, baseDir); err != nil {
			return err
		}
	}
	return nil
}

// walkAndResolve recorre recursivamente la lista de nodos, resolviendo includes
// y reemplazando cada nodo <include> por los children del componente referenciado.
func walkAndResolve(children *[]ast.ElementNode, s *AppState, doc *ast.Document, baseDir string) error {
	var resolved []ast.ElementNode

	for i := range *children {
		node := &(*children)[i]

		if node.IncludeSrc != "" {
			// (a) Leer archivo del include (ruta relativa al .hml actual)
			compPath := filepath.Join(baseDir, node.IncludeSrc)
			data, err := os.ReadFile(compPath)
			if err != nil {
				log.Printf("Warning: include %s no encontrado: %v", node.IncludeSrc, err)
				continue
			}

			// (b) Parsear el contenido del componente
			compDoc, compStyles, parseErr := parser.ParseHML(data)
			if parseErr != nil {
				log.Printf("Warning: error parseando include %s: %v", node.IncludeSrc, parseErr)
				continue
			}
			if len(compDoc.Pages) == 0 || len(compDoc.Pages[0].Children) == 0 {
				continue
			}

			// (c) Extraer style blocks del componente y mergear en s.StyleRules
			for _, block := range compStyles {
				rules, hssErr := parser.ParseHSS(block)
				if hssErr == nil {
					s.StyleRules = append(s.StyleRules, rules...)
				}
			}

			// (d) Mergear atributos: los del include ganan (override)
			compChildren := compDoc.Pages[0].Children
			compRoot := &compChildren[0]
			if node.Attrs != nil {
				if compRoot.Attrs != nil {
					mergeAttrs(compRoot.Attrs, node.Attrs)
				} else {
					compRoot.Attrs = node.Attrs
				}
			}

			// (e) Si el include tiene id, establecerlo en el root del componente
			if id, ok := node.Attrs["id"]; ok && id != "" {
				compRoot.ID = id
			}

			// Resolver includes anidados dentro del componente
			compBaseDir := filepath.Dir(compPath)
			for ci := range compChildren {
				if err := walkAndResolve(&compChildren[ci].Children, s, doc, compBaseDir); err != nil {
					return err
				}
			}

			// Reemplazar el nodo include por los children del componente
			resolved = append(resolved, compChildren...)
		} else {
			// Recurrir en hijos del nodo actual
			if err := walkAndResolve(&node.Children, s, doc, baseDir); err != nil {
				return err
			}
			resolved = append(resolved, *node)
		}
	}

	*children = resolved
	return nil
}

// mergeAttrs copia atributos de source a target. Source gana en caso de conflicto.
func mergeAttrs(target, source map[string]string) {
	for k, v := range source {
		target[k] = v
	}
}
