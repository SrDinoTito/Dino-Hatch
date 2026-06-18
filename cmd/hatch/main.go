// Hatch CLI — entry point: flags → parse → tcell → RunLoop.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/parser"
)

func main() {
	logFile, err := os.OpenFile("hatch.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: hatch run <archivo.hml>\n")
	}
	flag.Parse()
	if flag.NArg() < 2 || flag.Arg(0) != "run" {
		flag.Usage()
		os.Exit(1)
	}

	hmlPath := flag.Arg(1)
	data, err := os.ReadFile(hmlPath)
	if err != nil {
		log.Fatalf("Error leyendo %s: %v", hmlPath, err)
	}

	// Parsear HML
	doc, styleBlocks, err := parser.ParseHML(data)
	if err != nil {
		log.Fatalf("Error parseando HML: %v", err)
	}

	// Parsear HSS
	var rules []ast.StyleRule
	var allStyleContent string
	for _, block := range styleBlocks {
		allStyleContent += block + "\n"
		r, err := parser.ParseHSS(block)
		if err != nil {
			log.Printf("Warning: error en bloque style: %v", err)
			continue
		}
		rules = append(rules, r...)
	}
	// Wave 6: Extraer CSS variables de los bloques :root
	if allStyleContent != "" {
		if vars := parser.ParseCSSVars(allStyleContent); len(vars) > 0 {
			doc.ThemeVars = vars
		}
	}
	doc = parser.ComputeStyles(doc, rules, doc.ThemeVars)

	// Inicializar tcell
	tcellScr, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creando screen: %v", err)
	}
	if err := tcellScr.Init(); err != nil {
		log.Fatalf("Error inicializando screen: %v", err)
	}
	tcellScr.EnableMouse(tcell.MouseMotionEvents)
	defer tcellScr.Fini()
	baseDir := filepath.Dir(hmlPath)
	state = NewAppState(doc, rules, tcellScr, baseDir)
	defer state.Eng.Close()
	defer state.EventBus.Close()

	// Si stdin es pipe, publicar lineas como eventos al bus
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				state.EventBus.Publish("stdin", scanner.Text())
			}
		}()
	}
	RunLoop(state)
}
