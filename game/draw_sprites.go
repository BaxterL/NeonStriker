package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawPlayerShip(x, y float64, scale float64, alpha float32) {
	s := float32(scale)
	xf := float32(x)
	yf := float32(y)
	bc := color.RGBA{80, 200, 255, uint8(255 * alpha)}
	dc := color.RGBA{30, 110, 200, uint8(220 * alpha)}
	ac := color.RGBA{200, 245, 255, uint8(255 * alpha)}
	gc := color.RGBA{50, 150, 255, uint8(70 * alpha)}

	vector.DrawFilledCircle(g.offscreen, xf, yf, 30*s, gc, false)

	g.fillTriangle(xf, yf-24*s, xf-18*s, yf+2*s, xf+18*s, yf+2*s, bc)
	g.fillTriangle(xf-18*s, yf+2*s, xf+18*s, yf+2*s, xf, yf+18*s, dc)

	g.fillTriangle(xf-18*s, yf-2*s, xf-34*s, yf+18*s, xf-14*s, yf+14*s, bc)
	g.fillTriangle(xf+18*s, yf-2*s, xf+34*s, yf+18*s, xf+14*s, yf+14*s, bc)
	g.fillTriangle(xf-18*s, yf+2*s, xf-34*s, yf+18*s, xf-14*s, yf+14*s, dc)
	g.fillTriangle(xf-30*s, yf+16*s, xf-22*s, yf+14*s, xf-26*s, yf+20*s, dc)
	g.fillTriangle(xf+30*s, yf+16*s, xf+22*s, yf+14*s, xf+26*s, yf+20*s, dc)

	g.fillOval(xf, yf-4*s, 5*s, 9*s, ac)
	vector.DrawFilledCircle(g.offscreen, xf, yf-7*s, 2*s, color.RGBA{255, 255, 255, uint8(200 * alpha)}, false)

	flame := float32(10 + math.Sin(g.Time*28)*5)
	fc := color.RGBA{255, 140, 50, uint8(240 * alpha)}
	fc2 := color.RGBA{255, 220, 120, uint8(240 * alpha)}
	fc3 := color.RGBA{255, 255, 220, uint8(240 * alpha)}
	g.fillTriangle(xf-8*s, yf+16*s, xf+8*s, yf+16*s, xf, yf+16*s+flame*s, fc)
	g.fillTriangle(xf-4*s, yf+16*s, xf+4*s, yf+16*s, xf, yf+16*s+flame*0.7*s, fc2)
	g.fillTriangle(xf-2*s, yf+16*s, xf+2*s, yf+16*s, xf, yf+16*s+flame*0.4*s, fc3)

	g.fillTriangle(xf-10*s, yf+14*s, xf-7*s, yf+14*s, xf-9*s, yf+14*s+flame*0.5*s, fc)
	g.fillTriangle(xf+10*s, yf+14*s, xf+7*s, yf+14*s, xf+9*s, yf+14*s+flame*0.5*s, fc)
}

func (g *Game) drawEnemy(e Enemy) {
	xf := float32(e.X)
	yf := float32(e.Y)
	c := hsvToRgb(e.Hue, 0.75, 1.0)
	dc := hsvToRgb(e.Hue, 0.9, 0.7)
	gc := hsvToRgb(e.Hue, 0.5, 1.0)
	gc.A = 90

	switch e.EnemyType {
	case 0:
		g.drawEnemyBasic(xf, yf, e, c, dc, gc)
	case 1:
		g.drawEnemyElite(xf, yf, e, c, dc, gc)
	case 2:
		g.drawEnemyShooter(xf, yf, e, c, dc, gc)
	case 3:
		g.drawEnemySpinner(xf, yf, e, c, dc, gc)
	case 4:
		g.drawBoss(xf, yf, e, c, dc)
	case 5:
		g.drawEnemyRusher(xf, yf, e, c, dc, gc)
	}
}

func (g *Game) createEnemyDeathEffect(e Enemy) {
	base := int(e.Width / 10)
	if base < 4 {
		base = 4
	}
	if base > 16 {
		base = 16
	}
	if e.EnemyType == 0 || e.EnemyType == 5 {
		g.createHitParticles(e.X, e.Y, color.RGBA{255, 180, 120, 255}, base/2)
		return
	}
	if e.EnemyType == 1 {
		g.createExplosion(e.X, e.Y, 1)
		return
	}
	if e.EnemyType == 2 {
		g.createExplosion(e.X, e.Y, 2)
		return
	}
	if e.EnemyType == 3 {
		g.createExplosion(e.X, e.Y, 2)
		g.createHitParticles(e.X, e.Y, color.RGBA{255, 220, 120, 255}, base)
		return
	}
	g.createExplosion(e.X, e.Y, 3)
}

func (g *Game) drawEnemyBasic(xf, yf float32, e Enemy, c, dc, gc color.RGBA) {
	w := float32(e.Width / 2)
	h := float32(e.Height / 2)

	vector.DrawFilledCircle(g.offscreen, xf, yf, w+4, gc, false)

	g.fillTriangle(xf, yf-h, xf-w, yf, xf+w, yf, c)
	g.fillTriangle(xf-w, yf, xf+w, yf, xf, yf+h, dc)

	g.fillTriangle(xf-w*0.3, yf-h*0.3, xf-w, yf-2, xf-w*0.7, yf+h*0.2, dc)
	g.fillTriangle(xf+w*0.3, yf-h*0.3, xf+w, yf-2, xf+w*0.7, yf+h*0.2, dc)

	ey := yf - h*0.15
	vector.DrawFilledCircle(g.offscreen, xf-5, ey, 3.5, color.RGBA{255, 255, 255, 220}, false)
	vector.DrawFilledCircle(g.offscreen, xf+5, ey, 3.5, color.RGBA{255, 255, 255, 220}, false)
	vector.DrawFilledCircle(g.offscreen, xf-5, ey, 1.8, color.RGBA{255, 80, 80, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf+5, ey, 1.8, color.RGBA{255, 80, 80, 255}, false)

	g.fillTriangle(xf-3, yf+h*0.6, xf+3, yf+h*0.6, xf, yf+h+2, color.RGBA{255, 200, 100, 220})
}

func (g *Game) drawEnemyElite(xf, yf float32, e Enemy, c, dc, gc color.RGBA) {
	w := float32(e.Width / 2)
	h := float32(e.Height / 2)

	vector.DrawFilledCircle(g.offscreen, xf, yf, w+5, gc, false)

	pts := make([][2]float32, 6)
	for i := 0; i < 6; i++ {
		ang := float64(i) * math.Pi / 3
		pts[i] = [2]float32{float32(math.Cos(ang) * float64(w)), float32(math.Sin(ang) * float64(h))}
	}
	g.fillPoly(xf, yf, pts, c)

	innerPts := make([][2]float32, 6)
	for i := 0; i < 6; i++ {
		ang := float64(i)*math.Pi/3 + math.Pi/6
		innerPts[i] = [2]float32{float32(math.Cos(ang) * float64(w*0.65)), float32(math.Sin(ang) * float64(h*0.65))}
	}
	g.fillPoly(xf, yf, innerPts, dc)

	for i := 0; i < 6; i++ {
		ang := float64(i) * math.Pi / 3
		ox := float32(math.Cos(ang)) * w
		oy := float32(math.Sin(ang)) * h
		ix := float32(math.Cos(ang)) * w * 0.6
		iy := float32(math.Sin(ang)) * h * 0.6
		g.fillTriangle(ox, oy, ix-2, iy-2, ix+2, iy+2, c)
	}

	corePulse := float32(0.7 + 0.3*math.Sin(g.Time*5))
	vector.DrawFilledCircle(g.offscreen, xf, yf, 9*corePulse, color.RGBA{255, 255, 220, 220}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 4.5, color.RGBA{255, 255, 255, 255}, false)

	g.drawEnemyHpBar(e, xf, yf)
}

func (g *Game) drawEnemyShooter(xf, yf float32, e Enemy, c, dc, gc color.RGBA) {
	w := float32(e.Width / 2)
	h := float32(e.Height / 2)

	vector.DrawFilledCircle(g.offscreen, xf, yf, w+7, gc, false)

	g.fillOval(xf, yf, w*0.85, h*0.7, c)

	ring := 10
	for i := 0; i < ring; i++ {
		ang := float64(i) * 2 * math.Pi / float64(ring)
		ox := float32(math.Cos(ang)) * w * 0.9
		oy := float32(math.Sin(ang)) * h * 0.8
		vector.DrawFilledCircle(g.offscreen, xf+ox, yf+oy, 3, dc, false)
	}

	barrelAngles := []float64{-math.Pi / 2, -math.Pi / 6, math.Pi / 6}
	for _, ang := range barrelAngles {
		bx := float32(math.Cos(ang)) * w * 0.55
		by := float32(math.Sin(ang)) * h * 0.55
		tx := float32(math.Cos(ang)) * w * 1.05
		ty := float32(math.Sin(ang)) * h * 1.05
		g.fillTriangle(bx-4, by, bx+4, by, tx, ty, dc)
		g.fillTriangle(bx-2, by, bx+2, by, tx*0.92, ty*0.92, color.RGBA{40, 40, 60, 220})
	}

	corePulse := float32(0.7 + 0.3*math.Sin(g.Time*4))
	vector.DrawFilledCircle(g.offscreen, xf, yf, 12*corePulse, color.RGBA{255, 200, 80, 200}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 7, color.RGBA{255, 255, 200, 240}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 3, color.RGBA{255, 255, 255, 255}, false)

	g.drawEnemyHpBar(e, xf, yf)
}

func (g *Game) drawEnemySpinner(xf, yf float32, e Enemy, c, dc, gc color.RGBA) {
	w := float32(e.Width / 2)
	h := float32(e.Height / 2)

	vector.DrawFilledCircle(g.offscreen, xf, yf, w+10, gc, false)

	spikes := 10
	pts := make([][2]float32, spikes*2)
	for i := 0; i < spikes*2; i++ {
		ang := float64(i)*math.Pi/float64(spikes) + g.Time*1.5
		r := float64(w)
		if i%2 == 1 {
			r *= 0.7
		}
		ry := float64(h)
		if i%2 == 1 {
			ry *= 0.7
		}
		pts[i] = [2]float32{float32(math.Cos(ang) * r), float32(math.Sin(ang) * ry)}
	}
	g.fillPoly(xf, yf, pts, c)

	innerSpikes := 8
	innerPts := make([][2]float32, innerSpikes*2)
	for i := 0; i < innerSpikes*2; i++ {
		ang := float64(i)*math.Pi/float64(innerSpikes) - g.Time*2
		r := float64(w * 0.55)
		if i%2 == 1 {
			r *= 0.75
		}
		ry := float64(h * 0.55)
		if i%2 == 1 {
			ry *= 0.75
		}
		innerPts[i] = [2]float32{float32(math.Cos(ang) * r), float32(math.Sin(ang) * ry)}
	}
	g.fillPoly(xf, yf, innerPts, dc)

	eyeCount := 6
	for i := 0; i < eyeCount; i++ {
		ang := g.Time*1.2 + float64(i)*2*math.Pi/float64(eyeCount)
		ox := float32(math.Cos(ang)) * w * 0.35
		oy := float32(math.Sin(ang)) * h * 0.35
		vector.DrawFilledCircle(g.offscreen, xf+ox, yf+oy, 3.5, color.RGBA{255, 240, 240, 230}, false)
		vector.DrawFilledCircle(g.offscreen, xf+ox, yf+oy, 1.8, color.RGBA{255, 60, 60, 255}, false)
	}

	vector.DrawFilledCircle(g.offscreen, xf, yf, 11, color.RGBA{255, 100, 100, 250}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 6, color.RGBA{255, 200, 200, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 3, color.RGBA{255, 255, 255, 255}, false)

	g.drawEnemyHpBar(e, xf, yf)
}

func (g *Game) drawEnemyRusher(xf, yf float32, e Enemy, c, dc, gc color.RGBA) {
	w := float32(e.Width / 2)
	h := float32(e.Height / 2)

	vector.DrawFilledCircle(g.offscreen, xf, yf, w+5, gc, false)

	diamond := [][2]float32{
		{0, -h},
		{w * 0.75, 0},
		{0, h},
		{-w * 0.75, 0},
	}
	g.fillPoly(xf, yf, diamond, c)

	innerDiamond := [][2]float32{
		{0, -h * 0.55},
		{w * 0.4, 0},
		{0, h * 0.55},
		{-w * 0.4, 0},
	}
	g.fillPoly(xf, yf, innerDiamond, dc)

	finAngles := []float64{
		-math.Pi / 2,
		math.Pi / 2,
		0,
		math.Pi,
	}
	for _, ang := range finAngles {
		fx := float32(math.Cos(ang)) * w * 0.65
		fy := float32(math.Sin(ang)) * h * 0.65
		tx := float32(math.Cos(ang)) * w * 1.1
		ty := float32(math.Sin(ang)) * h * 1.1
		px := float32(math.Cos(ang+math.Pi/2)) * w * 0.15
		py := float32(math.Sin(ang+math.Pi/2)) * h * 0.15
		g.fillTriangle(tx, ty, fx+px, fy+py, fx-px, fy-py, dc)
	}

	vector.DrawFilledCircle(g.offscreen, xf, yf, 5, color.RGBA{255, 255, 200, 230}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, 2.5, color.RGBA{255, 80, 80, 255}, false)
}

func (g *Game) drawBoss(xf, yf float32, e Enemy, c, dc color.RGBA) {
	bw := float32(e.Width / 2)
	bh := float32(e.Height / 2)

	glowPulse := float32(0.85 + 0.15*math.Sin(g.Time*3))
	vector.DrawFilledCircle(g.offscreen, xf, yf, (bw+25)*glowPulse, color.RGBA{255, 50, 100, 50}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf, bw+12, color.RGBA{255, 80, 120, 70}, false)

	bodyPts := [][2]float32{
		{0, -bh},
		{bw * 0.9, -bh * 0.5},
		{bw, bh * 0.2},
		{bw * 0.7, bh},
		{-bw * 0.7, bh},
		{-bw, bh * 0.2},
		{-bw * 0.9, -bh * 0.5},
	}
	g.fillPoly(xf, yf, bodyPts, c)

	armorPts := [][2]float32{
		{0, -bh * 0.85},
		{bw * 0.7, -bh * 0.4},
		{bw * 0.75, bh * 0.1},
		{0, bh * 0.5},
		{-bw * 0.75, bh * 0.1},
		{-bw * 0.7, -bh * 0.4},
	}
	g.fillPoly(xf, yf, armorPts, dc)

	g.fillTriangle(xf-bw, yf-bh*0.3, xf-bw-30, yf-15, xf-bw-30, yf+20, dc)
	g.fillTriangle(xf+bw, yf-bh*0.3, xf+bw+30, yf-15, xf+bw+30, yf+20, dc)
	g.fillTriangle(xf-bw-5, yf-12, xf-bw-25, yf-8, xf-bw-20, yf+15, c)
	g.fillTriangle(xf+bw+5, yf-12, xf+bw+25, yf-8, xf+bw+20, yf+15, c)

	vector.DrawFilledCircle(g.offscreen, xf-bw-30, yf+5, 5, color.RGBA{255, 200, 100, 230}, false)
	vector.DrawFilledCircle(g.offscreen, xf+bw+30, yf+5, 5, color.RGBA{255, 200, 100, 230}, false)
	vector.DrawFilledCircle(g.offscreen, xf-bw-30, yf+5, 2.5, color.RGBA{255, 80, 80, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf+bw+30, yf+5, 2.5, color.RGBA{255, 80, 80, 255}, false)

	corePulse := float32(0.75 + 0.25*math.Sin(g.Time*4))
	vector.DrawFilledCircle(g.offscreen, xf, yf+bh*0.15, 22*corePulse, color.RGBA{255, 255, 180, 180}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf+bh*0.15, 14, color.RGBA{255, 255, 220, 230}, false)
	vector.DrawFilledCircle(g.offscreen, xf, yf+bh*0.15, 7, color.RGBA{255, 255, 255, 255}, false)

	eyeY := yf - bh*0.25
	eyeDX := bw * 0.35
	vector.DrawFilledCircle(g.offscreen, xf-eyeDX, eyeY, 8, color.RGBA{255, 240, 220, 240}, false)
	vector.DrawFilledCircle(g.offscreen, xf+eyeDX, eyeY, 8, color.RGBA{255, 240, 220, 240}, false)
	vector.DrawFilledCircle(g.offscreen, xf-eyeDX, eyeY, 4.5, color.RGBA{255, 60, 80, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf+eyeDX, eyeY, 4.5, color.RGBA{255, 60, 80, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf-eyeDX-1.5, eyeY-1.5, 1.5, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(g.offscreen, xf+eyeDX-1.5, eyeY-1.5, 1.5, color.RGBA{255, 255, 255, 255}, false)

	hornPts := [][2]float32{
		{-bw * 0.4, -bh},
		{-bw * 0.25, -bh - 20},
		{-bw * 0.1, -bh},
	}
	g.fillPoly(xf, yf, hornPts, dc)
	hornPts2 := [][2]float32{
		{bw * 0.4, -bh},
		{bw * 0.25, -bh - 20},
		{bw * 0.1, -bh},
	}
	g.fillPoly(xf, yf, hornPts2, dc)

	g.drawBossHpBar(e)
}

func (g *Game) drawSlashWave(s SlashWave) {
	alpha := float32(s.Life) / 60.0
	outer := float32(s.Radius)
	inner := outer - 18
	if inner < 0 {
		inner = 0
	}
	mid := (outer + inner) / 2
	center := color.RGBA{s.Color.R, s.Color.G, s.Color.B, uint8(150 * alpha)}
	soft := color.RGBA{s.Color.R, s.Color.G, s.Color.B, uint8(60 * alpha)}

	vector.StrokeCircle(g.offscreen, float32(s.X), float32(s.Y), outer, 4, soft, false)
	vector.StrokeCircle(g.offscreen, float32(s.X), float32(s.Y), mid, 3, center, false)
	vector.StrokeCircle(g.offscreen, float32(s.X), float32(s.Y), inner, 1.5, color.RGBA{s.Color.R, s.Color.G, s.Color.B, uint8(100 * alpha)}, false)

	x1 := float32(s.X + math.Cos(s.Angle)*float64(outer))
	y1 := float32(s.Y + math.Sin(s.Angle)*float64(outer))
	x2 := float32(s.X - math.Cos(s.Angle)*float64(outer))
	y2 := float32(s.Y - math.Sin(s.Angle)*float64(outer))
	vector.DrawFilledCircle(g.offscreen, x1, y1, 9, color.RGBA{255, 255, 255, uint8(120 * alpha)}, false)
	vector.DrawFilledCircle(g.offscreen, x1, y1, 5, color.RGBA{255, 255, 255, uint8(200 * alpha)}, false)
	vector.DrawFilledCircle(g.offscreen, x2, y2, 9, color.RGBA{255, 255, 255, uint8(120 * alpha)}, false)
	vector.DrawFilledCircle(g.offscreen, x2, y2, 5, color.RGBA{255, 255, 255, uint8(200 * alpha)}, false)

	perps := 5
	for i := 0; i < perps; i++ {
		t := float64(i)/float64(perps-1) - 0.5
		cx := float32(s.X + math.Cos(s.Angle)*float64(mid)*t*2)
		cy := float32(s.Y + math.Sin(s.Angle)*float64(mid)*t*2)
		size := float32(3 + math.Abs(t)*4)
		vector.DrawFilledCircle(g.offscreen, cx, cy, size, color.RGBA{255, 255, 255, uint8(80 * alpha)}, false)
	}
}
