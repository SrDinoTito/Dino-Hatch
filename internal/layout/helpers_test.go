package layout

import (
	"testing"

	"github.com/srdino/dino-hatch/internal/ast"
)

func assertBox(t *testing.T, got, want ast.BoundBox) {
	t.Helper()
	if got != want {
		t.Errorf("BoundBox = %+v, want %+v", got, want)
	}
}

func assertEq(t *testing.T, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %d, want %d", msg, got, want)
	}
}
