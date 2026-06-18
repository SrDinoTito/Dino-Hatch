package parser

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// pageFromStartElement extrae atributos de <page> y devuelve un Page.
func pageFromStartElement(start xml.StartElement) (*ast.Page, error) {
	page := &ast.Page{
		Children: []ast.ElementNode{},
	}
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			page.Name = attr.Value
		case "width":
			w, err := strconv.Atoi(attr.Value)
			if err != nil {
				return nil, fmt.Errorf("parser: width invalido en <page>: %q", attr.Value)
			}
			page.Width = w
		case "height":
			h, err := strconv.Atoi(attr.Value)
			if err != nil {
				return nil, fmt.Errorf("parser: height invalido en <page>: %q", attr.Value)
			}
			page.Height = h
		}
	}
	return page, nil
}

// parseStyleContent lee tokens hasta </style> y devuelve el contenido raw.
func parseStyleContent(dec *xml.Decoder) (string, error) {
	var b strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("parser: error en <style>: %w", err)
		}
		switch t := tok.(type) {
		case xml.EndElement:
			return strings.TrimSpace(b.String()), nil
		case xml.CharData:
			b.Write(t)
		}
	}
}

// parseElement parsea recursivamente un elemento (box o text).
// <text> no puede contener elementos hijos, <box> si.
func parseElement(dec *xml.Decoder, start xml.StartElement) (*ast.ElementNode, error) {
	node := &ast.ElementNode{
		Tag:   start.Name.Local,
		Attrs: attrsFromSlice(start.Attr),
	}

	// Tag include: solo guarda src, no procesa hijos
	if node.Tag == "include" {
		if src, ok := node.Attrs["src"]; ok {
			node.IncludeSrc = src
			delete(node.Attrs, "src")
		}
		// Consumir tokens hasta el cierre
		for {
			tok, err := dec.Token()
			if err != nil {
				return nil, fmt.Errorf("parser: error en <include>: %w", err)
			}
			if _, ok := tok.(xml.EndElement); ok {
				return node, nil
			}
		}
	}

	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("parser: error en <%s>: %w", start.Name.Local, err)
		}

		switch t := tok.(type) {
		case xml.EndElement:
			// Procesar atributos de eventos antes de retornar
			for key, val := range node.Attrs {
				switch key {
				case "onclick", "onchange", "onfocus", "onblur":
					if node.Events == nil {
						node.Events = make(map[string]string)
					}
					node.Events[strings.TrimPrefix(key, "on")] = val
					delete(node.Attrs, key)
				}
			}
			// Setear ID desde Attrs si existe
			if id, ok := node.Attrs["id"]; ok {
				node.ID = id
			}
			return node, nil

		case xml.StartElement:
			if node.Tag == "text" {
				return nil, fmt.Errorf("parser: <text> no puede contener elementos hijos")
			}
			child, err := parseElement(dec, t)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, *child)

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
			if node.Tag == "text" || node.Tag == "button" || node.Tag == "input" || node.Tag == "textarea" {
				node.Text = text
			}
				// <box> y otros ignoran texto directo
			}

		case xml.Comment:
			// ignorar comentarios XML
		}
	}
}

// attrsFromSlice convierte []xml.Attr a map[string]string.
// Retorna nil si no hay atributos.
func attrsFromSlice(attrs []xml.Attr) map[string]string {
	if len(attrs) == 0 {
		return nil
	}
	m := make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}
