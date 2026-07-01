package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	shakeX, shakeY := 0.0, 0.0
	if g.ScreenShake > 0 {
		shakeX = (rand.Float64() - 0.5) * float64(g.ScreenShake)
		shakeY = (rand.Float64() - 0.5) * float64(g.ScreenShake)
	}
	if g.offscreen == nil {
		g.offscreen = ebiten.NewImage(ScreenWidth, ScreenHeight)
	}
	g.offscreen.Fill(color.RGBA{8, 5, 18, 255})

	g.drawStars()

	switch g.State {
	case StateMenu:
		g.drawMenu()
	case StatePlaying, StatePaused:
		g.drawGame()
		if g.State == StatePaused {
			g.drawPaused()
		}
	case StateLevelUp:
		g.drawGame()
		g.drawLevelUp()
	case StateGameOver:
		g.drawGame()
		g.drawGameOver()
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(g.offscreen, op)
}

func (g *Game) drawStars() {
	for _, s := range g.Stars {
		c := color.RGBA{200, 220, 255, uint8(255 * s.Alpha)}
		vector.DrawFilledRect(g.offscreen, float32(s.X), float32(s.Y), float32(s.Size), float32(s.Size), c, false)
	}
}

func (g *Game) drawMenu() {
	f := fontNormal()
	title := "霓虹幸存者"
	b := text.BoundString(f, title)
	text.Draw(g.offscreen, title, f, ScreenWidth/2-b.Dx()/2, 180, color.RGBA{100, 220, 255, 255})

	subtitle := "NEON SURVIVOR"
	b = text.BoundString(f, subtitle)
	text.Draw(g.offscreen, subtitle, f, ScreenWidth/2-b.Dx()/2, 210, color.RGBA{200, 150, 255, 255})

	g.drawPlayerShip(ScreenWidth/2, 310, 1.8, 1)

	prompt := "点击鼠标 或 按空格 开始游戏"
	b = text.BoundString(f, prompt)
	flash := float32(0.5 + 0.5*math.Sin(g.Time*3))
	c := color.RGBA{uint8(100 + 155*flash), uint8(255 * flash), uint8(200 * flash), 255}
	text.Draw(g.offscreen, prompt, f, ScreenWidth/2-b.Dx()/2, 440, c)

	coinText := fmt.Sprintf("局外金币: %d", g.Meta.Coins)
	if g.Meta.InfiniteGold {
		coinText = "局外金币: 无限"
	}
	b = text.BoundString(f, coinText)
	text.Draw(g.offscreen, coinText, f, ScreenWidth/2-b.Dx()/2, 485, color.RGBA{255, 220, 100, 255})

	hs := fmt.Sprintf("最高分: %d", g.HighScore)
	b = text.BoundString(f, hs)
	text.Draw(g.offscreen, hs, f, ScreenWidth/2-b.Dx()/2, 510, color.RGBA{255, 200, 100, 255})

	controls := []string{
		"操作说明:",
		"鼠标  - 控制战机移动",
		"自动  - 武器自动射击",
		"经验  - 击杀敌人获得，升级变强",
		"升级  - 三选一，打造你的最强战机",
		"空格/右键 - 暂停游戏",
		fmt.Sprintf("1 初始血量+25  价格%d  当前+%d", metaHPCost, g.Meta.HPLevel*25),
		fmt.Sprintf("2 初始子弹+1  价格%d  当前%d/3", metaBulletCost, 1+g.Meta.BulletLevel),
		fmt.Sprintf("3 加速卡  价格%d  %s", metaBoostCost, ownedText(g.Meta.BoostCard)),
		fmt.Sprintf("4 初始抽卡+1  价格%d  当前%d/3", metaCardBaseCost*(g.Meta.InitialCards+1), g.Meta.InitialCards),
		"彩蛋: 上上下下左右左右BA 解锁无限金币",
	}
	for i, line := range controls {
		text.Draw(g.offscreen, line, f, 30, ScreenHeight-230+i*20, color.RGBA{150, 150, 200, 255})
	}
}

func ownedText(owned bool) string {
	if owned {
		return "已购买"
	}
	return "未购买"
}

func (g *Game) drawGame() {
	for _, b := range g.Bullets {
		if b.Beam {
			x1 := float32(b.X)
			y1 := float32(b.Y)
			x2 := float32(b.X + b.VX*28)
			y2 := float32(b.Y + b.VY*28)
			pulse := float32(0.75 + 0.25*math.Sin(g.Time*30+float64(b.X)))
			glow := color.RGBA{b.Color.R, b.Color.G, b.Color.B, 45}
			mid := color.RGBA{b.Color.R, b.Color.G, b.Color.B, uint8(130 * pulse)}
			core := color.RGBA{255, 255, 255, uint8(220 * pulse)}
			vector.StrokeLine(g.offscreen, x1, y1, x2, y2, float32(b.Size)*2.8, glow, false)
			vector.StrokeLine(g.offscreen, x1, y1, x2, y2, float32(b.Size)*1.5, mid, false)
			vector.StrokeLine(g.offscreen, x1, y1, x2, y2, max(float32(2), float32(b.Size)*0.35), core, false)
			vector.DrawFilledCircle(g.offscreen, x2, y2, float32(b.Size)*0.8*pulse, core, false)
			continue
		}
		if b.Player {
			if b.Size > 6 {
				vector.DrawFilledCircle(g.offscreen, float32(b.X), float32(b.Y), float32(b.Size+3), color.RGBA{b.Color.R, b.Color.G, b.Color.B, 80}, false)
			}
			vector.DrawFilledCircle(g.offscreen, float32(b.X), float32(b.Y), float32(b.Size), b.Color, false)
			vector.DrawFilledCircle(g.offscreen, float32(b.X), float32(b.Y), float32(b.Size)*0.5, color.RGBA{255, 255, 255, 200}, false)
		} else {
			vector.DrawFilledCircle(g.offscreen, float32(b.X), float32(b.Y), float32(b.Size), b.Color, false)
			vector.DrawFilledCircle(g.offscreen, float32(b.X), float32(b.Y), float32(b.Size)*0.5, color.RGBA{255, 255, 255, 180}, false)
		}
	}

	for _, e := range g.Enemies {
		g.drawEnemy(e)
	}

	for _, pu := range g.PowerUps {
		g.drawPowerUp(pu)
	}

	for _, c := range g.Coins {
		vector.DrawFilledCircle(g.offscreen, float32(c.X), float32(c.Y), float32(c.Size)*0.7, color.RGBA{255, 210, 80, 255}, false)
	}

	for _, e := range g.ExpOrbs {
		vector.DrawFilledCircle(g.offscreen, float32(e.X), float32(e.Y), float32(e.Size)*0.7, color.RGBA{110, 255, 150, 255}, false)
	}

	for _, s := range g.SlashWaves {
		g.drawSlashWave(s)
	}

	for _, p := range g.Particles {
		alpha := float32(p.Life) / float32(p.MaxLife)
		c := color.RGBA{p.Color.R, p.Color.G, p.Color.B, uint8(float32(p.Color.A) * alpha)}
		vector.DrawFilledCircle(g.offscreen, float32(p.X), float32(p.Y), float32(p.Size), c, false)
	}

	if g.State == StatePlaying || g.State == StatePaused || g.State == StateLevelUp {
		flash := 1.0
		if g.Player.Invincible > 0 && g.Player.Invincible%6 < 3 {
			flash = 0.4
		}
		g.drawPlayerShip(g.Player.X, g.Player.Y, 1.0, float32(flash))

		if g.Player.Shield > 0 {
			a := float32(0.3 + 0.2*math.Sin(g.Time*5))
			vector.StrokeCircle(g.offscreen, float32(g.Player.X), float32(g.Player.Y), 32, 2.5, color.RGBA{100, 200, 255, uint8(255 * a)}, false)
		}


	}

	g.drawHUD()
}

func (g *Game) fillTriangle(x0, y0, x1, y1, x2, y2 float32, clr color.Color) {
	var path vector.Path
	path.MoveTo(x0, y0)
	path.LineTo(x1, y1)
	path.LineTo(x2, y2)
	path.Close()
	vector.FillPath(g.offscreen, &path, nil, &vector.DrawPathOptions{ColorScale: colorScale(clr)})
}

func (g *Game) fillPoly(cx, cy float32, pts [][2]float32, clr color.Color) {
	var path vector.Path
	for i, pt := range pts {
		if i == 0 {
			path.MoveTo(cx+pt[0], cy+pt[1])
		} else {
			path.LineTo(cx+pt[0], cy+pt[1])
		}
	}
	path.Close()
	vector.FillPath(g.offscreen, &path, nil, &vector.DrawPathOptions{ColorScale: colorScale(clr)})
}

func (g *Game) fillOval(cx, cy, rx, ry float32, clr color.Color) {
	var path vector.Path
	segments := 20
	for i := 0; i <= segments; i++ {
		angle := float32(i) * 2 * math.Pi / float32(segments)
		x := cx + rx*float32(math.Cos(float64(angle)))
		y := cy + ry*float32(math.Sin(float64(angle)))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()
	vector.FillPath(g.offscreen, &path, nil, &vector.DrawPathOptions{ColorScale: colorScale(clr)})
}
