# dino-hatch 🦕

**Framework TUI declarativo para Go.**  
Parsea archivos `.hml` (TermML + HSS) en un AST, calcula layout con flexbox simplificado, y renderiza a terminal usando [tcell v2](https://github.com/gdamore/tcell). Sin Bubble Tea, sin Lip Gloss — una dependencia externa.

```
Versión: [0.1.1]
Go:      1.24.4
Módulo:  github.com/srdino/dino-hatch
```

📖 **Glosario completo**: [`docs/glossary.md`](docs/glossary.md) — referencia de elementos HML, propiedades HSS, eventos, acciones y tipos.

---

## Pipeline

```
.hml → Parser (HML+HSS) → AST raw → resolveIncludes → ComputeStyles (con CSS vars)
→ Layout (flexbox + min/max clamping + padding/margin)
→ Render (CellBuffer persistente + Diff) → tcell.Show()
```

---

## Features

### Core (estables)

| Feature | Descripción |
|---------|-------------|
| **Parseo HML** | XML-like a AST (`Document`, `Page`, `ElementNode`) |
| **Parseo HSS** | CSS-like en bloques `<style>`, incluyendo `:root` |
| **Component system** | `<include src="...">` con merge de atributos y estilos, includes anidados |
| **Layout flexbox** | `direction`, `grow`, `gap`, `align`, `justify`, `padding`, `margin`, `min-width/height`, `max-width/height` |
| **Overflow scroll** | Per-container `overflow: scroll` con scrollbar visual y clipping de hijos |
| **CellBuffer + Diff** | Buffer persistente, solo envía celdas cambiadas a tcell, dirty flag |
| **CSS variables / Theming** | `:root { --x: y }` + `var(--x)` en estilos, multi-tema dinámico |

### Interactividad

| Feature | Descripción |
|---------|-------------|
| **Eventos declarativos** | `onclick`, `onchange`, `onfocus`, `onblur` en elementos |
| **Data binding** | `<input bind="target-id">` sincroniza texto en vivo entre elementos |
| **Tab focus navigation** | Navegación Tab/Shift+Tab entre input, textarea y button |
| **Scroll** | Por página (PgUp/PgDn, flechas) y por contenedor (mouse wheel con throttle 30fps) |
| **Modal overlay** | Sistema de modales con apertura/cierre desde eventos |

### Conectividad externa (experimental)

| Feature | Descripción |
|---------|-------------|
| **`exec:` actions** | `onclick="exec:ls -la"` ejecuta comandos, stdout a `#exec-log` |
| **`curl:` actions** | `onclick="curl:https://api.example.com"` peticiones HTTP GET |
| **`theme:` actions** | `onclick="theme:dark"` cambia tema dinámicamente |
| **Event bus** | Pub/Sub thread-safe para integración con goroutines background |
| **Stdin pipe** | `echo "data" \| hatch demo.hml` publica líneas como eventos en el bus |

### Optimizaciones

- **Dirty flag + LayoutDirty**: layout solo cuando hay cambios reales
- **Mouse throttle**: eventos mouse a 30fps, no fuerzan re-layout
- **Frame budget tracking**: salta frames lentos para evitar acumulación de retraso
- **ForceFullRedraw**: zero allocations en navegación entre páginas
- **ContentHeight cache**: altura de contenido cacheadas por nodo
- **Scrollbar eficiente**: solo se renderiza cuando cambia el scroll offset

---

## Cómo usar

### Build

```bash
go build -o bin/hatch ./cmd/hatch
```

### Run

```bash
./bin/hatch run canva/demo.hml
```

### Atajos de teclado

| Tecla | Acción |
|-------|--------|
| Tab / Shift+Tab | Navegar entre inputs/buttons/textareas |
| ↑ ↓ ← → | Scroll |
| PgUp / PgDn | Scroll rápido |
| Home / End | Ir al inicio / final |
| Shift+D | Debug colors mode |
| ESC | Cerrar modal |

---

## Demo incluida

La demo en `canva/demo.hml` incluye:

- **7 componentes demo**: scroll, events, bind, theme, textarea, layout, tabs
- **4 páginas**: inicio, proyectos, config, ayuda
- Cada componente documenta qué prueba y el resultado esperado

Componentes reutilizables en `canva/components/`, páginas en `canva/pages/`.

---

## Estructura del proyecto

```
dino-hatch/
├── AGENTS.md              ← Documentación del agente (arquitectura, specs, tareas)
├── Makefile               ← Build, test, run, coverage
├── go.mod / go.sum        ← Módulo Go (tcell v2.13.10 como única dep externa)
├── cmd/hatch/             ← Entry point CLI (22 archivos, ~1900 líneas)
│   ├── main.go            ← Entry point: flags → parse → tcell → RunLoop
│   ├── pipeline.go        ← AppState, LoadPage, NewAppState
│   ├── eventloop.go       ← RunLoop: dirty flag, cell buf diff, scrollbars
│   ├── render.go          ← renderDoc, renderPage, renderNode (dispatcher)
│   ├── render_elements.go ← renderBoxContent, renderButton, renderText, etc.
│   ├── render_helpers.go  ← drawBorder, drawOverlay, renderTextarea, etc.
│   ├── render_clip.go     ← isChildVisible, clipHeight (overflow clipping)
│   ├── scrollbar.go       ← drawScrollbar (global), drawContainerScrollbar
│   ├── scroll_container.go← findScrollContainer (hit test)
│   ├── include.go         ← resolveIncludes, walkAndResolve, mergeAttrs
│   ├── navigate.go        ← ContentHeight(), FindElementByID()
│   ├── focus.go           ← buildFocusOrder()
│   ├── input.go           ← handleInputKey, cursorLineCol, moveCursorUp/Down
│   ├── interactive.go     ← hitTest, copyToClipboard, getElementText
│   ├── handler.go         ← executeEvent, executeDataBinding
│   ├── state.go           ← inputState, normalizedSelectionRect
│   ├── keyboard.go        ← handleKeyEvent (Tab, arrows, scroll)
│   ├── mouse_handler.go   ← handleMouseEvent, handleMousePress
│   ├── mouse_release.go   ← handleMouseRelease, copy-on-release
│   ├── colors.go          ← toggleRandomColors, collectBoxes
│   ├── textarea_scrollbar.go ← renderTextareaScrollbar
│   └── diag_test.go       ← Tests de textarea Enter, shrink, layout stability
├── internal/
│   ├── ast/               ← Tipos AST (node.go, node_test.go)
│   │   ├── node.go        ← Document, Page, ElementNode, ComputedStyle, BoundBox
│   │   └── node_test.go   ← Tests AST (100% cobertura)
│   ├── parser/            ← Parser HML+HSS, ComputeStyles, CSS vars
│   │   ├── hml.go         ← ParseHML, preprocessStyleBlocks
│   │   ├── hml_parse.go   ← pageFromStartElement, parseElement, attrsFromSlice
│   │   ├── hss.go         ← ParseHSS, ParseCSSVars, parseProperties
│   │   ├── compute.go     ← ComputeStyles, resolveNode
│   │   ├── props.go       ← applyProps, resolveCSSVars
│   │   ├── inherit.go     ← inheritProps (herencia de estilos)
│   │   ├── ..._test.go    ← ~12 archivos de test (88% cobertura)
│   │   └── testdata/      ← Fixtures .hml para tests
│   ├── layout/            ← Layout flexbox
│   │   ├── flex.go        ← Layout, layoutPage, layoutNode, post-expansion
│   │   ├── layout.go      ← layoutChildren (algoritmo flexbox principal)
│   │   ├── content.go     ← ContentHeight, ContentWidth (para scroll)
│   │   ├── measure.go     ← intrinsicSize
│   │   └── ..._test.go    ← ~8 archivos de test (78% cobertura)
│   ├── render/            ← CellBuffer, Screen driver tcell
│   │   ├── cellbuffer.go  ← CellBuffer, Diff, Resize, Fill
│   │   ├── screen.go      ← Screen interface + tcell driver
│   │   └── ..._test.go    ← Tests con mockScreen (73% cobertura)
│   ├── handler/           ← Key binding engine
│   │   ├── engine.go      ← Engine, ActionFunc, Binding, Config, Load, Handle
│   │   └── ..._test.go    ← Tests engine + keys + logging
│   ├── actions/           ← Sistema extensible de acciones
│   │   ├── actions.go     ← Tipos base: Callbacks, Handler, Registry
│   │   ├── exec.go        ← HandlerExec (shell), HandlerCurl (HTTP GET)
│   │   ├── theme.go       ← HandlerTheme (cambio dinámico de tema)
│   │   └── actions_test.go← Tests (98% cobertura)
│   ├── events/            ← Event bus Pub/Sub thread-safe
│   │   ├── bus.go         ← Bus, Subscribe, Publish, Consume
│   │   └── bus_test.go    ← Tests (100% cobertura)
│   └── theme/             ← Theme manager multi-tema
│       ├── manager.go     ← Manager, AddTheme, SwitchTheme
│       └── manager_test.go← Tests (100% cobertura)
├── docs/                 ← Documentación (glosario, referencias)
│   └── glossary.md       ← Glosario completo de features
├── canva/                 ← Demo
│   ├── demo.hml           ← Archivo HML principal
│   ├── handler.json       ← Bindings de teclado
│   ├── components/        ← Componentes reutilizables (9 archivos .hml)
│   └── pages/             ← Páginas de la app (4 archivos .hml)
└── agent/specs/           ← SDD specs (requirements, design, tasks)
    └── hatch-core/
        ├── requirements.md
        ├── design.md
        └── tasks.md
```

---

## Tests

```bash
go test -count=1 ./...
```

### Cobertura actual

| Paquete | Cobertura |
|---------|-----------|
| `internal/ast` | 100% |
| `internal/actions` | 100% |
| `internal/events` | 100% |
| `internal/theme` | 100% |
| `internal/parser` | 92.8% |
| `internal/layout` | 90.5% |
| `internal/render` | 79.1% |

Ejecutar con cobertura:

```bash
go test -count=1 -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Convenciones del proyecto

- **Dependencias mínimas**: solo tcell v2 + stdlib. No agregar sin aprobación explícita.
- **Comentarios**: en español, explican el "por qué", no el "qué".
- **Límite de 150 líneas** por archivo `.go`. Excepción: `flex.go` y `layout.go` (core algorítmico).
- **Commits**: formato `[0.1.X] descripción descriptiva`.
- **Tests obligatorios**: todo paquete debe tener `_test.go`.

---

## Licencia

MIT
