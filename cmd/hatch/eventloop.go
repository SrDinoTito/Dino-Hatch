// Event loop principal: layout, render, flush y dispatch (teclado + resize).
// B1: CellBuffer persistente + Diff(). B2: Dirty flag. B3: Mouse throttle. B4: ContentHeight cache.
package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/layout"
	"github.com/srdino/dino-hatch/internal/render"
)

// RunLoop ejecuta el ciclo principal: layout, render, flush y dispatch de eventos.
// Usa dirty flag (B2) para saltar layout/render cuando no hay cambios,
// CellBuffer persistente con Diff (B1) para enviar solo celdas modificadas,
// y frame budget tracking para evitar acumulacion de frames.
func RunLoop(s *AppState) {
	var prevRenderDuration time.Duration
	const frameBudget = 16 * time.Millisecond // ~60 FPS

	for s.Running {
		// CRITICO: Procesar eventos del bus SIEMPRE, incluso si Dirty=false.
		// Si solo se procesara dentro del bloque dirty, background goroutines
		// que publican eventos via PublishEvent se perderian silenciosamente.
		s.EventBus.ProcessAll()

		// ===== SOLO SI DIRTY =====
		if s.Dirty {
			// Frame budget: si el render anterior tomo >2x budget, saltamos
			// este frame para no acumular retraso (dropeo de frames controlado).
			if prevRenderDuration > frameBudget*2 {
				s.Dirty = false
				prevRenderDuration = 0
				continue
			}
			frameStart := time.Now()

			if s.LayoutDirty {
				layout.Layout(s.Doc, s.W, s.H)
				s.MaxScroll = max(0, s.ContentHeight()-s.H)
				if s.ScrollY > s.MaxScroll {
					s.ScrollY = s.MaxScroll
				}
				s.LayoutDirty = false
			}

			if s.CurrCB.Width() != s.W || s.CurrCB.Height() != s.H {
				s.CurrCB.Resize(s.W, s.H)
			}
			s.CurrCB.Fill(' ', tcell.StyleDefault)
			renderDoc(s.CurrCB, s.Doc, s.ScrollY)
			if s.ModalOpen && s.ModalDoc != nil {
				layout.Layout(s.ModalDoc, s.W, s.H)
				drawOverlay(s.CurrCB)
				if len(s.ModalDoc.Pages) > 0 && len(s.ModalDoc.Pages[0].Children) > 0 {
					bg := &s.ModalDoc.Pages[0].Children[0]
					for i := range bg.Children {
						renderNode(s.CurrCB, &bg.Children[i], tcell.ColorReset, 0)
					}
				}
			}

			// C1: Scrollbar visual — ultima columna cuando hay scroll disponible
			drawScrollbar(s.CurrCB, s.H, s.ScrollY, s.MaxScroll)

			// D3: Scrollbars de contenedores con overflow:scroll
			drawAllContainerScrollbars(s.CurrCB, s.Doc, s.ScrollY)

			// B1: Diff o push segun corresponda
			if s.ForceFullRedraw || s.PrevCB.Width() != s.W || s.PrevCB.Height() != s.H {
				// full push (navegacion o resize/primer frame)
				if s.ForceFullRedraw {
					s.ForceFullRedraw = false
				}
				if s.PrevCB.Width() != s.W || s.PrevCB.Height() != s.H {
					s.PrevCB.Resize(s.W, s.H)
				}
				cells := s.CurrCB.Cells()
				for i, cell := range cells {
					x := i % s.W
					y := i / s.W
					s.Screen.SetContent(x, y, cell.Rune, nil, cell.Style)
				}
			} else {
				// Diff: solo celdas cambiadas
				_, updates := s.PrevCB.Diff(s.CurrCB)
				for _, u := range updates {
					s.Screen.SetContent(u.X, u.Y, u.Cell.Rune, nil, u.Cell.Style)
				}
			}
			s.PrevCB, s.CurrCB = s.CurrCB, s.PrevCB
			s.Screen.Show()
			s.Dirty = false

			prevRenderDuration = time.Since(frameStart)
		}

		// ===== SIEMPRE: esperar evento (bloqueante) =====
		ev := s.Screen.PollEvent()

		switch e := ev.(type) {
		case *tcell.EventMouse:
			// B3: actualizar estado, decidir si re-renderizar.
			// Solo se marca Dirty si hay clic, wheel o cambio de hover.
			// Sin throttle temporal: renderNode no usa hover para estilos,
			// por lo que refrescar cada 33ms sin cambio visual es desperdicio.
			prevHovered := s.HoveredElement
			isButtonEvent := e.Buttons() != tcell.ButtonNone &&
				e.Buttons() != tcell.WheelUp && e.Buttons() != tcell.WheelDown
			isWheel := e.Buttons() == tcell.WheelUp || e.Buttons() == tcell.WheelDown

			handleMouseEvent(s, e)

			if isButtonEvent || isWheel {
				s.Dirty = true
			} else if s.HoveredElement != prevHovered {
				s.Dirty = true
			}
			s.LastMouseEv = time.Now()

		case *tcell.EventResize:
			s.W, s.H = e.Size()
			s.PrevCB = render.NewCellBuffer(s.W, s.H)
			s.CurrCB = render.NewCellBuffer(s.W, s.H)
			s.Dirty = true

		case *tcell.EventKey:
			handleKeyEvent(s, e)
			s.Dirty = true
		}
	}
}
