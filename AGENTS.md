# AGENTS.md — dino-hatch

## Descripcion
Hatch es un framework TUI declarativo para Go. Parsea archivos .hml (TermML + HSS)
en un AST, resuelve includes de componentes, calcula layout con flexbox simplificado,
soporta overflow scroll, eventos declarativos, data binding, CSS variables, y
renderiza a terminal usando tcell. No requiere Bubble Tea ni Lip Gloss.

- **Modulo**: `github.com/srdino/dino-hatch`
- **Go version**: 1.24.4
- **Estado**: MVP v0.1.0-alpha

## Pipeline
`.hml` → Parser (HML+HSS) → AST raw → resolveIncludes → ComputeStyles (con CSS vars)
→ Layout (flexbox + min/max clamping + padding/margin) → Render (CellBuffer persistente
+ Diff) → tcell.Show()

## Skills
- **architecture-workspace**: Navegacion de arquitectura del proyecto
- **backend-workspace**: Contexto de backend Go
- **project-conventions**: Convenciones de estructura del proyecto

## Specs activos (agent/specs/)
- **hatch-core** (status: draft, v0.1.0-alpha): Nucleo del framework Hatch. Parseo HML/HSS,
  componentes reutilizables (include), eventos declarativos, data binding, CSS variables,
  scroll, layout flexbox, renderizado tcell. Contiene requirements, design y tasks.

## Designs (agent/specs/hatch-core/design.md)
- **DSG-001 — AST**: Tipos Document, Page, ElementNode, ComputedStyle, BoundBox, StyleRule.
  Parseo HML con encoding/xml, parseo HSS manual con scanner. Campos: id, events, overflow,
  scroll, padding, margin, min/max, theme vars.
- **DSG-002 — Layout Flexbox Simplificado**: Algoritmo de layout con grow, gap, align,
  justify, direction. Padding, margin, min/max clamping. Overflow truncado. ContentHeight cache.
- **DSG-003 — Cell Buffer y Diff**: CellBuffer persistente (PrevCB/CurrCB en AppState),
  diff celda por celda, dirty flag para saltar frames sin cambios.
- **DSG-004 — tcell Screen Driver**: Interfaz Screen (Init, Flush, Close, Size).
  Implementaciones: tcellScreen (real), mockScreen (tests).
- **DSG-005 — Component System**: Resolucion de <include src="..."> antes de ComputeStyles.
  mergeAttrs, merge de estilos desde componentes incluidos, include anidados.
- **DSG-006 — Overflow Scroll Containers**: Per-container scroll (ScrollY interno),
  clipping de hijos con isChildVisible/clipHeight, scrollbar visual en contenedor,
  wheel event routing a scroll container mas especifico.
- **DSG-007 — Declarative Events + Data Binding**: Acciones onclick (page:, modal:,
  action:), bind attribute para sincronizar input→target en 5 puntos de input.go.
- **DSG-008 — CSS Variables / Theming**: ParseCSSVars() extrae :root variables,
  resolveCSSVars() reemplaza var(--name) en applyProps, vars propagadas via
  ComputeStyles(…, vars).
- **DSG-009 — Tab Focus Navigation**: FocusOrder construido por DFS, FocusIndex,
  Tab/Shift+Tab en keyboard.go, cursor en input/textarea al focus.

## Tasks (agent/specs/hatch-core/tasks.md)
| Task | Descripcion | Prioridad | Estado |
|------|-------------|-----------|--------|
| TASK-001 | Inicializar proyecto Go | Alta | ✅ Done |
| TASK-002 | Definir tipos AST | Alta | ✅ Done |
| TASK-003 | Implementar parser HML | Alta | ✅ Done |
| TASK-004 | Implementar parser HSS | Alta | ✅ Done |
| TASK-005 | Merge de estilos (computados) | Alta | ✅ Done |
| TASK-006 | Implementar layout flexbox | Alta | ✅ Done |
| TASK-007 | Implementar cell buffer | Alta | ✅ Done |
| TASK-008 | Interfaz Screen y driver tcell | Alta | ✅ Done |
| TASK-009 | Pipeline completo (main.go) | Alta | ✅ Done |
| TASK-010 | Crear archivo .hml de ejemplo | Baja | ✅ Done |
| TASK-011 | Crear AGENTS.md | Alta | ✅ Done |
| TASK-012 | Configurar Makefile | Media | ✅ Done |
| TASK-013 | Padding, Margin, Min/Max en AST+layout | Alta | ✅ Done |
| TASK-014 | Refactor main.go en modulos (pipeline, eventloop, etc.) | Alta | ✅ Done |
| TASK-015 | CellBuffer persistente + dirty flag | Alta | ✅ Done |
| TASK-016 | Scrollbar visual + throttle mouse | Media | ✅ Done |
| TASK-017 | Component system (include) | Alta | ✅ Done |
| TASK-018 | Overflow scroll containers | Alta | ✅ Done |
| TASK-019 | Eventos declarativos + data binding | Alta | ✅ Done |
| TASK-020 | CSS variables / theming | Media | ✅ Done |
| TASK-021 | Tab focus navigation | Media | ✅ Done |
| TASK-022 | Demo redisenada (7 componentes, 4 paginas) | Media | ✅ Done |

## Convenciones

### Codigo
- **Comentarios**: en espanol, minimos, solo explican "por que", no "que" (REQ-009)
- **Limite de archivos**: maximo 150 lineas por archivo .go. Si se excede, crear subcarpeta
  y dividir (REQ-010)
- **Dependencias**: solo tcell v3 + stdlib. No agregar dependencias externas sin approval (REQ-008)
- **Excepcion 150 lineas**: flex.go (129L), layout.go (213L) — permisibles por ser el core
  algoritmico del layout engine

### Testing
- **Cobertura minima**: 70% en parser, 70% en layout, 70% en render, 70% en ast (REQ-011)
- **Tests obligatorios**: todo paquete debe tener `_test.go` (REQ-011)
- **Mock Screen**: tests de render usan mockScreen, no terminal real (DSG-004)
- **Test fixtures**: archivos .hml de prueba en `internal/parser/testdata/`

### Git
- **Formato de commits**: `[0.1.X] descripcion descriptiva` (REQ-012)
- **Version actual**: v0.1.0-alpha
- **Commits atomicos**: un cambio conceptual por commit

## Estructura del proyecto
```
dino-hatch/
├── AGENTS.md                          ← Este archivo
├── Makefile                           ← Build, test, run, coverage (TASK-012)
├── go.mod / go.sum                    ← Modulo Go (tcell v2.13.10)
├── .gitignore
├── bin/                               ← Binarios compilados (gitignored)
├── canva/
│   ├── demo.hml                       ← Archivo HML principal (TASK-010)
│   ├── handler.json                   ← Bindings de teclado
│   ├── logs/                          ← Logs de handler
│   ├── components/                    ← Componentes reutilizables
│   │   ├── header.hml
│   │   ├── modal.hml
│   │   ├── scroll_demo.hml
│   │   ├── events_demo.hml
│   │   ├── bind_demo.hml
│   │   ├── theme_demo.hml
│   │   ├── textarea_demo.hml
│   │   ├── layout_demo.hml
│   │   └── tabs_demo.hml
│   └── pages/                         ← Paginas de la app
│       ├── inicio.hml
│       ├── proyectos.hml
│       ├── config.hml
│       └── ayuda.hml
├── cmd/hatch/                         ← Entry point CLI (22 archivos)
│   ├── main.go                        ← Entry point (TASK-009): flags→parse→tcell→RunLoop
│   ├── pipeline.go                    ← AppState, LoadPage, NewAppState
│   ├── eventloop.go                   ← RunLoop: dirty flag, cell buf diff, scrollbars
│   ├── state.go                       ← inputState, normalizedSelectionRect
│   ├── render.go                      ← renderDoc, renderPage, renderNode (dispatcher)
│   ├── render_elements.go             ← renderBoxContent, renderButton, renderText, etc.
│   ├── render_helpers.go              ← drawBorder, drawOverlay, renderTextarea, etc.
│   ├── render_clip.go                 ← isChildVisible, clipHeight (overflow clipping)
│   ├── scrollbar.go                   ← drawScrollbar (global), drawContainerScrollbar
│   ├── scroll_container.go            ← findScrollContainer (hit test)
│   ├── textarea_scrollbar.go          ← renderTextareaScrollbar (barra interna)
│   ├── handler.go                     ← executeEvent, executeDataBinding
│   ├── include.go                     ← resolveIncludes, walkAndResolve, mergeAttrs
│   ├── navigate.go                    ← ContentHeight(), FindElementByID()
│   ├── focus.go                       ← buildFocusOrder()
│   ├── input.go                       ← handleInputKey, cursorLineCol, moveCursorUp/Down
│   ├── interactive.go                 ← hitTest, copyToClipboard, getElementText
│   ├── mouse_handler.go               ← handleMouseEvent, handleMousePress
│   ├── mouse_release.go               ← handleMouseRelease, copy-on-release
│   ├── keyboard.go                    ← handleKeyEvent (Tab, arrows, scroll)
│   ├── colors.go                      ← toggleRandomColors, collectBoxes
│   └── diag_test.go                   ← Tests de textarea Enter, shrink, layout stability
├── internal/
│   ├── ast/
│   │   ├── node.go                    ← Tipos AST: Document, Page, ElementNode,
│   │   │                                ComputedStyle, BoundBox, StyleRule, DefaultStyle
│   │   └── node_test.go               ← Tests AST
│   ├── parser/
│   │   ├── hml.go                     ← ParseHML (entry point), preprocessStyleBlocks
│   │   ├── hml_parse.go               ← pageFromStartElement, parseElement, attrsFromSlice
│   │   ├── hss.go                     ← ParseHSS, ParseCSSVars, parseProperties, isKnownProp
│   │   ├── compute.go                 ← ComputeStyles, resolveNode, parseInt/parseFloat/etc.
│   │   ├── props.go                   ← applyProps, resolveCSSVars
│   │   ├── inherit.go                 ← inheritProps (herencia de estilos)
│   │   ├── hml_test.go, hml_error_test.go
│   │   ├── hss_test.go, hss_props_test.go, hss_vars_test.go
│   │   ├── compute_test.go, compute_edge_test.go, compute_helpers_test.go
│   │   ├── props_values_test.go, props_border_test.go, props_invalid_test.go
│   │   ├── props_overflow_test.go, props_vars_test.go
│   │   ├── inherit_test.go, inherit_cascade_test.go, inherit_allprops_test.go
│   │   └── testdata/                  ← Fixtures .hml para tests
│   ├── handler/
│   │   ├── engine.go                  ← Engine, ActionFunc, Binding, Config, Load, Handle
│   │   ├── engine_test.go             ← Tests engine
│   │   ├── keys.go                    ← keyNames tcell→string
│   │   ├── keys_test.go               ← Tests keys
│   │   └── logging_test.go            ← Tests logging
│   ├── layout/
│   │   ├── flex.go                    ← Layout, layoutPage, layoutNode, post-expansion
│   │   ├── layout.go                  ← layoutChildren (algoritmo flexbox principal)
│   │   ├── content.go                 ← ContentHeight, ContentWidth (para scroll)
│   │   ├── measure.go                 ← intrinsicSize
│   │   ├── flex_test.go, flex_basic_test.go, flex_distribution_test.go
│   │   ├── flex_align_test.go, flex_edge_test.go
│   │   ├── helpers_test.go, measure_test.go
│   └── render/
│       ├── cellbuffer.go              ← CellBuffer, Diff, Resize, Fill
│       ├── cellbuffer_test.go         ← Tests cell buffer
│       ├── screen.go                  ← Screen interface + tcell driver
│       └── screen_test.go             ← Tests de render con mockScreen
└── agent/
    └── specs/
        └── hatch-core/
            ├── requirements.md         ← Requerimientos funcionales y no funcionales
            ├── design.md               ← Diseno arquitectonico (DSG-001..009)
            └── tasks.md                ← Plan de tareas (TASK-001..022)
```

## Requerimientos funcionales (hatch-core)
| ID | Descripcion | Prioridad | Estado |
|----|-------------|-----------|--------|
| REQ-001 | Parseo de HML (XML-like a AST) | Alta | ✅ |
| REQ-002 | Parseo de HSS (bloques `<style>` CSS-like) | Alta | ✅ |
| REQ-003 | AST con estilos computados (herencia + defaults) | Alta | ✅ |
| REQ-004 | Layout flexbox (grow, gap, align, justify, direction) | Alta | ✅ |
| REQ-005 | Renderizado a terminal con tcell v3 | Alta | ✅ |
| REQ-006 | Cell buffer con diff (dirty rect) | Media | ✅ |
| REQ-007 | CLI `hatch run <archivo.hml>` | Media | ✅ |
| REQ-013 | Component system (include src) | Alta | ✅ |
| REQ-014 | Overflow scroll en contenedores | Alta | ✅ |
| REQ-015 | Eventos declarativos (onclick, onchange) | Alta | ✅ |
| REQ-016 | Data binding (bind attribute) | Alta | ✅ |
| REQ-017 | CSS variables / themes (:root + var()) | Media | ✅ |
| REQ-018 | Tab focus navigation | Media | ✅ |

## Requerimientos no funcionales (hatch-core)
| ID | Descripcion | Prioridad |
|----|-------------|-----------|
| REQ-008 | Dependencias minimas (solo tcell + stdlib) | Alta |
| REQ-009 | Comentarios en espanol | Media |
| REQ-010 | Limite de 150 lineas por archivo .go | Media |
| REQ-011 | Tests obligatorios, cobertura >=70% en parser, layout, render, ast | Alta |
| REQ-012 | Commits versionados `[0.1.X]` | Media |

## Dependencias
- **tcell v3** (`github.com/gdamore/tcell/v2` v2.13.10): Terminal cell buffer, eventos, resize
- **encoding/xml** (stdlib): Parseo de HML
- **encoding/json** (stdlib): Parseo de handler.json
- **Go 1.24.4**: Version del compilador
