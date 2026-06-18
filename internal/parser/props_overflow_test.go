package parser

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func TestApplyProps_OverflowHidden(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"overflow": "hidden"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Overflow != "hidden" {
		t.Errorf("overflow=hidden: got %q", s.Overflow)
	}
}

func TestApplyProps_OverflowVisible(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag:   "box",
			Attrs: map[string]string{"overflow": "visible"},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	s := result.Pages[0].Children[0].Style
	if s.Overflow != "visible" {
		t.Errorf("overflow=visible: got %q", s.Overflow)
	}
}
