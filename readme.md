# Netwalk —— Go 复刻版

> *一款将所有终端接入主服务器的连线解谜游戏*
> 原作作者：**Bryan Lynn**（C 语言版） https://crypto.stanford.edu/~blynn/play/netwalk

> 本项目：使用 **Go + Ebiten** 重写

---

## 🎮 游戏简介

Netwalk 是一款经典 2D 益智游戏。棋盘由若干方块（Tile）组成，每个方块包含 0\~4 根电缆端口，不同方向预示可连通的方向。玩家的目标是在 **最少步数** 内旋转方块，将所有终端（Client）与中央服务器（Server） **连成一张无环网络**，并确保没有悬空电缆。

---

## 📜 游戏规则

| 元素               | 说明                          |
| ---------------- | --------------------------- |
| **服务器 (Server)** | 棋盘上唯一的发光核心，必须与所有终端保持连通      |
| **终端 (Client)**  | 其他发光方块，全部需要接入服务器            |
| **电缆 (Wire)**    | 连接线，允许横/竖/斜三类方向（NESW 对角线无效） |
| **空端口**          | 如果棋盘终局仍出现“裸线”，即失败           |

* **鼠标左键**：顺时针旋转 90°

---

## 🏗️ 项目结构

```
netwalk-go
├── cmd/
│   └── netwalk/     # 入口：解析 flags、启动 Ebiten
├── internal/
│   ├── game/        # 纯逻辑：Tile、Board、DFS 连通检测、关卡生成
│   ├── ui/ebiten/   # 渲染与输入：实现 ebiten.Game
├── go.mod
└── README.md
```

---

## 🚀 快速开始

> 依赖：Go ≥ 1.22

```bash
# 克隆仓库
$ git clone https://github.com/H1W0XXX/netwalk.git

$ cd netwalk
# 编译并运行（桌面版）
$ go run ./cmd/netwalk/main.go

```
