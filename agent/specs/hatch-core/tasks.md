# HATCH-CORE — Tasks

## Version
0.1.0-alpha

---

## TASK-001 — Inicializar proyecto Go
**Prioridad**: Alta
**Dependencias**: Ninguna
**Archivos afectados**: go.mod, .gitignore, Makefile
**Riesgo**: Bajo
**Validación**: `go build ./...` compila sin errores
**Criterio de done**: Proyecto Go inicializado con dependencia tcell v3
**Referencias**: REQ-008
**Estado**: ✅ Done

## TASK-002 — Definir tipos AST
**Prioridad**: Alta
**Dependencias**: Ninguna
**Archivos afectados**: internal/ast/node.go
**Riesgo**: Bajo
**Validación**: Los tipos compilan y representan el diseño
**Criterio de done**: Tipos Document, Page, ElementNode, ComputedStyle, BoundBox, StyleRule,
DefaultStyle definidos con todos los campos necesarios (id, events, overflow, scroll,
padding, margin, min/max, include, theme vars)
**Referencias**: REQ-001, REQ-003, DSG-001
**Estado**: ✅ Done

## TASK-003 — Implementar parser HML
**Prioridad**: Alta
**Dependencias**: TASK-002
**Archivos afectados**: internal/parser/hml.go, internal/parser/hml_parse.go,
internal/parser/hml_test.go, internal/parser/hml_error_test.go
**Riesgo**: Medio — encoding/xml puede no ser ideal para HML
**Validación**: Tests pasan con archivos .hml válidos e inválidos
**Criterio de done**: Parseo de `<page>`, `<box>`, `<text>`, `<input>`, `<textarea>`,
`<button>`, `<include>`, `<style>`. Extracción de eventos (onclick) de Attrs a Events.
Soporte de IncludeSrc.
**Referencias**: REQ-001, DSG-001
**Estado**: ✅ Done

## TASK-004 — Implementar parser HSS
**Prioridad**: Alta
**Dependencias**: TASK-002
**Archivos afectados**: internal/parser/hss.go, internal/parser/hss_test.go,
internal/parser/hss_props_test.go, internal/parser/hss_vars_test.go
**Riesgo**: Medio — parseo manual de CSS-like
**Validación**: Tests pasan con reglas HSS válidas
**Criterio de done**: Bloques `<style>` se parsean a StyleRules. Propiedades inválidas
se ignoran con warning. ParseCSSVars extrae variables :root. Selector :root no genera
StyleRule.
**Referencias**: REQ-002, REQ-017, DSG-001, DSG-008
**Estado**: ✅ Done

## TASK-005 — Merge de estilos (computados)
**Prioridad**: Alta
**Dependencias**: TASK-003, TASK-004
**Archivos afectados**: internal/parser/compute.go, internal/parser/props.go,
internal/parser/inherit.go, tests varios
**Riesgo**: Bajo
**Validación**: Tests de herencia y especificidad de selector
**Criterio de done**: ComputeStyles recibe vars map[string]string. applyProps soporta
overflow, padding, margin, min/max. resolveCSSVars reemplaza var(--name). inheritProps
hereda overflow, text-align, valign.
**Referencias**: REQ-003, REQ-017, DSG-001, DSG-008
**Estado**: ✅ Done

## TASK-006 — Implementar layout flexbox
**Prioridad**: Alta
**Dependencias**: TASK-005
**Archivos afectados**: internal/layout/flex.go, internal/layout/layout.go,
internal/layout/measure.go, internal/layout/content.go, tests varios
**Riesgo**: Alto — algoritmo de layout con overflow, grow, gap, padding, margin, min/max
**Validación**: Tests con direction=row y direction=column. Casos de 1, 2, N hijos,
grow mixto, gap, padding, margin, min/max clamping
**Criterio de done**: Layout calcula BoundBox correcto para cada nodo. Overflow truncado.
Post-expansion de contenedores sin width/height. ContentHeight/ContentWidth para scroll.
**Referencias**: REQ-004, DSG-002
**Estado**: ✅ Done

## TASK-007 — Implementar cell buffer
**Prioridad**: Alta
**Dependencias**: TASK-002
**Archivos afectados**: internal/render/cellbuffer.go
**Riesgo**: Medio — rendimiento de diff
**Validación**: Cell buffer almacena y recupera celdas correctamente
**Criterio de done**: Set(x, y, rune, style), Get(x, y), Diff() produce lista de cambios,
Fill(), Resize().
**Referencias**: REQ-006, DSG-003
**Estado**: ✅ Done

## TASK-008 — Implementar interfaz Screen y driver tcell
**Prioridad**: Alta
**Dependencias**: TASK-007
**Archivos afectados**: internal/render/screen.go, internal/render/screen_test.go
**Riesgo**: Bajo — tcell API bien documentada
**Validación**: mockScreen funciona en tests, tcellScreen pinta en terminal real
**Criterio de done**: Interfaz Screen con Init, Flush, Close, Size. Implementación
tcell y mock.
**Referencias**: REQ-005, DSG-004
**Estado**: ✅ Done

## TASK-009 — Pipeline completo (main.go)
**Prioridad**: Alta
**Dependencias**: TASK-003..008
**Archivos afectados**: cmd/hatch/main.go
**Riesgo**: Medio — integración de todos los módulos
**Validación**: `hatch run canva/demo.hml` funciona
**Criterio de done**: Pipeline completo: flags → parse → tcell → RunLoop con event
dispatch, dirty flag, cell buffer diff.
**Referencias**: REQ-007
**Estado**: ✅ Done

## TASK-010 — Crear archivo .hml de ejemplo
**Prioridad**: Baja
**Dependencias**: Ninguna
**Archivos afectados**: canva/demo.hml
**Riesgo**: Bajo
**Validación**: Archivo HML válido que pueda ser parseado por TASK-003
**Criterio de done**: Archivo demo.hml con page, box, text, style.
**Referencias**: REQ-001
**Estado**: ✅ Done

## TASK-011 — Crear AGENTS.md
**Prioridad**: Alta
**Dependencias**: Ninguna
**Archivos afectados**: AGENTS.md
**Riesgo**: Bajo
**Validación**: Sigue el formato estándar de project-conventions
**Criterio de done**: AGENTS.md completo con descripción, skills, specs, reglas de
testing y convenciones.
**Referencias**: REQ-009, REQ-010, REQ-011, REQ-012
**Estado**: ✅ Done

## TASK-012 — Configurar Makefile
**Prioridad**: Media
**Dependencias**: TASK-001
**Archivos afectados**: Makefile
**Riesgo**: Bajo
**Validación**: `make build`, `make test`, `make run` funcionan
**Criterio de done**: Makefile con targets build, test, lint, run, clean, coverage.
**Referencias**: REQ-008
**Estado**: ✅ Done

## TASK-013 — Padding, Margin, Min/Max en AST+layout
**Prioridad**: Alta
**Dependencias**: TASK-005, TASK-006
**Archivos afectados**: internal/ast/node.go, internal/parser/props.go,
internal/layout/flex.go, internal/layout/layout.go
**Riesgo**: Medio — cambios en layout afectan todos los renders
**Validación**: Tests de padding/margin/min/max en compute y layout
**Criterio de done**: ComputedStyle tiene Padding, Margin, MinWidth, MinHeight,
MaxWidth, MaxHeight. applyProps los parsea. layout aplica padding como reducción
del área de contenido, margin como espacio entre hermanos, y min/max clamping
después de grow.
**Referencias**: REQ-003, REQ-004, DSG-001, DSG-002
**Estado**: ✅ Done

## TASK-014 — Refactor main.go en módulos
**Prioridad**: Alta
**Dependencias**: TASK-009
**Archivos afectados**: cmd/hatch/main.go → cmd/hatch/pipeline.go, eventloop.go,
mouse_handler.go, mouse_release.go, keyboard.go, state.go
**Riesgo**: Medio — refactor de código existente
**Validación**: `hatch run canva/demo.hml` sigue funcionando idéntico
**Criterio de done**: main.go reducido de 514L a ~70L. AppState centralizado en
pipeline.go. RunLoop en eventloop.go. Mouse y teclado en archivos separados.
**Referencias**: REQ-010
**Estado**: ✅ Done

## TASK-015 — CellBuffer persistente + dirty flag
**Prioridad**: Alta
**Dependencias**: TASK-007, TASK-014
**Archivos afectados**: cmd/hatch/pipeline.go, cmd/hatch/eventloop.go,
internal/render/cellbuffer.go
**Riesgo**: Bajo — cambio localizado en event loop
**Validación**: Frames sin cambios no ejecutan layout/render. Frames con cambios
solo actualizan celdas modificadas.
**Criterio de done**: AppState.PrevCB y CurrCB persistentes. Diff() en vez de
push completo. Dirty flag salta layout+render cuando false.
**Referencias**: REQ-006, DSG-003
**Estado**: ✅ Done

## TASK-016 — Scrollbar visual + throttle mouse
**Prioridad**: Media
**Dependencias**: TASK-014
**Archivos afectados**: cmd/hatch/eventloop.go, cmd/hatch/scrollbar.go
**Riesgo**: Bajo
**Validación**: Scrollbar aparece cuando MaxScroll > 0. Botón de rueda marca Dirty.
**Criterio de done**: drawScrollbar en última columna, barra proporcional. Mouse
throttle a 33ms (≈30 FPS) para mouse-move en idle.
**Referencias**: DSG-003
**Estado**: ✅ Done

## TASK-017 — Component system (include)
**Prioridad**: Alta
**Dependencias**: TASK-005, TASK-014
**Archivos afectados**: cmd/hatch/include.go, cmd/hatch/pipeline.go,
internal/parser/hml_parse.go
**Riesgo**: Medio — resolución recursiva de includes
**Validación**: Componente con <include src="..."> se resuelve correctamente.
Include anidado funciona.
**Criterio de done**: resolveIncludes reemplaza <include> por AST del componente.
mergeAttrs. Merge de estilos. Includes anidados.
**Referencias**: REQ-013, DSG-005
**Estado**: ✅ Done

## TASK-018 — Overflow scroll containers
**Prioridad**: Alta
**Dependencias**: TASK-006, TASK-014, TASK-017
**Archivos afectados**: cmd/hatch/render_clip.go, cmd/hatch/scroll_container.go,
cmd/hatch/scrollbar.go, cmd/hatch/render.go, cmd/hatch/mouse_handler.go,
internal/layout/content.go, internal/ast/node.go, internal/parser/props.go
**Riesgo**: Medio — clipping de hijos afecta render de todos los contenedores
**Validación**: Contenedor overflow:scroll scrollea independientemente. Hijos se
clipean al viewport. Scrollbar visual aparece.
**Criterio de done**: isChildVisible, clipHeight, findScrollContainer,
drawContainerScrollbar. Wheel routing a scroll container más específico.
**Referencias**: REQ-014, DSG-006
**Estado**: ✅ Done

## TASK-019 — Eventos declarativos + data binding
**Prioridad**: Alta
**Dependencias**: TASK-014, TASK-017
**Archivos afectados**: cmd/hatch/handler.go, cmd/hatch/mouse_handler.go,
cmd/hatch/input.go, internal/parser/hml_parse.go
**Riesgo**: Medio — eventos afectan navegación y estado global
**Validación**: onclick="page:proyectos" navega. bind="target" sincroniza texto.
**Criterio de done**: executeEvent con soporte page/modal/action. executeDataBinding
en 5 puntos de input.go. onclick en handleMousePress.
**Referencias**: REQ-015, REQ-016, DSG-007
**Estado**: ✅ Done

## TASK-020 — CSS variables / theming
**Prioridad**: Media
**Dependencias**: TASK-005
**Archivos afectados**: internal/parser/hss.go, internal/parser/props.go,
internal/parser/compute.go, cmd/hatch/pipeline.go, cmd/hatch/main.go
**Riesgo**: Bajo — cambio localizado en parser
**Validación**: :root { --bg: #333; } + box { bg: var(--bg); } produce bg=#333.
**Criterio de done**: ParseCSSVars extrae variables :root. resolveCSSVars reemplaza
var(--name). vars pasadas a ComputeStyles. Variables globales en doc.ThemeVars.
**Referencias**: REQ-017, DSG-008
**Estado**: ✅ Done

## TASK-021 — Tab focus navigation
**Prioridad**: Media
**Dependencias**: TASK-014
**Archivos afectados**: cmd/hatch/focus.go, cmd/hatch/keyboard.go,
cmd/hatch/render_elements.go, cmd/hatch/render_helpers.go
**Riesgo**: Bajo — cambio localizado en keyboard dispatch
**Validación**: Tab navega entre input/textarea/button. Shift+Tab reversa.
**Criterio de done**: buildFocusOrder (DFS), FocusIndex, Tab/Shift+Tab en keyboard.go,
cursor visible en input/textarea al focus.
**Referencias**: REQ-018, DSG-009
**Estado**: ✅ Done

## TASK-022 — Demo rediseñada (7 componentes, 4 páginas)
**Prioridad**: Media
**Dependencias**: TASK-017, TASK-018, TASK-019, TASK-020, TASK-021
**Archivos afectados**: canva/demo.hml, canva/pages/*.hml, canva/components/*.hml
**Riesgo**: Bajo — solo archivos de assets
**Validación**: hatch corre y muestra la demo con navegación entre páginas
**Criterio de done**: 7 componentes nuevos (scroll, events, bind, theme, textarea,
layout, tabs). 4 páginas (inicio dashboard, proyectos scroll, config bind+events+theme,
ayuda textarea+tabs). demo.hml con onclick en navbar/sidebar. Componentes envueltos
en <page>.
**Referencias**: REQ-013, REQ-014, REQ-015, REQ-016, REQ-017, REQ-018
**Estado**: ✅ Done

---

## Flags de paralelización

| [PARALELIZABLE] | Tasks |
|----------------|-------|
| Sí | TASK-001, TASK-010, TASK-011, TASK-012 |
| Sí (mismo dominio) | TASK-002, TASK-010 |
| Sí (distinto dominio) | TASK-003 + TASK-004 (parser), TASK-007 (render) |
| No (misma dependencia) | TASK-003 → TASK-005 → TASK-006 (cadena) |
| No (integración final) | TASK-009 (depende de todas las anteriores) |
| No (refactor secuencial) | TASK-013 → TASK-014 → TASK-015 → TASK-016 (cadena) |
| Sí (paralelo a refactor) | TASK-017 (include) + TASK-020 (CSS vars) |
| No (dependencia de include) | TASK-017 → TASK-018 (scroll necesita AST con includes) |
| No (integración eventos) | TASK-019 (depende de include y pipeline) |
| No (demo final) | TASK-022 (depende de todas las features previas) |
