package parser

import (
	"testing"
)

func TestParseHSS_UnknownProperty(t *testing.T) {
	input := `box {
		color: red;
		margin: 2;
		padding: 1;
		gap: 2;
		zoom: 10;
	}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("se esperaba 1 regla, se obtuvo %d", len(rules))
	}
	if _, hasZoom := rules[0].Properties["zoom"]; hasZoom {
		t.Error("'zoom' no deberia estar en las propiedades")
	}
	if v, ok := rules[0].Properties["color"]; !ok || v != "red" {
		t.Errorf("esperaba color=red, obtuvo color=%s", v)
	}
	if v, ok := rules[0].Properties["gap"]; !ok || v != "2" {
		t.Errorf("esperaba gap=2, obtuvo gap=%s", v)
	}
	if v, ok := rules[0].Properties["margin"]; !ok || v != "2" {
		t.Errorf("esperaba margin=2, obtuvo margin=%s", v)
	}
	if v, ok := rules[0].Properties["padding"]; !ok || v != "1" {
		t.Errorf("esperaba padding=1, obtuvo padding=%s", v)
	}
}

func TestParseHSS_ColorNames(t *testing.T) {
	input := `text {
		color: red;
	}
	title {
		color: blue;
	}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("se esperaban 2 reglas, se obtuvo %d", len(rules))
	}
	checkRule(t, rules[0], "text", map[string]string{"color": "red"})
	checkRule(t, rules[1], "title", map[string]string{"color": "blue"})
}

func TestParseHSS_AllProperties(t *testing.T) {
	input := `box {
		width: 80;
		height: 24;
		grow: 1;
		gap: 2;
		direction: row;
		align: center;
		justify: space-between;
		color: white;
		bg: blue;
		border: true;
		padding: 2;
		margin: 1;
		min-width: 10;
		min-height: 5;
		max-width: 200;
		max-height: 100;
		overflow: scroll;
	}`
	rules, err := ParseHSS(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("se esperaba 1 regla, se obtuvo %d", len(rules))
	}
	expected := map[string]string{
		"width":      "80",
		"height":     "24",
		"grow":       "1",
		"gap":        "2",
		"direction":  "row",
		"align":      "center",
		"justify":    "space-between",
		"color":      "white",
		"bg":         "blue",
		"border":     "true",
		"padding":    "2",
		"margin":     "1",
		"min-width":  "10",
		"min-height": "5",
		"max-width":  "200",
		"max-height": "100",
		"overflow":   "scroll",
	}
	checkRule(t, rules[0], "box", expected)
}
