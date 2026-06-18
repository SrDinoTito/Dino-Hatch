package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readTestData lee un archivo de fixture del directorio testdata/
func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no se pudo leer %s: %v", path, err)
	}
	return data
}

func TestParseHML_Basic(t *testing.T) {
	data := readTestData(t, "basic.hml")
	doc, styles, err := ParseHML(data)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	if len(doc.Pages) != 1 {
		t.Fatalf("esperaba 1 page, got %d", len(doc.Pages))
	}

	page := doc.Pages[0]
	if page.Name != "main" {
		t.Errorf("page.Name = %q, esperaba %q", page.Name, "main")
	}
	if page.Width != 80 {
		t.Errorf("page.Width = %d, esperaba %d", page.Width, 80)
	}
	if page.Height != 24 {
		t.Errorf("page.Height = %d, esperaba %d", page.Height, 24)
	}
	if len(styles) != 0 {
		t.Errorf("esperaba 0 styles, got %d", len(styles))
	}

	// Verificar estructura: page > box > (text, text)
	if len(page.Children) != 1 {
		t.Fatalf("esperaba 1 hijo en page, got %d", len(page.Children))
	}
	box := page.Children[0]
	if box.Tag != "box" {
		t.Errorf("Tag = %q, esperaba %q", box.Tag, "box")
	}
	if len(box.Children) != 2 {
		t.Fatalf("esperaba 2 hijos en box, got %d", len(box.Children))
	}

	text1 := box.Children[0]
	if text1.Tag != "text" {
		t.Errorf("Tag = %q, esperaba %q", text1.Tag, "text")
	}
	if text1.Text != "Hola" {
		t.Errorf("Text = %q, esperaba %q", text1.Text, "Hola")
	}

	text2 := box.Children[1]
	if text2.Tag != "text" {
		t.Errorf("Tag = %q, esperaba %q", text2.Tag, "text")
	}
	if text2.Text != "Mundo" {
		t.Errorf("Text = %q, esperaba %q", text2.Text, "Mundo")
	}
}

func TestParseHML_WithStyle(t *testing.T) {
	data := readTestData(t, "with_style.hml")
	doc, styles, err := ParseHML(data)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	if len(styles) != 1 {
		t.Fatalf("esperaba 1 style block, got %d", len(styles))
	}
	if !strings.Contains(styles[0], "direction: column") {
		t.Errorf("style deberia contener 'direction: column', got %q", styles[0])
	}
	if !strings.Contains(styles[0], "gap: 1") {
		t.Errorf("style deberia contener 'gap: 1', got %q", styles[0])
	}

	// Verificar que <style> no aparece como elemento hijo
	page := doc.Pages[0]
	for _, child := range page.Children {
		if child.Tag == "style" {
			t.Errorf("<style> no deberia estar en Children de page")
		}
	}
}

func TestParseHML_Attributes(t *testing.T) {
	data := readTestData(t, "basic.hml")
	doc, _, err := ParseHML(data)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	page := doc.Pages[0]

	// Atributos de <page>
	if page.Name != "main" {
		t.Errorf("page.Name = %q", page.Name)
	}
	if page.Width != 80 {
		t.Errorf("page.Width = %d", page.Width)
	}
	if page.Height != 24 {
		t.Errorf("page.Height = %d", page.Height)
	}

	// Atributos de <box>
	box := page.Children[0]
	if box.Attrs["direction"] != "row" {
		t.Errorf("box direction = %q", box.Attrs["direction"])
	}
	if box.Attrs["gap"] != "2" {
		t.Errorf("box gap = %q", box.Attrs["gap"])
	}

	// Atributos de <text>
	text1 := box.Children[0]
	if text1.Attrs["color"] != "green" {
		t.Errorf("text color = %q", text1.Attrs["color"])
	}
}


