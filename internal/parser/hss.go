// Package parser parsea archivos .hml (HML + HSS).
package parser

import (
	"log"
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// ParseHSS parsea el contenido de un bloque <style> y devuelve reglas.
// Recibe el texto raw entre las etiquetas <style> y </style>.
// Propiedades desconocidas se ignoran con warning (log.Printf).
// Selectores :root se ignoran (no generan StyleRule);
// usa ParseCSSVars para extraer sus variables CSS.
func ParseHSS(styleContent string) ([]ast.StyleRule, error) {
	content := strings.TrimSpace(styleContent)
	if content == "" {
		return nil, nil
	}

	var rules []ast.StyleRule
	i := 0

	for i < len(content) {
		// Avanzar espacios/blancos entre reglas
		for i < len(content) && isWhitespace(content[i]) {
			i++
		}
		if i >= len(content) {
			break
		}

		// Leer selector hasta '{'
		selStart := i
		for i < len(content) && content[i] != '{' {
			i++
		}
		if i >= len(content) {
			break
		}

		selector := strings.TrimSpace(content[selStart:i])

		// Saltar '{'
		i++

		// Leer cuerpo hasta '}' balanceado
		bodyStart := i
		depth := 1
		for i < len(content) && depth > 0 {
			switch content[i] {
			case '{':
				depth++
			case '}':
				depth--
			}
			if depth > 0 {
				i++
			}
		}

		body := content[bodyStart:i]
		props := parseProperties(body)

		if selector == ":root" {
			// :root no genera StyleRule, solo CSS variables
			// (extraer con ParseCSSVars si es necesario)
		} else if len(props) > 0 {
			rules = append(rules, ast.StyleRule{
				Selector:   selector,
				Properties: props,
			})
		}

		if i < len(content) {
			i++ // saltar '}'
		}
	}

	return rules, nil
}

// ParseCSSVars extrae CSS variables (claves "--") del selector :root
// en el contenido de un bloque <style>.
func ParseCSSVars(styleContent string) map[string]string {
	content := strings.TrimSpace(styleContent)
	if content == "" {
		return nil
	}

	vars := make(map[string]string)
	i := 0

	for i < len(content) {
		for i < len(content) && isWhitespace(content[i]) {
			i++
		}
		if i >= len(content) {
			break
		}

		selStart := i
		for i < len(content) && content[i] != '{' {
			i++
		}
		if i >= len(content) {
			break
		}

		selector := strings.TrimSpace(content[selStart:i])

		i++ // saltar '{'

		bodyStart := i
		depth := 1
		for i < len(content) && depth > 0 {
			switch content[i] {
			case '{':
				depth++
			case '}':
				depth--
			}
			if depth > 0 {
				i++
			}
		}

		body := content[bodyStart:i]

		if selector == ":root" {
			decls := strings.Split(body, ";")
			for _, decl := range decls {
				decl = strings.TrimSpace(decl)
				if decl == "" {
					continue
				}
				idx := strings.IndexByte(decl, ':')
				if idx < 0 {
					continue
				}
				name := strings.TrimSpace(decl[:idx])
				value := strings.TrimSpace(decl[idx+1:])
				if strings.HasPrefix(name, "--") {
					vars[name] = value
				}
			}
		}

		if i < len(content) {
			i++ // saltar '}'
		}
	}

	if len(vars) == 0 {
		return nil
	}
	return vars
}

// parseProperties descompone el cuerpo de una regla en pares propiedad:valor.
func parseProperties(body string) map[string]string {
	props := make(map[string]string)
	decls := strings.Split(body, ";")

	for _, decl := range decls {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}

		idx := strings.IndexByte(decl, ':')
		if idx < 0 {
			continue
		}

		name := strings.TrimSpace(decl[:idx])
		value := strings.TrimSpace(decl[idx+1:])

		if !isKnownProp(name) {
			log.Printf("HSS warning: propiedad desconocida '%s' ignorada", name)
			continue
		}

		props[name] = value
	}

	return props
}

// isKnownProp verifica si una propiedad HSS es soportada.
func isKnownProp(name string) bool {
	switch name {
	case "width", "height", "grow", "gap",
		"direction", "align", "justify",
		"color", "bg", "border",
		"text-align", "valign",
		"padding", "margin",
		"min-width", "min-height",
		"max-width", "max-height",
		"overflow":
		return true
	}
	return false
}

// isWhitespace retorna true si el byte es espacio, tab, CR o LF.
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}
