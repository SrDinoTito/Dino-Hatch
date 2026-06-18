package parser

import "testing"

func TestParseHML_InvalidXML(t *testing.T) {
	data := readTestData(t, "invalid.hml")
	_, _, err := ParseHML(data)
	if err == nil {
		t.Fatal("esperaba error por XML invalido, got nil")
	}
}

func TestParseHML_EmptyFile(t *testing.T) {
	data := readTestData(t, "empty.hml")
	_, _, err := ParseHML(data)
	if err == nil {
		t.Fatal("esperaba error por archivo vacio, got nil")
	}
}

func TestParseHML_TextWithChildren(t *testing.T) {
	data := readTestData(t, "text_with_children.hml")
	_, _, err := ParseHML(data)
	if err == nil {
		t.Fatal("esperaba error por <text> con hijos, got nil")
	}
}

func TestParseHML_MissingClosePage(t *testing.T) {
	// <box> sin </box> causa XML mal formado (</page> no cierra <box>)
	input := []byte(`<page name="test"><box></page>`)
	_, _, err := ParseHML(input)
	if err == nil {
		t.Fatal("esperaba error por XML mal formado (falta cierre de <box>), got nil")
	}
}

func TestParseHML_WrongRoot(t *testing.T) {
	// Elemento raiz no es <page>: debe retornar error
	input := []byte(`<box name="test"></box>`)
	_, _, err := ParseHML(input)
	if err == nil {
		t.Fatal("esperaba error por raiz que no es <page>, got nil")
	}
}
