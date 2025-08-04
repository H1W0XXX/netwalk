package game

import (
	"math/bits"
	"math/rand"
	"time"
)

// ----------------------------------
// 构造与生成
// ----------------------------------
func NewBoard(w, h int) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{W: w, H: h, Grid: make([]Tile, w*h)}

	// 中央双格电源
	b.rtX, b.rtY = w/2, h/2
	b.setConn(b.rtX, b.rtY, Down) // rootTop ↔ rootBot
	b.setConn(b.rtX, b.rtY+1, Up)

	// 改进的生成算法：先生成基础树，再添加分支
	b.generateEnhancedNetwork()

	return b
}

// generateEnhancedNetwork 生成更多终端的网络
func (b *Board) generateEnhancedNetwork() {
	// 第一阶段：生成基础连通树
	visited := make([]bool, b.W*b.H)
	visited[idx(b.W, b.rtX, b.rtY)] = true
	visited[idx(b.W, b.rtX, b.rtY+1)] = true

	// 使用广度优先搜索创建更平衡的树结构
	queue := []struct{ x, y int }{{b.rtX, b.rtY}, {b.rtX, b.rtY + 1}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// 随机选择要连接的方向数量 (1-3个方向)
		dirs := []Dir{Up, Right, Down, Left}
		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

		// 限制每个节点的连接数，增加终端概率
		maxConnections := 1 + rand.Intn(3) // 1-3个连接
		connections := 0

		for _, d := range dirs {
			if connections >= maxConnections {
				break
			}

			nx, ny := step(curr.x, curr.y, d)
			if !in(nx, ny, b.W, b.H) {
				continue
			}

			ni := idx(b.W, nx, ny)
			if visited[ni] {
				continue
			}

			// 随机决定是否连接（增加分支的随机性）
			if rand.Float32() < 0.7 { // 70%概率连接
				b.setConn(curr.x, curr.y, d)
				b.setConn(nx, ny, opposite(d))
				visited[ni] = true
				queue = append(queue, struct{ x, y int }{nx, ny})
				connections++
			}
		}
	}

	// 第二阶段：添加额外的分支来增加终端数量
	b.addExtraBranches(visited)

	// 第三阶段：随机添加一些短分支
	b.addShortBranches(visited)
}

// addExtraBranches 添加额外分支来增加终端数量
func (b *Board) addExtraBranches(visited []bool) {
	// 找到所有已连接的节点
	connectedNodes := []struct{ x, y int }{}
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			if visited[idx(b.W, x, y)] && b.Grid[idx(b.W, x, y)].Conn != 0 {
				connectedNodes = append(connectedNodes, struct{ x, y int }{x, y})
			}
		}
	}

	// 从已连接的节点随机选择一些来添加分支
	rand.Shuffle(len(connectedNodes), func(i, j int) {
		connectedNodes[i], connectedNodes[j] = connectedNodes[j], connectedNodes[i]
	})

	branchesToAdd := len(connectedNodes) / 3 // 大约1/3的节点添加分支
	for i := 0; i < branchesToAdd && i < len(connectedNodes); i++ {
		node := connectedNodes[i]
		b.tryAddBranch(node.x, node.y, visited)
	}
}

// tryAddBranch 尝试从指定节点添加一个分支
func (b *Board) tryAddBranch(x, y int, visited []bool) {
	dirs := []Dir{Up, Right, Down, Left}
	rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

	for _, d := range dirs {
		if b.hasConn(x, y, d) {
			continue // 这个方向已经有连接了
		}

		nx, ny := step(x, y, d)
		if !in(nx, ny, b.W, b.H) {
			continue
		}

		ni := idx(b.W, nx, ny)
		if visited[ni] {
			continue
		}

		// 添加分支
		b.setConn(x, y, d)
		b.setConn(nx, ny, opposite(d))
		visited[ni] = true

		// 有概率继续延伸这个分支
		if rand.Float32() < 0.4 { // 40%概率延伸
			b.extendBranch(nx, ny, visited, 1+rand.Intn(3)) // 延伸1-3步
		}
		break
	}
}

// extendBranch 延伸分支
func (b *Board) extendBranch(x, y int, visited []bool, maxSteps int) {
	for i := 0; i < maxSteps; i++ {
		dirs := []Dir{Up, Right, Down, Left}
		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

		extended := false
		for _, d := range dirs {
			if b.hasConn(x, y, d) {
				continue
			}

			nx, ny := step(x, y, d)
			if !in(nx, ny, b.W, b.H) {
				continue
			}

			ni := idx(b.W, nx, ny)
			if visited[ni] {
				continue
			}

			b.setConn(x, y, d)
			b.setConn(nx, ny, opposite(d))
			visited[ni] = true
			x, y = nx, ny
			extended = true
			break
		}

		if !extended {
			break
		}
	}
}

// addShortBranches 添加一些短分支来确保足够的终端数量
func (b *Board) addShortBranches(visited []bool) {
	attempts := b.W * b.H / 4 // 尝试次数

	for attempt := 0; attempt < attempts; attempt++ {
		// 随机选择一个已访问的位置
		candidates := []struct{ x, y int }{}
		for y := 0; y < b.H; y++ {
			for x := 0; x < b.W; x++ {
				if visited[idx(b.W, x, y)] {
					candidates = append(candidates, struct{ x, y int }{x, y})
				}
			}
		}

		if len(candidates) == 0 {
			break
		}

		baseNode := candidates[rand.Intn(len(candidates))]

		// 检查这个节点是否还能添加连接
		availableDirs := []Dir{}
		for _, d := range allDirs {
			if !b.hasConn(baseNode.x, baseNode.y, d) {
				nx, ny := step(baseNode.x, baseNode.y, d)
				if in(nx, ny, b.W, b.H) && !visited[idx(b.W, nx, ny)] {
					availableDirs = append(availableDirs, d)
				}
			}
		}

		if len(availableDirs) > 0 {
			d := availableDirs[rand.Intn(len(availableDirs))]
			nx, ny := step(baseNode.x, baseNode.y, d)

			b.setConn(baseNode.x, baseNode.y, d)
			b.setConn(nx, ny, opposite(d))
			visited[idx(b.W, nx, ny)] = true
		}
	}
}

// ----------------------------------
// API
// ----------------------------------

func (b *Board) Rotate(x, y int) {
	i := idx(b.W, x, y)
	c := b.Grid[i].Conn
	b.Grid[i].Conn = ((c << 1) | (c >> 3)) & 0xF
}

func (b *Board) Shuffle() {
	for i := range b.Grid {
		r := rand.Intn(4)
		c := b.Grid[i].Conn
		b.Grid[i].Conn = ((c << r) | (c >> (4 - r))) & 0xF
	}
}

func (b *Board) Solved() bool {
	seen := make([]bool, len(b.Grid))
	var dfs func(x, y int)
	dfs = func(x, y int) {
		i := idx(b.W, x, y)
		if seen[i] {
			return
		}
		seen[i] = true
		for _, d := range allDirs {
			if b.hasConn(x, y, d) {
				nx, ny := step(x, y, d)
				if in(nx, ny, b.W, b.H) && b.hasConn(nx, ny, opposite(d)) {
					dfs(nx, ny)
				}
			}
		}
	}
	dfs(b.rtX, b.rtY)
	dfs(b.rtX, b.rtY+1)
	for i, t := range b.Grid {
		if t.Conn != 0 && !seen[i] {
			return false
		}
	}
	return true
}

func (b *Board) IsRoot(x, y int) bool {
	return x == b.rtX && (y == b.rtY || y == b.rtY+1)
}

func (b *Board) IsEndpoint(x, y int) bool {
	if b.IsRoot(x, y) {
		return false
	}
	c := b.Grid[idx(b.W, x, y)].Conn
	return bits.OnesCount8(uint8(c)) == 1
}

// Reachable returns a slice marking which tiles are connected to the power source.
func (b *Board) Reachable() []bool {
	seen := make([]bool, len(b.Grid))
	var dfs func(x, y int)
	dfs = func(x, y int) {
		i := idx(b.W, x, y)
		if seen[i] {
			return
		}
		seen[i] = true
		for _, d := range allDirs {
			if b.hasConn(x, y, d) {
				nx, ny := step(x, y, d)
				if in(nx, ny, b.W, b.H) && b.hasConn(nx, ny, opposite(d)) {
					dfs(nx, ny)
				}
			}
		}
	}
	dfs(b.rtX, b.rtY)
	dfs(b.rtX, b.rtY+1)
	return seen
}

// CountEndpoints 返回棋盘上"终端"(degree==1) 的数量。
func (b *Board) CountEndpoints() int {
	cnt := 0
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			if b.IsEndpoint(x, y) {
				cnt++
			}
		}
	}
	return cnt
}

// NewBoardWithMinEndpoints 生成直到终端数 ≥ want。
// maxTry 限制循环次数，避免极端情况下长时间卡住。
func NewBoardWithMinEndpoints(w, h, want, maxTry int) *Board {
	var best *Board
	maxLeaf := -1

	for i := 0; i < maxTry; i++ {
		b := NewBoard(w, h)
		leaf := b.CountEndpoints()

		if leaf >= want {
			b.Shuffle() // 最后再随机旋转一次
			return b
		}

		if leaf > maxLeaf {
			best, maxLeaf = b, leaf
		}
	}

	if best != nil {
		best.Shuffle()
		return best
	}

	// 如果没有找到合适的，创建一个新的
	b := NewBoard(w, h)
	b.Shuffle()
	return b
}

// NewBoardWithExactEndpoints 若想精确控制终端数量可用此版本。
func NewBoardWithExactEndpoints(w, h, exact, maxTry int) *Board {
	var (
		best        *Board
		bestDiffAbs = 1<<31 - 1
		bestLeaf    = -1
	)

	for i := 0; i < maxTry; i++ {
		b := NewBoard(w, h)
		leaf := b.CountEndpoints()

		if leaf == exact {
			b.Shuffle()
			return b
		}

		diff := leaf - exact
		if diff < 0 {
			diff = -diff
		}

		if diff < bestDiffAbs || (diff == bestDiffAbs && leaf > bestLeaf) {
			best, bestDiffAbs, bestLeaf = b, diff, leaf
		}
	}

	if best != nil {
		best.Shuffle()
		return best
	}

	// 如果没有找到合适的，创建一个新的
	b := NewBoard(w, h)
	b.Shuffle()
	return b
}

// ----------------------------------
// 辅助函数
// ----------------------------------

func idx(w, x, y int) int    { return y*w + x }
func in(x, y, w, h int) bool { return x >= 0 && x < w && y >= 0 && y < h }
func In(x, y, w, h int) bool { return in(x, y, w, h) }

func step(x, y int, d Dir) (int, int) {
	switch d {
	case Up:
		return x, y - 1
	case Down:
		return x, y + 1
	case Left:
		return x - 1, y
	default: // Right
		return x + 1, y
	}
}

func opposite(d Dir) Dir {
	return Dir(((d << 2) | (d >> 2)) & 0xF)
}

func (b *Board) setConn(x, y int, d Dir) {
	b.Grid[idx(b.W, x, y)].Conn |= d
}

func (b *Board) hasConn(x, y int, d Dir) bool {
	return b.Grid[idx(b.W, x, y)].Conn&d != 0
}
