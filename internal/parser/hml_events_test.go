package parser

import (
	"testing"
)

// TestParseHML_Onclick verifica que onclick se extraiga como Events["click"]
// y se elimine de Attrs.
func TestParseHML_Onclick(t *testing.T) {
	input := []byte(`<page name="test"><box onclick="page:next" /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Events["click"] != "page:next" {
		t.Errorf("Events['click'] = %q, esperaba %q", node.Events["click"], "page:next")
	}
	if _, ok := node.Attrs["onclick"]; ok {
		t.Errorf("Attrs['onclick'] deberia estar eliminado, got %q", node.Attrs["onclick"])
	}
}

// TestParseHML_MultipleEvents verifica que onclick + onchange + onfocus
// se extraigan correctamente en un mismo nodo.
func TestParseHML_MultipleEvents(t *testing.T) {
	input := []byte(`<page name="test"><input onclick="action:log" onchange="action:save" onfocus="action:focus" /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Events["click"] != "action:log" {
		t.Errorf("Events['click'] = %q, esperaba %q", node.Events["click"], "action:log")
	}
	if node.Events["change"] != "action:save" {
		t.Errorf("Events['change'] = %q, esperaba %q", node.Events["change"], "action:save")
	}
	if node.Events["focus"] != "action:focus" {
		t.Errorf("Events['focus'] = %q, esperaba %q", node.Events["focus"], "action:focus")
	}
	if len(node.Events) != 3 {
		t.Errorf("len(Events) = %d, esperaba 3", len(node.Events))
	}
}

// TestParseHML_NoEvents verifica que un nodo sin eventos tenga Events == nil.
func TestParseHML_NoEvents(t *testing.T) {
	input := []byte(`<page name="test"><box id="simple" /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Events != nil {
		t.Errorf("Events deberia ser nil, got %v", node.Events)
	}
}

// TestParseHML_Onblur verifica que onblur se extraiga como Events["blur"].
func TestParseHML_Onblur(t *testing.T) {
	input := []byte(`<page name="test"><textarea onblur="action:validate" /></page>`)
	doc, _, err := ParseHML(input)
	if err != nil {
		t.Fatalf("ParseHML fallo: %v", err)
	}

	node := doc.Pages[0].Children[0]
	if node.Events["blur"] != "action:validate" {
		t.Errorf("Events['blur'] = %q, esperaba %q", node.Events["blur"], "action:validate")
	}
}
