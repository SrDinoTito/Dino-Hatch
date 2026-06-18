package parser

import (
	"reflect"
	"testing"
)

func TestParseCSSVars_Empty(t *testing.T) {
	vars := ParseCSSVars("")
	if vars != nil {
		t.Errorf("contenido vacio: esperaba nil, obtuvo %v", vars)
	}
}

func TestParseCSSVars_Whitespace(t *testing.T) {
	vars := ParseCSSVars("   \n\t   ")
	if vars != nil {
		t.Errorf("solo whitespace: esperaba nil, obtuvo %v", vars)
	}
}

func TestParseCSSVars_NoRoot(t *testing.T) {
	input := `box {
		color: red;
	}`
	vars := ParseCSSVars(input)
	if vars != nil {
		t.Errorf("sin :root: esperaba nil, obtuvo %v", vars)
	}
}

func TestParseCSSVars_SingleVar(t *testing.T) {
	input := `:root {
		--bg: #333;
	}`
	vars := ParseCSSVars(input)
	expected := map[string]string{"--bg": "#333"}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf(":root simple: esperaba %v, obtuvo %v", expected, vars)
	}
}

func TestParseCSSVars_MultipleVars(t *testing.T) {
	input := `:root {
		--bg: #333;
		--fg: #fff;
		--primary: #ff0000;
	}`
	vars := ParseCSSVars(input)
	expected := map[string]string{
		"--bg":      "#333",
		"--fg":      "#fff",
		"--primary": "#ff0000",
	}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf(":root multiple: esperaba %v, obtuvo %v", expected, vars)
	}
}

func TestParseCSSVars_RootWithOtherRules(t *testing.T) {
	input := `:root {
		--bg: #222;
	}
	box {
		color: red;
	}
	text {
		color: white;
	}`
	vars := ParseCSSVars(input)
	expected := map[string]string{"--bg": "#222"}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf(":root con otras reglas: esperaba %v, obtuvo %v", expected, vars)
	}
}

func TestParseCSSVars_OnlyNonRoot(t *testing.T) {
	input := `box {
		color: red;
	}
	text {
		color: white;
	}`
	vars := ParseCSSVars(input)
	if vars != nil {
		t.Errorf("solo reglas no :root: esperaba nil, obtuvo %v", vars)
	}
}

func TestParseCSSVars_NoVarsInRoot(t *testing.T) {
	input := `:root {
		color: red;
	}`
	vars := ParseCSSVars(input)
	// ":root" existe pero no hay vars "--" -> nil
	if vars != nil {
		t.Errorf(":root sin vars --: esperaba nil, obtuvo %v", vars)
	}
}

func TestParseCSSVars_Minified(t *testing.T) {
	input := `:root{--bg:#333;--fg:#fff;}box{color:red;}`
	vars := ParseCSSVars(input)
	expected := map[string]string{"--bg": "#333", "--fg": "#fff"}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf("minified: esperaba %v, obtuvo %v", expected, vars)
	}
}

func TestParseCSSVars_MalformedDecl(t *testing.T) {
	input := `:root {
		--bg: #333;
		invalid-line;
		--fg: #fff;
	}`
	vars := ParseCSSVars(input)
	expected := map[string]string{"--bg": "#333", "--fg": "#fff"}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf("decl malformed: esperaba %v, obtuvo %v", expected, vars)
	}
}
