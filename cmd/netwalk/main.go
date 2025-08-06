package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"netwalk/internal/game"
	"netwalk/internal/ui"
)

func main() {
	const (
		w, h     = 10, 9
		tileInit = 32
		fps      = 10
	)
	ebiten.SetWindowSize(w*tileInit, h*tileInit)
	ebiten.SetWindowTitle("NetWalk – Go Edition")
	ebiten.SetWindowResizable(true)
	//ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum) // 让 Draw 频率跟随 TPS
	ebiten.SetMaxTPS(fps)

	//b := game.NewBoard(w, h)
	b := game.NewBoardWithMinEndpoints(w, h, 30, 3000)
	b.Shuffle()
	g := ui.NewGameView(b)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// go build -ldflags="-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o netwalk.exe .\cmd\netwalk\main.go
