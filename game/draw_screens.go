package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawHUD() {
	f := fontSmall()
	p := g.Player

	scoreText := fmt.Sprintf("分数: %d", g.Score)
	text.Draw(g.offscreen, scoreText, f, 10, 15, color.RGBA{255, 255, 200, 255})

	hsText := fmt.Sprintf("最高: %d", g.HighScore)
	b := text.BoundString(f, hsText)
	text.Draw(g.offscreen, hsText, f, ScreenWidth-b.Dx()-10, 15, color.RGBA{255, 200, 100, 255})

	waveText := fmt.Sprintf("第 %d 波", g.Wave)
	b = text.BoundString(f, waveText)
	text.Draw(g.offscreen, waveText, f, ScreenWidth/2-b.Dx()/2, 15, color.RGBA{200, 150, 255, 255})

	barX := float32(10)
	barY := float32(22)
	barW := float32(140)
	barH := float32(10)
	hpRatio := float32(p.Health / p.MaxHealth)
	vector.DrawFilledRect(g.offscreen, barX, barY, barW, barH, color.RGBA{60, 30, 40, 200}, false)
	vector.DrawFilledRect(g.offscreen, barX, barY, barW*hpRatio, barH, color.RGBA{255, 80, 120, 255}, false)
	vector.StrokeRect(g.offscreen, barX, barY, barW, barH, 1, color.RGBA{200, 100, 140, 255}, false)

	hpText := fmt.Sprintf("%.0f/%.0f", p.Health, p.MaxHealth)
	b = text.BoundString(f, hpText)
	text.Draw(g.offscreen, hpText, f, 10+140/2-b.Dx()/2, 31, color.RGBA{255, 255, 255, 255})

	expBarY := float32(36)
	expRatio := float32(p.Exp) / float32(p.ExpToNext)
	vector.DrawFilledRect(g.offscreen, barX, expBarY, barW, 6, color.RGBA{30, 60, 40, 200}, false)
	vector.DrawFilledRect(g.offscreen, barX, expBarY, barW*expRatio, 6, color.RGBA{100, 255, 150, 255}, false)

	lvlText := fmt.Sprintf("等级 %d", p.Level)
	text.Draw(g.offscreen, lvlText, f, 155, 43, color.RGBA{150, 255, 200, 255})

	statY := 58
	text.Draw(g.offscreen, fmt.Sprintf("伤害: %d", p.Damage), f, 10, statY, color.RGBA{255, 180, 180, 255})
	text.Draw(g.offscreen, fmt.Sprintf("射速: %.1f/秒", 60.0/float64(p.FireRate)), f, 10, statY+15, color.RGBA{255, 220, 180, 255})
	text.Draw(g.offscreen, fmt.Sprintf("暴击: %d%%", int(p.CritChance*100)), f, 10, statY+30, color.RGBA{255, 180, 255, 255})

	text.Draw(g.offscreen, fmt.Sprintf("击杀: %d", p.Kills), f, ScreenWidth-100, statY, color.RGBA{200, 200, 255, 255})
	text.Draw(g.offscreen, fmt.Sprintf("金币: %d", p.Coins), f, ScreenWidth-100, statY+15, color.RGBA{255, 220, 100, 255})

	if p.Shield > 0 {
		t := fmt.Sprintf("护盾: %d秒", p.Shield/60)
		text.Draw(g.offscreen, t, f, ScreenWidth-100, statY+30, color.RGBA{100, 200, 255, 255})
	}
}

func (g *Game) drawPaused() {
	f := fontNormal()
	vector.DrawFilledRect(g.offscreen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 150}, false)
	t := "游戏暂停"
	b := text.BoundString(f, t)
	text.Draw(g.offscreen, t, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2, color.RGBA{255, 255, 255, 255})
	t2 := "按空格 或 鼠标右键 继续"
	b = text.BoundString(f, t2)
	text.Draw(g.offscreen, t2, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2+35, color.RGBA{200, 200, 200, 255})
}

func (g *Game) drawLevelUp() {
	f := fontNormal()
	f2 := fontSmall()
	vector.DrawFilledRect(g.offscreen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 200}, false)

	t := "升级了！"
	b := text.BoundString(f, t)
	text.Draw(g.offscreen, t, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2-130, color.RGBA{100, 255, 150, 255})

	t2 := "选择一项强化"
	b = text.BoundString(f2, t2)
	text.Draw(g.offscreen, t2, f2, ScreenWidth/2-b.Dx()/2, ScreenHeight/2-100, color.RGBA{200, 200, 255, 255})

	for i, u := range g.UpgradeChoices {
		btnW := float32(180)
		gap := float32(24)
		totalW := btnW*3 + gap*2
		startX := float32(ScreenWidth)/2 - totalW/2
		btnX := startX + float32(i)*(btnW+gap)
		btnY := float32(ScreenHeight/2 - 90)
		btnH := float32(180)

		var borderCol color.RGBA
		switch u.Rarity {
		case 1:
			borderCol = color.RGBA{150, 200, 255, 255}
		case 2:
			borderCol = color.RGBA{200, 150, 255, 255}
		case 3:
			borderCol = color.RGBA{255, 200, 100, 255}
		default:
			borderCol = color.RGBA{200, 200, 200, 255}
		}

		vector.DrawFilledRect(g.offscreen, btnX, btnY, btnW, btnH, color.RGBA{30, 30, 50, 230}, false)
		vector.StrokeRect(g.offscreen, btnX, btnY, btnW, btnH, 2, borderCol, false)

		name := u.Name
		b = text.BoundString(f, name)
		text.Draw(g.offscreen, name, f, int(btnX)+int(btnW)/2-b.Dx()/2, int(btnY)+35, borderCol)

		g.drawUpgradeIcon(btnX+btnW/2, btnY+85, u)

		desc := u.Desc
		b = text.BoundString(f2, desc)
		text.Draw(g.offscreen, desc, f2, int(btnX)+int(btnW)/2-b.Dx()/2, int(btnY)+135, color.RGBA{255, 255, 255, 255})

		rarityText := []string{"", "普通", "稀有", "传说"}[u.Rarity]
		b = text.BoundString(f2, rarityText)
		text.Draw(g.offscreen, rarityText, f2, int(btnX)+int(btnW)/2-b.Dx()/2, int(btnY)+165, borderCol)
	}
}

func (g *Game) drawUpgradeIcon(cx, cy float32, u Upgrade) {
	c := color.RGBA{100, 200, 255, 255}
	switch u.Rarity {
	case 2:
		c = color.RGBA{200, 150, 255, 255}
	case 3:
		c = color.RGBA{255, 200, 100, 255}
	}
	vector.DrawFilledCircle(g.offscreen, cx, cy, 22, color.RGBA{c.R, c.G, c.B, 60}, false)
	vector.DrawFilledCircle(g.offscreen, cx, cy, 15, c, false)
	vector.DrawFilledCircle(g.offscreen, cx, cy, 7, color.RGBA{255, 255, 255, 200}, false)
}

func (g *Game) drawGameOver() {
	f := fontNormal()
	f2 := fontSmall()
	vector.DrawFilledRect(g.offscreen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 180}, false)

	t := "游戏结束"
	b := text.BoundString(f, t)
	text.Draw(g.offscreen, t, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2-90, color.RGBA{255, 80, 80, 255})

	scoreText := fmt.Sprintf("最终分数: %d", g.Score)
	b = text.BoundString(f, scoreText)
	text.Draw(g.offscreen, scoreText, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2-50, color.RGBA{255, 255, 200, 255})

	hsText := fmt.Sprintf("最高分数: %d", g.HighScore)
	b = text.BoundString(f, hsText)
	text.Draw(g.offscreen, hsText, f, ScreenWidth/2-b.Dx()/2, ScreenHeight/2-20, color.RGBA{255, 200, 100, 255})

	lvlText := fmt.Sprintf("等级: %d  |  击杀: %d", g.Player.Level, g.Player.Kills)
	b = text.BoundString(f2, lvlText)
	text.Draw(g.offscreen, lvlText, f2, ScreenWidth/2-b.Dx()/2, ScreenHeight/2+10, color.RGBA{200, 200, 255, 255})

	btnX := float32(ScreenWidth/2 - 80)
	btnY := float32(ScreenHeight/2 + 40)
	btnW, btnH := float32(160), float32(45)
	vector.DrawFilledRect(g.offscreen, btnX, btnY, btnW, btnH, color.RGBA{80, 100, 150, 200}, false)
	vector.StrokeRect(g.offscreen, btnX, btnY, btnW, btnH, 2, color.RGBA{100, 200, 255, 255}, false)
	t = "重新开始"
	b = text.BoundString(f, t)
	text.Draw(g.offscreen, t, f, ScreenWidth/2-b.Dx()/2, int(btnY)+30, color.RGBA{255, 255, 255, 255})

	btnY2 := float32(ScreenHeight/2 + 105)
	vector.DrawFilledRect(g.offscreen, btnX, btnY2, btnW, btnH, color.RGBA{80, 100, 150, 200}, false)
	vector.StrokeRect(g.offscreen, btnX, btnY2, btnW, btnH, 2, color.RGBA{200, 150, 255, 255}, false)
	t = "返回主菜单"
	b = text.BoundString(f, t)
	text.Draw(g.offscreen, t, f, ScreenWidth/2-b.Dx()/2, int(btnY2)+30, color.RGBA{255, 255, 255, 255})
}
