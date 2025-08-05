package ui

import "github.com/hajimehoshi/ebiten/v2"

// View 是 Update/Draw/Layout 的集合
type View interface {
	Update() error
	Draw(*ebiten.Image)
	Layout(int, int) (int, int)
}

// Wrapper 持有当前要渲染的 View，可在运行时替换
type Wrapper struct {
	View
}

func (w *Wrapper) Update() error             { return w.View.Update() }
func (w *Wrapper) Draw(screen *ebiten.Image) { w.View.Draw(screen) }
func (w *Wrapper) Layout(wid, hei int) (int, int) {
	return w.View.Layout(wid, hei)
}
