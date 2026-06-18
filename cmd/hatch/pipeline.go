// Estado global de la aplicacion Hatch y funciones de pipeline.
package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/actions"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/events"
	"github.com/srdino/dino-hatch/internal/handler"
	"github.com/srdino/dino-hatch/internal/parser"
	"github.com/srdino/dino-hatch/internal/render"
	"github.com/srdino/dino-hatch/internal/theme"
)

var state *AppState

type AppState struct {
	Doc, ModalDoc     *ast.Document
	StyleRules        []ast.StyleRule
	PageFiles         map[string]string
	CurrentPage       string
	Screen            tcell.Screen
	W, H, ScrollY, MaxScroll int
	PrevButtons       tcell.ButtonMask
	ModalOpen         bool
	HoveredElement    *ast.ElementNode
	FocusedElement    *ast.ElementNode
	InputStates       map[*ast.ElementNode]*inputState
	SelStartX, SelStartY, SelEndX, SelEndY int
	SelActive         bool
	RandomColorsMode  bool
	BoxColors         map[*ast.ElementNode]tcell.Color
	Eng               *handler.Engine
	Running           bool
	ActionsCB         *actions.Callbacks
	PrevCB, CurrCB      *render.CellBuffer // B1: persistente+Diff
	Dirty, LayoutDirty  bool              // dirty flags
	ForceFullRedraw     bool              // forzar push completo (navegacion)
	LastMouseEv       time.Time          // B3: Mouse throttle
	FocusOrder        []*ast.ElementNode // C2: orden Tab (DFS)
	FocusIndex        int                // -1 = ninguno
	BaseDir, FilePath string             // directorio base y ruta .hml
	EventBus          *events.Bus
	ThemeManager      *theme.Manager // Gestor de temas para cambio dinamico
}

func NewAppState(doc *ast.Document, rules []ast.StyleRule, tcellScr tcell.Screen, baseDir string) *AppState {
	var modalDoc *ast.Document
	if mData, err := os.ReadFile("canva/components/modal.hml"); err == nil {
		if md, mStyleBlocks, parseErr := parser.ParseHML(mData); parseErr == nil {
			for _, block := range mStyleBlocks {
				if r, hssErr := parser.ParseHSS(block); hssErr == nil {
					modalDoc = parser.ComputeStyles(md, r, nil)
				}
			}
		}
	} else {
		log.Printf("Warning: no se pudo cargar modal.hml: %v", err)
	}
	pageFiles := map[string]string{
		"inicio": "canva/pages/inicio.hml", "proyectos": "canva/pages/proyectos.hml",
		"config": "canva/pages/config.hml", "ayuda": "canva/pages/ayuda.hml",
	}
	w, h := tcellScr.Size()
	s := &AppState{
		Doc: doc, StyleRules: rules, ModalDoc: modalDoc,
		PageFiles: pageFiles, CurrentPage: "inicio",
		Screen: tcellScr, W: w, H: h,
		InputStates:  make(map[*ast.ElementNode]*inputState),
		FocusOrder:   []*ast.ElementNode{}, FocusIndex: -1,
		BaseDir: baseDir, Running: true,
		EventBus: events.New(), ThemeManager: theme.New(),
	}
	state = s
	initActions(s)
	eng := handler.New()
	eng.Register("quit", func() error { s.Running = false; return nil })
	eng.Register("random_colors", func() error { toggleRandomColors(s); return nil })
	if err := eng.Load("canva/handler.json"); err != nil {
		log.Printf("Warning: handler.json: %v", err)
	}
	s.Eng = eng
	s.LoadPage("inicio")
	s.EventBus.Publish("app.state", "loaded")
	s.CurrCB, s.PrevCB = render.NewCellBuffer(w, h), render.NewCellBuffer(w, h)
	s.Dirty, s.LayoutDirty = true, true
	return s
}

func (s *AppState) LoadPage(name string) {
	path, ok := s.PageFiles[name]
	if !ok { return }
	s.FilePath = path
	data, err := os.ReadFile(path)
	if err != nil { log.Printf("Warning: no se pudo cargar %s: %v", path, err); return }
	pageDoc, styleBlocks, err := parser.ParseHML(data)
	if err != nil { log.Printf("Warning: error parseando %s: %v", path, err); return }
	if len(pageDoc.Pages) == 0 || len(pageDoc.Pages[0].Children) == 0 { return }

	if err := s.resolveIncludes(pageDoc); err != nil {
		log.Printf("Warning: error resolviendo includes: %v", err)
	}
	contentArea := s.FindElementByID("content-area")
	if contentArea == nil { log.Printf("Warning: no se encontro content-area"); return }

	// Extraer CSS variables y guardar tema default
	var all string
	for _, b := range styleBlocks { all += b + "\n" }
	if vars := parser.ParseCSSVars(all); len(vars) > 0 {
		if s.Doc.ThemeVars == nil {
			s.Doc.ThemeVars = vars
		} else {
			for k, v := range vars { s.Doc.ThemeVars[k] = v }
		}
		s.ThemeManager.AddTheme("default", s.Doc.ThemeVars)
	}

	mergedRules := append([]ast.StyleRule{}, s.StyleRules...)
	for _, block := range styleBlocks {
		if r, err := parser.ParseHSS(block); err == nil {
			mergedRules = append(mergedRules, r...)
		}
	}
	contentArea.Children = pageDoc.Pages[0].Children[0].Children
	parser.ComputeStyles(s.Doc, mergedRules, s.Doc.ThemeVars)
	s.buildFocusOrder()
	s.Screen.Clear()
	s.Dirty, s.LayoutDirty = true, true
	s.ForceFullRedraw = true // Forzar push completo sin alocar buffer extra (OPT-001)
}

