package parser

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestParseHSS_Basic(t *testing.T) {
	input := `box {
		direction: column;
		gap: 1;
	}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("se esperaba 1 regla, se obtuvo %d", len(rules))
	}
	checkRule(t, rules[0], "box", map[string]string{
		"direction": "column",
		"gap":       "1",
	})
}

func TestParseHSS_MultipleRules(t *testing.T) {
	input := `box {
		direction: column;
	}
	text {
		color: white;
	}
	title {
		color: yellow;
	}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("se esperaban 3 reglas, se obtuvo %d", len(rules))
	}
	checkRule(t, rules[0], "box", map[string]string{"direction": "column"})
	checkRule(t, rules[1], "text", map[string]string{"color": "white"})
	checkRule(t, rules[2], "title", map[string]string{"color": "yellow"})
}

func TestParseHSS_Empty(t *testing.T) {
	rules, err := ParseHSS("")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("se esperaban 0 reglas, se obtuvo %d", len(rules))
	}
}

func TestParseHSS_EmptyWhitespace(t *testing.T) {
	rules, err := ParseHSS("   \n\t   ")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("se esperaban 0 reglas, se obtuvo %d", len(rules))
	}
}

func TestParseHSS_Minified(t *testing.T) {
	input := `box{width:80;height:24;}text{color:red;}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("se esperaban 2 reglas, se obtuvo %d", len(rules))
	}
	checkRule(t, rules[0], "box", map[string]string{"width": "80", "height": "24"})
	checkRule(t, rules[1], "text", map[string]string{"color": "red"})
}

// checkRule verifica que una regla tenga el selector y propiedades esperados.
func checkRule(t *testing.T, rule ast.StyleRule, expectedSelector string, expectedProps map[string]string) {
	t.Helper()
	if rule.Selector != expectedSelector {
		t.Errorf("selector: esperado '%s', obtuvo '%s'", expectedSelector, rule.Selector)
	}
	if len(rule.Properties) != len(expectedProps) {
		t.Fatalf("numero de propiedades: esperado %d, obtuvo %d", len(expectedProps), len(rule.Properties))
	}
	for k, v := range expectedProps {
		got, ok := rule.Properties[k]
		if !ok {
			t.Errorf("propiedad '%s' no encontrada en resultado", k)
			continue
		}
		if got != v {
			t.Errorf("propiedad '%s': esperado '%s', obtuvo '%s'", k, v, got)
		}
	}
}
