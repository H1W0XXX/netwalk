package game

// ----------------------------------
// 底层类型
// ----------------------------------

type Dir uint8

const (
	Up Dir = 1 << iota
	Right
	Down
	Left
)

var allDirs = []Dir{Up, Right, Down, Left}

type Tile struct {
	Conn Dir // 四位 bitmask，指示四方向是否连线
}

type Board struct {
	W, H     int
	Grid     []Tile // 行主序，一维数组
	rtX, rtY int
}

// ----------------------------------
// 构造与生成
// ----------------------------------
//func NewBoard(w, h int) *Board {
//	rand.Seed(time.Now().UnixNano())
//	b := &Board{W: w, H: h, Grid: make([]Tile, w*h)}
//
//	// 中央双格电源
//	b.rtX, b.rtY = w/2, h/2
//	b.setConn(b.rtX, b.rtY, Down) // rootTop ↔ rootBot
//	b.setConn(b.rtX, b.rtY+1, Up)
//
//	// 深度优先生成树
//	type cell struct{ x, y int }
//	stack := []cell{{b.rtX, b.rtY}}
//	visited := make([]bool, w*h)
//	visited[idx(w, b.rtX, b.rtY)] = true
//	for len(stack) > 0 {
//		c := &stack[len(stack)-1]
//		// 随机打乱四方向
//		dirs := []Dir{Up, Right, Down, Left}
//		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
//
//		extended := false
//		for _, d := range dirs {
//			nx, ny := step(c.x, c.y, d)
//			if !in(nx, ny, w, h) {
//				continue
//			}
//			ni := idx(w, nx, ny)
//			if visited[ni] {
//				continue
//			}
//			// 连边
//			b.setConn(c.x, c.y, d)
//			b.setConn(nx, ny, opposite(d))
//
//			visited[ni] = true
//			stack = append(stack, cell{nx, ny})
//			extended = true
//			break
//		}
//		if !extended {
//			// 回溯
//			stack = stack[:len(stack)-1]
//		}
//	}
//	return b
//}

// ----------------------------------
// API
// ----------------------------------
//
//func (b *Board) Rotate(x, y int) {
//	i := idx(b.W, x, y)
//	c := b.Grid[i].Conn
//	b.Grid[i].Conn = ((c << 1) | (c >> 3)) & 0xF
//}
//
//func (b *Board) Shuffle() {
//	for i := range b.Grid {
//		r := rand.Intn(4)
//		c := b.Grid[i].Conn
//		b.Grid[i].Conn = ((c << r) | (c >> (4 - r))) & 0xF
//	}
//}
//
//func (b *Board) Solved(_ int) bool {
//	seen := make([]bool, len(b.Grid))
//	var dfs func(x, y int)
//	dfs = func(x, y int) {
//		i := idx(b.W, x, y)
//		if seen[i] {
//			return
//		}
//		seen[i] = true
//		for _, d := range allDirs {
//			if b.hasConn(x, y, d) {
//				nx, ny := step(x, y, d)
//				if in(nx, ny, b.W, b.H) && b.hasConn(nx, ny, opposite(d)) {
//					dfs(nx, ny)
//				}
//			}
//		}
//	}
//	dfs(b.rtX, b.rtY)
//	dfs(b.rtX, b.rtY+1)
//	for i, t := range b.Grid {
//		if t.Conn != 0 && !seen[i] {
//			return false
//		}
//	}
//	return true
//}
//
//func (b *Board) IsRoot(x, y int) bool {
//	return x == b.rtX && (y == b.rtY || y == b.rtY+1)
//}
//
//func (b *Board) IsEndpoint(x, y int) bool {
//	if b.IsRoot(x, y) {
//		return false
//	}
//	c := b.Grid[idx(b.W, x, y)].Conn
//	return bits.OnesCount8(uint8(c)) == 1
//}
//
//// Reachable returns a slice marking which tiles are connected to the power source.
//func (b *Board) Reachable() []bool {
//	seen := make([]bool, len(b.Grid))
//	var dfs func(x, y int)
//	dfs = func(x, y int) {
//		i := idx(b.W, x, y)
//		if seen[i] {
//			return
//		}
//		seen[i] = true
//		for _, d := range allDirs {
//			if b.hasConn(x, y, d) {
//				nx, ny := step(x, y, d)
//				if in(nx, ny, b.W, b.H) && b.hasConn(nx, ny, opposite(d)) {
//					dfs(nx, ny)
//				}
//			}
//		}
//	}
//	dfs(b.rtX, b.rtY)
//	dfs(b.rtX, b.rtY+1)
//	return seen
//}

// ----------------------------------
// 辅助
// ----------------------------------
//
//func idx(w, x, y int) int    { return y*w + x }
//func in(x, y, w, h int) bool { return x >= 0 && x < w && y >= 0 && y < h }
//func In(x, y, w, h int) bool { return in(x, y, w, h) }
//func step(x, y int, d Dir) (int, int) {
//	switch d {
//	case Up:
//		return x, y - 1
//	case Down:
//		return x, y + 1
//	case Left:
//		return x - 1, y
//	default: // Right
//		return x + 1, y
//	}
//}
//func opposite(d Dir) Dir { return Dir(((d << 2) | (d >> 2)) & 0xF) }
//
//func (b *Board) setConn(x, y int, d Dir) {
//	b.Grid[idx(b.W, x, y)].Conn |= d
//}
//func (b *Board) hasConn(x, y int, d Dir) bool {
//	return b.Grid[idx(b.W, x, y)].Conn&d != 0
//}
