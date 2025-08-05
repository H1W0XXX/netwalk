// cmd/warmup/main.go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type warm struct {
	frame int
}

func (w *warm) Update() error {
	w.frame++
	if w.frame > 10 { // 跑 10 帧就够
		return ebiten.Termination // v2.5+ 结束 RunGame（无错误）
	}
	return nil
}

func (w *warm) Draw(dst *ebiten.Image) {
	// 触发 3 套常用管线的编译
	ebitenutil.DrawLine(dst, 1, 1, 15, 15, color.White)          // 线
	ebitenutil.DrawRect(dst, 2, 4, 8, 8, color.RGBA{1, 1, 1, 1}) // 矩形
	ebitenutil.DebugPrint(dst, "warm")                           // 字体
}

func (w *warm) Layout(int, int) (int, int) { return 16, 16 }

func main() {
	ebiten.SetWindowSize(16, 16) // 极小不可见窗口
	_ = ebiten.RunGame(&warm{})
}

// go build -ldflags="-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o warmup.exe .\cmd\warmup\main.go
