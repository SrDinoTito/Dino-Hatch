package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestParseInt_Valid(t *testing.T) {
	n, ok := parseInt("42")
	if !ok || n != 42 {
		t.Errorf("parseInt('42'): got (%d, %v), want (42, true)", n, ok)
	}
}

func TestParseInt_TrimSpaces(t *testing.T) {
	n, ok := parseInt("  7  ")
	if !ok || n != 7 {
		t.Errorf("parseInt('  7  '): got (%d, %v), want (7, true)", n, ok)
	}
}

func TestParseInt_Invalid(t *testing.T) {
	_, ok := parseInt("abc")
	if ok {
		t.Error("parseInt('abc'): got ok=true, want false")
	}
}

func TestParseInt_Empty(t *testing.T) {
	_, ok := parseInt("")
	if ok {
		t.Error("parseInt(''): got ok=true, want false")
	}
}

func TestParseFloat_Valid(t *testing.T) {
	f, ok := parseFloat("3.14")
	if !ok || f != 3.14 {
		t.Errorf("parseFloat('3.14'): got (%f, %v), want (3.14, true)", f, ok)
	}
}

func TestParseFloat_Integer(t *testing.T) {
	f, ok := parseFloat("5")
	if !ok || f != 5.0 {
		t.Errorf("parseFloat('5'): got (%f, %v), want (5.0, true)", f, ok)
	}
}

func TestParseFloat_Invalid(t *testing.T) {
	_, ok := parseFloat("not-a-number")
	if ok {
		t.Error("parseFloat('not-a-number'): got ok=true, want false")
	}
}

func TestParseFloat_Empty(t *testing.T) {
	_, ok := parseFloat("")
	if ok {
		t.Error("parseFloat(''): got ok=true, want false")
	}
}

func TestParseColor_Valid(t *testing.T) {
	c, ok := parseColor("red")
	if !ok || c != tcell.GetColor("red") {
		t.Errorf("parseColor('red'): got (%v, %v), want red", c, ok)
	}
}

func TestParseColor_Default(t *testing.T) {
	c, ok := parseColor("default")
	if !ok || c != tcell.ColorDefault {
		t.Errorf("parseColor('default'): got (%v, %v), want ColorDefault", c, ok)
	}
}

func TestParseColor_DefaultMixedCase(t *testing.T) {
	c, ok := parseColor("DEFAULT")
	if !ok || c != tcell.ColorDefault {
		t.Errorf("parseColor('DEFAULT'): got (%v, %v), want ColorDefault", c, ok)
	}
}

func TestParseColor_Invalid(t *testing.T) {
	_, ok := parseColor("not-a-color")
	if ok {
		t.Error("parseColor('not-a-color'): got ok=true, want false")
	}
}

func TestParseColor_Hex(t *testing.T) {
	c, ok := parseColor("#ff0000")
	if !ok || c == tcell.ColorDefault {
		t.Errorf("parseColor('#ff0000'): got (%v, %v), want a valid color", c, ok)
	}
}

func TestParseBool_True(t *testing.T) {
	b, ok := parseBool("true")
	if !ok || !b {
		t.Errorf("parseBool('true'): got (%v, %v), want (true, true)", b, ok)
	}
}

func TestParseBool_Yes(t *testing.T) {
	b, ok := parseBool("yes")
	if !ok || !b {
		t.Errorf("parseBool('yes'): got (%v, %v), want (true, true)", b, ok)
	}
}

func TestParseBool_One(t *testing.T) {
	b, ok := parseBool("1")
	if !ok || !b {
		t.Errorf("parseBool('1'): got (%v, %v), want (true, true)", b, ok)
	}
}

func TestParseBool_False(t *testing.T) {
	b, ok := parseBool("false")
	if !ok || b {
		t.Errorf("parseBool('false'): got (%v, %v), want (false, true)", b, ok)
	}
}

func TestParseBool_No(t *testing.T) {
	b, ok := parseBool("no")
	if !ok || b {
		t.Errorf("parseBool('no'): got (%v, %v), want (false, true)", b, ok)
	}
}

func TestParseBool_Zero(t *testing.T) {
	b, ok := parseBool("0")
	if !ok || b {
		t.Errorf("parseBool('0'): got (%v, %v), want (false, true)", b, ok)
	}
}

func TestParseBool_Invalid(t *testing.T) {
	_, ok := parseBool("maybe")
	if ok {
		t.Error("parseBool('maybe'): got ok=true, want false")
	}
}

func TestParseBool_MixedCase(t *testing.T) {
	b, ok := parseBool("True")
	if !ok || !b {
		t.Errorf("parseBool('True'): got (%v, %v), want (true, true)", b, ok)
	}
}
