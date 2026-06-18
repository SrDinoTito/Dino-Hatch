package parser

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

func TestInherit_OverflowCascade(t *testing.T) {
	// abuelo -> padre -> hijo: overflow debe propagarse
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "grandparent",
			Attrs: map[string]string{
				"overflow": "hidden",
			},
			Children: []ast.ElementNode{{
				Tag: "parent",
				Children: []ast.ElementNode{{
					Tag: "child",
				}},
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	gp := result.Pages[0].Children[0]
	p := gp.Children[0]
	c := p.Children[0]

	if gp.Style.Overflow != "hidden" {
		t.Errorf("abuelo overflow: got %q, want hidden", gp.Style.Overflow)
	}
	if p.Style.Overflow != "hidden" {
		t.Errorf("padre deberia heredar hidden: got %q", p.Style.Overflow)
	}
	if c.Style.Overflow != "hidden" {
		t.Errorf("hijo deberia heredar hidden: got %q", c.Style.Overflow)
	}
}

func TestInherit_OverflowCascadeVisible(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "grandparent",
			Children: []ast.ElementNode{{
				Tag:   "parent",
				Attrs: map[string]string{"overflow": "scroll"},
				Children: []ast.ElementNode{{
					Tag: "child",
				}},
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	gp := result.Pages[0].Children[0]
	p := gp.Children[0]
	c := p.Children[0]

	if gp.Style.Overflow != "visible" {
		t.Errorf("abuelo overflow: got %q, want visible (default)", gp.Style.Overflow)
	}
	if p.Style.Overflow != "scroll" {
		t.Errorf("padre overflow: got %q, want scroll", p.Style.Overflow)
	}
	if c.Style.Overflow != "scroll" {
		t.Errorf("hijo deberia heredar scroll: got %q", c.Style.Overflow)
	}
}

func TestInherit_ColorAndBgWithOverflow(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "parent",
			Attrs: map[string]string{
				"color":    "red",
				"bg":       "blue",
				"overflow": "scroll",
			},
			Children: []ast.ElementNode{{
				Tag: "child",
				Attrs: map[string]string{
					"overflow": "hidden",
				},
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	p := result.Pages[0].Children[0]
	c := p.Children[0]

	if c.Style.Color != tcell.GetColor("red") {
		t.Errorf("hijo deberia heredar color: got %v, want red", c.Style.Color)
	}
	if c.Style.BgColor != tcell.GetColor("blue") {
		t.Errorf("hijo deberia heredar bg: got %v, want blue", c.Style.BgColor)
	}
	if c.Style.Overflow != "hidden" {
		t.Errorf("hijo overflow propio: got %q, want hidden", c.Style.Overflow)
	}
}

func TestInherit_ChildDefinedPropsNotInherited(t *testing.T) {
	doc := &ast.Document{
		Pages: []ast.Page{{Children: []ast.ElementNode{{
			Tag: "parent",
			Attrs: map[string]string{
				"color": "red",
			},
			Children: []ast.ElementNode{{
				Tag:   "child",
				Attrs: map[string]string{"color": "blue"},
			}},
		}}}},
	}
	result := ComputeStyles(doc, nil, nil)
	c := result.Pages[0].Children[0].Children[0]
	if c.Style.Color != tcell.GetColor("blue") {
		t.Errorf("hijo deberia tener su propio color: got %v, want blue", c.Style.Color)
	}
}
