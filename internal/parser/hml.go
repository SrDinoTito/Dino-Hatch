// Package parser implementa el parseo de archivos .hml (TermML + HSS).
//
// Usa encoding/xml con token-based parsing para convertir el XML-like
// del HML en un AST, extrayendo los bloques <style> como texto raw
// para que el parser HSS los procese despues.
package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// ParseHML parsea un archivo .hml completo y devuelve un Document.
// styles contiene el contenido raw de cada bloque <style> para
// que el parser HSS lo procese despues.
func ParseHML(data []byte) (*ast.Document, []string, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil, fmt.Errorf("parser: archivo vacio")
	}

	dec := xml.NewDecoder(strings.NewReader(string(data)))
	doc := &ast.Document{}
	var styles []string

	// El primer token debe ser <page>
	tok, err := dec.Token()
	if err != nil {
		return nil, nil, fmt.Errorf("parser: error leyendo XML: %w", err)
	}

	startEl, ok := tok.(xml.StartElement)
	if !ok {
		return nil, nil, fmt.Errorf("parser: se esperaba <page>, got %T", tok)
	}
	if startEl.Name.Local != "page" {
		return nil, nil, fmt.Errorf("parser: elemento raiz debe ser <page>, got <%s>", startEl.Name.Local)
	}

	page, err := pageFromStartElement(startEl)
	if err != nil {
		return nil, nil, err
	}

	// Procesar hijos de <page> hasta </page>
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return nil, nil, fmt.Errorf("parser: falta cierre </page>")
		}
		if err != nil {
			return nil, nil, fmt.Errorf("parser: error XML: %w", err)
		}

		switch t := tok.(type) {
		case xml.EndElement:
			if t.Name.Local == "page" {
				doc.Pages = append(doc.Pages, *page)
				return doc, styles, nil
			}
		case xml.StartElement:
			switch t.Name.Local {
			case "style":
				content, err := parseStyleContent(dec)
				if err != nil {
					return nil, nil, err
				}
				styles = append(styles, content)
			default:
				child, err := parseElement(dec, t)
				if err != nil {
					return nil, nil, err
				}
				page.Children = append(page.Children, *child)
			}
		// ignorar comentarios y texto fuera de elementos
		case xml.CharData:
		case xml.Comment:
		}
	}
}
