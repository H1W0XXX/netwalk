package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"netwalk/internal/game"
)

type GameView struct {
	board                     *game.Board
	solvedOnce                bool
	screenWidth, screenHeight int // 随窗口实时更新
}

// ---------------------- 构造 -------------------------

func NewGameView(b *game.Board) *GameView {
	// 初始尺寸 = main.go 里 SetWindowSize 设置的值
	return &GameView{
		board:        b,
		screenWidth:  b.W * 32,
		screenHeight: b.H * 32,
	}
}

// ---------------------- Ebiten 接口 ------------------

// Layout 允许窗口缩放，并随时记录当前像素尺寸
func (g *GameView) Layout(outsideW, outsideH int) (int, int) {
	// 最小允许像素（避免除零）
	if outsideW == 0 || outsideH == 0 {
		outsideW, outsideH = g.screenWidth, g.screenHeight
	}

	// 计算符合当前窗口、且保持宽高比的 tilePx
	tilePx := min(outsideW/g.board.W, outsideH/g.board.H)
	if tilePx < 1 {
		tilePx = 1
	}
	wantW := g.board.W * tilePx
	wantH := g.board.H * tilePx

	// 如果用户拖出了比例，立即把窗口吸附回等比尺寸
	if wantW != outsideW || wantH != outsideH {
		// ⚠️ 可以安全地在 Layout 里调用 SetWindowSize
		ebiten.SetWindowSize(wantW, wantH)
	}

	g.screenWidth, g.screenHeight = wantW, wantH
	return wantW, wantH
}

var frame int

// Update：每次点击仅旋转一次，依据当前 tile 大小计算格坐标
func (g *GameView) Update() error {
	frame++
	if frame == 1 { // 窗口已出现，改 FPSMode 不会卡住
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.solvedOnce {
		x, y := ebiten.CursorPosition()
		tilePx := min(g.screenWidth/g.board.W, g.screenHeight/g.board.H)
		tx, ty := x/tilePx, y/tilePx
		if game.In(tx, ty, g.board.W, g.board.H) {
			g.board.Rotate(tx, ty)
			g.solvedOnce = g.board.Solved()
		}
	}
	return nil
}

// Draw：按实时 tile 大小把整盘缩放绘制
func (g *GameView) Draw(screen *ebiten.Image) {
	tilePx := min(g.screenWidth/g.board.W, g.screenHeight/g.board.H)

	// 1) 背景网格
	gridCol := color.RGBA{192, 192, 192, 255}
	for x := 1; x < g.board.W; x++ {
		ebitenutil.DrawLine(screen,
			float64(x*tilePx), 0,
			float64(x*tilePx), float64(g.board.H*tilePx),
			gridCol)
	}
	for y := 1; y < g.board.H; y++ {
		ebitenutil.DrawLine(screen,
			0, float64(y*tilePx),
			float64(g.board.W*tilePx), float64(y*tilePx),
			gridCol)
	}

	// 2) 计算哪些格子被电源连通
	reach := g.board.Reachable()

	// 3) 绘制每个方块
	for y := 0; y < g.board.H; y++ {
		for x := 0; x < g.board.W; x++ {
			drawTile(screen, g.board, reach, x, y, tilePx)
		}
	}
}

// ------------ 渲染单块 ------------
func drawTile(dst *ebiten.Image, b *game.Board, reachable []bool, x, y, tile int) {
	// 1) 画线路
	cx := float64(x*tile + tile/2)
	cy := float64(y*tile + tile/2)
	conn := b.Grid[y*b.W+x].Conn
	half := float64(tile / 2)
	lineColor := color.RGBA{255, 127, 127, 255}
	if len(reachable) == len(b.Grid) { // 全盘通电后改为绿线
		lineColor = color.RGBA{0, 191, 0, 255}
	}
	if conn&game.Up != 0 {
		ebitenutil.DrawLine(dst, cx, cy, cx, cy-half, lineColor)
	}
	if conn&game.Right != 0 {
		ebitenutil.DrawLine(dst, cx, cy, cx+half, cy, lineColor)
	}
	if conn&game.Down != 0 {
		ebitenutil.DrawLine(dst, cx, cy, cx, cy+half, lineColor)
	}
	if conn&game.Left != 0 {
		ebitenutil.DrawLine(dst, cx, cy, cx-half, cy, lineColor)
	}

	// 2) 画电源和终端
	ox := float64(x * tile)
	oy := float64(y * tile)
	idx := y*b.W + x

	if b.IsRoot(x, y) {
		// 电源：蓝色
		ebitenutil.DrawRect(dst,
			ox+float64(tile)*0.25, oy+float64(tile)*0.25,
			float64(tile)*0.5, float64(tile)*0.5,
			color.RGBA{0, 0, 255, 255},
		)
	} else if b.IsEndpoint(x, y) {
		// 终端：未连通灰、连通黄
		col := color.RGBA{191, 191, 191, 255}
		if reachable[idx] {
			col = color.RGBA{255, 255, 0, 255}
		}
		ebitenutil.DrawRect(dst,
			ox+float64(tile)*0.29, oy+float64(tile)*0.29,
			float64(tile)*0.42, float64(tile)*0.42,
			col,
		)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
