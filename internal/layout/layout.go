package layout

import "github.com/srdino/dino-hatch/internal/ast"

// layoutChildren posiciona los hijos de un contenedor segun
// el estilo computado (flexbox simplificado: direction, grow,
// gap, align, justify).
func layoutChildren(children *[]ast.ElementNode, cx, cy, cw, ch int, s *ast.ComputedStyle) {
	n := len(*children)
	if n == 0 {
		return
	}
	isRow := s.Direction == "row"
	gap := s.Gap
	// Fase 1: calcular tamano primario de cada hijo
	sizes := make([]int, n)
	fixedPrimary := 0
	var totalGrow float64

	for i := range *children {
		child := &(*children)[i]
		if child.Style.Grow <= 0 {
			sz := child.Style.Width
			if !isRow {
				sz = child.Style.Height
			}
			sizes[i] = sz
			fixedPrimary += sz
		} else {
			totalGrow += child.Style.Grow
		}
	}

	// Tamano intrinseco para elementos sin tamano explicito ni grow
	for i := range *children {
		child := &(*children)[i]
		if child.Style.Grow <= 0 && sizes[i] <= 0 {
			isz := intrinsicSize(child, isRow)
			prev := child.BoundBox.H
			if isRow {
				prev = child.BoundBox.W
			}
			// Estabilizacion selectiva:
			// - Nodos con MaxHeight (textarea): prevenir shrinkage pero permitir
			//   crecimiento hasta MaxHeight cuando el contenido se expande.
			// - Nodos sin MaxHeight (contenedores): estabilizacion bilateral
			//   (mantener tamano del frame anterior, ni crecer ni decrecer).
			if prev > 0 {
				if child.Style.MaxHeight > 0 {
					// Textarea: crece con el contenido hasta MaxHeight, nunca se achica
					if isz > child.Style.MaxHeight {
						isz = child.Style.MaxHeight // Cap al tope
					}
					if isz < prev {
						isz = prev // No se achica
					}
				} else {
					// Non-textarea: estabilizacion bilateral (tamano fijo entre frames)
					isz = prev
				}
			}
			sizes[i] = isz
			fixedPrimary += sizes[i]
		}
	}

	containerPrimary := cw
	if !isRow {
		containerPrimary = ch
	}

	totalGap := 0
	if n > 1 {
		totalGap = gap * (n - 1)
	}

	// Calcular margen total de los hijos (cada hijo contribuye margin*2 en eje primario)
	totalMargin := 0
	for i := range *children {
		totalMargin += (*children)[i].Style.Margin * 2
	}

	remaining := containerPrimary - fixedPrimary - totalGap - totalMargin

	// Distribuir espacio restante entre hijos con grow
	if remaining > 0 && totalGrow > 0 {
		for i := range *children {
			child := &(*children)[i]
			if child.Style.Grow > 0 {
				sz := int(float64(remaining) * child.Style.Grow / totalGrow)
				if sz < 0 {
					sz = 0
				}
				sizes[i] = sz
			}
		}
	} else if remaining <= 0 {
		// Overflow: hijos con grow no reciben espacio
		for i := range *children {
			if (*children)[i].Style.Grow > 0 {
				sizes[i] = 0
			}
		}
	}

	// Clamping min/max en eje primario (despues de grow/intrinsic)
	for i := range *children {
		child := &(*children)[i]
		if isRow {
			if child.Style.MinWidth > 0 && sizes[i] < child.Style.MinWidth {
				sizes[i] = child.Style.MinWidth
			}
			if child.Style.MaxWidth > 0 && sizes[i] > child.Style.MaxWidth {
				sizes[i] = child.Style.MaxWidth
			}
		} else {
			if child.Style.MinHeight > 0 && sizes[i] < child.Style.MinHeight {
				sizes[i] = child.Style.MinHeight
			}
			if child.Style.MaxHeight > 0 && sizes[i] > child.Style.MaxHeight {
				sizes[i] = child.Style.MaxHeight
			}
		}
	}

	// Calcular espacio total usado para justify-content (incluyendo margenes)
	totalUsed := 0
	for i, sz := range sizes {
		totalUsed += sz
		totalUsed += (*children)[i].Style.Margin * 2
	}
	if n > 1 {
		totalUsed += gap * (n - 1)
	}
	extraSpace := containerPrimary - totalUsed
	if extraSpace < 0 {
		extraSpace = 0
	}

	// Aplicar justify-content sobre eje principal
	var startOffset, betweenGap int
	switch s.Justify {
	case "center":
		startOffset = extraSpace / 2
		betweenGap = gap
	case "end":
		startOffset = extraSpace
		betweenGap = gap
	case "space-between":
		if n > 1 {
			betweenGap = extraSpace / (n - 1)
		}
		// startOffset = 0
	default: // "start"
		betweenGap = gap
	}

	// Tamano del eje transversal (stretch por defecto)
	crossSize := ch
	if !isRow {
		crossSize = cw
	}

	// Posicionar cada hijo secuencialmente
	primaryPos := cx + startOffset
	if !isRow {
		primaryPos = cy + startOffset
	}

	for i := range *children {
		child := &(*children)[i]
		ps := sizes[i]
		if ps < 0 {
			ps = 0
		}

		// Aplicar margen del hijo en eje primario (espacio externo entre hermanos)
		margin := child.Style.Margin
		primaryPos += margin

		// Calcular tamano y offset en eje transversal
		cs := crossSize
		co := 0
		if s.Align != "stretch" && s.Align != "" {
			ecs := child.Style.Height // cross axis for row
			if !isRow {
				ecs = child.Style.Width // cross axis for column
			}
			if ecs > 0 && ecs < crossSize {
				cs = ecs
			}
			switch s.Align {
			case "center":
				co = (crossSize - cs) / 2
			case "end":
				co = crossSize - cs
			}
			// "start": co = 0
		}

		// Asignar BoundBox
		if isRow {
			child.BoundBox = ast.BoundBox{X: primaryPos, Y: cy + co, W: ps, H: cs}
			primaryPos += ps + margin + betweenGap
		} else {
			child.BoundBox = ast.BoundBox{X: cx + co, Y: primaryPos, W: cs, H: ps}
			primaryPos += ps + margin + betweenGap
		}

		// Recurrir en los hijos de este nodo
		layoutNode(child, child.BoundBox.X, child.BoundBox.Y, child.BoundBox.W, child.BoundBox.H)
	}
}
