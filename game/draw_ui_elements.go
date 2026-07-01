package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawEnemyHpBar(e Enemy, xf, yf float32) {
	barW := float32(e.Width)
	hpRatio := float32(e.Health) / float32(e.MaxHealth)
	by := yf - float32(e.Height/2) - 8
	vector.DrawFilledRect(g.offscreen, xf-barW/2, by, barW, 3, color.RGBA{60, 60, 80, 200}, false)
	vector.DrawFilledRect(g.offscreen, xf-barW/2, by, barW*hpRatio, 3, color.RGBA{100, 255, 150, 255}, false)
}

func (g *Game) drawBossHpBar(e Enemy) {
	barW := float32(ScreenWidth - 80)
	x := float32(40)
	y := float32(45)
	hpRatio := float32(e.Health) / float32(e.MaxHealth)
	vector.DrawFilledRect(g.offscreen, x, y, barW, 12, color.RGBA{60, 30, 40, 200}, false)
	vector.DrawFilledRect(g.offscreen, x, y, barW*hpRatio, 12, color.RGBA{255, 60, 100, 255}, false)
	vector.StrokeRect(g.offscreen, x, y, barW, 12, 2, color.RGBA{255, 150, 180, 255}, false)

	f := fontSmall()
	name := "首领"
	b := text.BoundString(f, name)
	text.Draw(g.offscreen, name, f, ScreenWidth/2-b.Dx()/2, 43, color.RGBA{255, 220, 220, 255})
}

func (g *Game) drawPowerUp(pu PowerUp) {
	pulse := float32(0.85 + 0.15*math.Sin(g.Time*6))
	half := float32(pu.Size/2) * pulse
	xf := float32(pu.X)
	yf := float32(pu.Y)

	var col color.RGBA
	var label string
	switch pu.Type {
	case 0:
		col = color.RGBA{100, 255, 120, 255}
		label = "血"
	case 1:
		col = color.RGBA{100, 200, 255, 255}
		label = "盾"
	case 2:
		col = color.RGBA{255, 220, 100, 255}
		label = "金"
	case 3:
		col = color.RGBA{255, 150, 255, 255}
		label = "强"
	case 4:
		col = color.RGBA{120, 220, 255, 255}
		label = "电"
	}

	vector.DrawFilledCircle(g.offscreen, xf, yf, half+5, color.RGBA{col.R, col.G, col.B, 60}, false)
	g.fillTriangle(xf, yf-half, xf-half, yf, xf+half, yf, col)
	g.fillTriangle(xf-half, yf, xf+half, yf, xf, yf+half, color.RGBA{col.R / 2, col.G / 2, col.B / 2, 255})

	f := fontSmall()
	b := text.BoundString(f, label)
	text.Draw(g.offscreen, label, f, int(pu.X)-b.Dx()/2, int(pu.Y)+5, color.RGBA{255, 255, 255, 255})
}
