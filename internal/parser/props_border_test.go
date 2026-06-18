package parser

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestApplyProps_BorderTrue(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "true"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if !s.Border {
		t.Errorf("border=true: got false, want true")
	}
}

func TestApplyProps_BorderYes(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "yes"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if !s.Border {
		t.Errorf("border=yes: got false, want true")
	}
}

func TestApplyProps_BorderOne(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "1"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if !s.Border {
		t.Errorf("border=1: got false, want true")
	}
}

func TestApplyProps_BorderFalse(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "false"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Border {
		t.Errorf("border=false: got true, want false")
	}
}

func TestApplyProps_BorderZero(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"border": "0"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Border {
		t.Errorf("border=0: got true, want false")
	}
}
