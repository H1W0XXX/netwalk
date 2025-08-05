// internal/ui/loading.go
package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// LoadingView 只画“Generating board…”
type LoadingView struct {
	w, h, tile int
}

func NewLoadingView(w, h, tile int) *LoadingView {
	return &LoadingView{w: w, h: h, tile: tile}
}

func (l *LoadingView) Update() error { return nil }

func (l *LoadingView) Draw(dst *ebiten.Image) {
	ebitenutil.DebugPrint(dst, "Generating board…")
}

func (l *LoadingView) Layout(outW, outH int) (int, int) {
	if outW == 0 || outH == 0 {
		outW, outH = l.w*l.tile, l.h*l.tile
	}
	return outW, outH
}
