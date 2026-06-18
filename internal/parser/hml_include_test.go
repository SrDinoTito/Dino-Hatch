package parser

import (
	"testing"
)

// TestParseHML_Include verifica que el parser extraiga IncludeSrc de <include>
// y elimine "src" de Attrs.
func TestParseHML_Include(t *testing.T) {
	input := []byte(`<page name="test">
		<include src="header.hml" />
	</page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	if len(doc.Pages) != 1 {
		t.Fatalf("esperaba 1 page, got %d", len(doc.Pages))
	}
	if len(doc.Pages[0].Children) != 1 {
		t.Fatalf("esperaba 1 hijo en page, got %d", len(doc.Pages[0].Children))
	}

	node := doc.Pages[0].Children[0]
	if node.Tag != "include" {
		t.Errorf("Tag = %q, esperaba %q", node.Tag, "include")
	}
	if node.IncludeSrc != "header.hml" {
		t.Errorf("IncludeSrc = %q, esperaba %q", node.IncludeSrc, "header.hml")
	}
	if _, ok := node.Attrs["src"]; ok {
		t.Errorf("Attrs['src'] deberia estar eliminado, got %q", node.Attrs["src"])
	}
}

// TestParseHML_IncludeNoSrc verifica que <include> sin atributo src
// deje IncludeSrc vacio.
func TestParseHML_IncludeNoSrc(t *testing.T) {
	input := []byte(`<page name="test"><include /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Tag != "include" {
		t.Fatalf("Tag = %q, esperaba %q", node.Tag, "include")
	}
	if node.IncludeSrc != "" {
		t.Errorf("IncludeSrc = %q, esperaba vacio", node.IncludeSrc)
	}
}

// TestParseHML_IncludeWithAttrs verifica que <include> con src y atributos
// adicionales mantenga los extras y solo elimine "src".
func TestParseHML_IncludeWithAttrs(t *testing.T) {
	input := []byte(`<page name="test"><include src="modal.hml" id="myModal" class="hidden" /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.IncludeSrc != "modal.hml" {
		t.Errorf("IncludeSrc = %q, esperaba %q", node.IncludeSrc, "modal.hml")
	}
	if node.Attrs["id"] != "myModal" {
		t.Errorf("Attrs['id'] = %q, esperaba %q", node.Attrs["id"], "myModal")
	}
	if node.Attrs["class"] != "hidden" {
		t.Errorf("Attrs['class'] = %q, esperaba %q", node.Attrs["class"], "hidden")
	}
	if _, ok := node.Attrs["src"]; ok {
		t.Errorf("Attrs['src'] deberia estar eliminado, got %q", node.Attrs["src"])
	}
}

// TestParseHML_IncludeNested verifica que <include> con contenido interno
// se consuma correctamente hasta </include>.
func TestParseHML_IncludeNested(t *testing.T) {
	input := []byte(`<page name="test"><include src="comp.hml"><inner>contenido</inner></include></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Tag != "include" {
		t.Errorf("Tag = %q, esperaba %q", node.Tag, "include")
	}
	if node.IncludeSrc != "comp.hml" {
		t.Errorf("IncludeSrc = %q, esperaba %q", node.IncludeSrc, "comp.hml")
	}
}
