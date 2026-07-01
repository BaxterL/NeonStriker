package game

import (
	"image/color"
	"math"
	"math/rand"
)

func (g *Game) fireWeapons() {
	p := g.Player
	baseDmg := p.Damage
	crit := rand.Float64() < p.CritChance
	if crit {
		baseDmg = int(float64(baseDmg) * p.CritMult)
	}

	bulletCol := color.RGBA{100, 255, 200, 255}
	if crit {
		bulletCol = color.RGBA{255, 220, 100, 255}
	}

	bulletSize := 4.0
	if p.Spread > 10 {
		bulletSize = 5.0
	}

	makeB := func(x, y, vx, vy float64, dmg int, c color.RGBA, sz float64, pierce int) Bullet {
		return Bullet{X: x, Y: y, VX: vx, VY: vy, Damage: dmg, Player: true, Size: sz, Color: c, Life: 120, Pierce: pierce}
	}

	speed := p.BulletSpeed
	count := p.BulletCount
	spreadDeg := float64(p.Spread)

	if p.LaserBeam {
		g.Bullets = append(g.Bullets, makeB(p.X, p.Y-p.Height/2, 0, -speed*1.8, baseDmg*2, color.RGBA{255, 80, 200, 255}, 10, 6))
		return
	}

	if p.FrontArc {
		for i := -2; i <= 2; i++ {
			angle := float64(i) * 12 * math.Pi / 180
			vx := math.Sin(angle) * speed
			vy := -math.Cos(angle) * speed
			g.Bullets = append(g.Bullets, makeB(p.X, p.Y-p.Height/2, vx, vy, baseDmg, bulletCol, bulletSize, 1))
		}
	}

	if count == 1 {
		g.Bullets = append(g.Bullets, makeB(p.X, p.Y-p.Height/2, 0, -speed, baseDmg, bulletCol, bulletSize, 0))
	} else {
		var totalSpread float64
		if count <= 3 {
			totalSpread = spreadDeg + float64(count-1)*10
		} else {
			totalSpread = spreadDeg + 20 + float64(count-3)*20
		}
		step := totalSpread / float64(count-1)
		startAngle := -totalSpread / 2
		for i := 0; i < count; i++ {
			angle := (startAngle + step*float64(i)) * math.Pi / 180
			vx := math.Sin(angle) * speed
			vy := -math.Cos(angle) * speed
			g.Bullets = append(g.Bullets, makeB(p.X, p.Y-p.Height/2, vx, vy, baseDmg, bulletCol, bulletSize, 0))
		}
	}

	// 左前/右前射口，增强扇面输出感。
	for i := 0; i < p.SideGuns; i++ {
		offset := float64(i+1) * 12
		g.Bullets = append(g.Bullets, makeB(p.X-offset, p.Y, -speed*0.3, -speed*0.95, baseDmg, color.RGBA{150, 200, 255, 255}, 3, 0))
		g.Bullets = append(g.Bullets, makeB(p.X+offset, p.Y, speed*0.3, -speed*0.95, baseDmg, color.RGBA{150, 200, 255, 255}, 3, 0))
	}

	if p.BackGun {
		g.Bullets = append(g.Bullets, makeB(p.X, p.Y+p.Height/2, 0, speed, baseDmg, color.RGBA{200, 150, 255, 255}, 4, 0))
	}

	if p.TwinShip {
		g.Bullets = append(g.Bullets, makeB(p.X-20, p.Y-p.Height/2, 0, -speed*1.1, baseDmg, color.RGBA{100, 255, 240, 255}, 4, 0))
		g.Bullets = append(g.Bullets, makeB(p.X+20, p.Y-p.Height/2, 0, -speed*1.1, baseDmg, color.RGBA{100, 255, 240, 255}, 4, 0))
	}

	if p.LaserLevel > 0 && !p.LaserBeam {
		laserDmg := p.LaserLevel * 3
		if crit {
			laserDmg = int(float64(laserDmg) * p.CritMult)
		}
		g.Bullets = append(g.Bullets, makeB(p.X, p.Y-p.Height/2, 0, -speed*1.5, laserDmg, color.RGBA{255, 100, 200, 255}, 8, 5))
	}
}

func (g *Game) updateBullets() {
	active := g.Bullets[:0]
	for _, b := range g.Bullets {
		b.X += b.VX
		b.Y += b.VY
		b.Life--
		if b.Life > 0 && b.Y > -30 && b.Y < ScreenHeight+30 && b.X > -30 && b.X < ScreenWidth+30 {
			active = append(active, b)
		}
	}
	g.Bullets = active
}

func (g *Game) spawnEnemies() {
	if g.BossActive {
		return
	}

	// 第一首领击败前，节奏更平缓；之后开始明显加速并引入新敌机类型。
	waveMul := waveGrowth(g.Wave) * spawnGrowthMultiplier(g.Wave)
	baseInterval := 48
	interval := float64(baseInterval) / waveMul
	if g.BossDefeated > 0 {
		interval = interval * 0.82
	}
	if interval < 3 {
		interval = 3
	}

	g.EnemySpawnTimer++
	if float64(g.EnemySpawnTimer) < interval {
		return
	}
	g.EnemySpawnTimer = 0

	roll := rand.Float64()
	enemyType := 0
	if g.BossDefeated == 0 {
		if g.Wave >= 3 && roll < 0.18 {
			enemyType = 1
		}
	} else {
		if g.Wave >= 2 && roll < 0.12 {
			enemyType = 1
		} else if g.Wave >= 4 && roll < 0.22 {
			enemyType = 2
		} else if g.Wave >= 6 && roll < 0.30 {
			enemyType = 3
		} else if g.Wave >= 8 && roll < 0.38 {
			enemyType = 5
		}
	}

	hpMult := waveGrowth(g.Wave) * enemyGrowthMultiplier(g.Wave)
	if g.Wave >= 7 {
		hpMult *= 1.6 + float64(g.Wave-7)*0.35
	}
	spdMult := 1.0 + float64(g.Wave-1)*0.06
	if g.BossDefeated > 0 {
		spdMult += 0.1 + float64(g.BossDefeated)*0.03
	}

	switch enemyType {
	case 0:
		hp := int(2 + 3*hpMult)
		g.Enemies = append(g.Enemies, Enemy{
			X: randRange(40, ScreenWidth-40),
			Y: -30, Width: 34, Height: 34,
			Health: hp, MaxHealth: hp,
			VX: (rand.Float64()-0.5)*1.2,
			VY: (1.6 + rand.Float64()) * spdMult,
			Score: 5 + g.Wave/2, ExpReward: 2 + g.Wave*3/4, CoinDrop: 1,
			EnemyType: 0, Hue: rand.Float64() * 360,
		})
	case 1:
		hp := int(8 * hpMult)
		if g.BossDefeated >= 3 {
			hp = max(hp, g.eliteHPFromPlayerPower(2.0))
		}
		g.Enemies = append(g.Enemies, Enemy{
			X: randRange(50, ScreenWidth-50),
			Y: -40, Width: 44, Height: 44,
			Health: hp, MaxHealth: hp,
			VX: (rand.Float64()-0.5)*3.2,
			VY: (1.3 + rand.Float64()*0.6) * spdMult,
			Score: 16 + g.Wave,
			ExpReward: 4 + g.Wave, CoinDrop: 3,
			EnemyType: 1, Hue: 280 + rand.Float64()*40,
		})
	case 2:
		hp := int(20 * hpMult)
		if g.BossDefeated >= 3 {
			hp = max(hp, g.eliteHPFromPlayerPower(4.0))
		}
		g.Enemies = append(g.Enemies, Enemy{
			X: randRange(50, ScreenWidth-50),
			Y: -50, Width: 54, Height: 54,
			Health: hp, MaxHealth: hp,
			VX: math.Sin(rand.Float64()*math.Pi*2) * 1.5,
			VY: (0.8 + rand.Float64()*0.4) * spdMult,
			Score: 40 + g.Wave*2, ExpReward: 9 + g.Wave*3/2, CoinDrop: 6,
			EnemyType: 2, FireCd: 120, FireTimer: 60,
			Hue: 200 + rand.Float64()*30,
		})
	case 3:
		hp := int(45 * hpMult)
		if g.BossDefeated >= 3 {
			hp = max(hp, g.eliteHPFromPlayerPower(6.0))
		}
		g.Enemies = append(g.Enemies, Enemy{
			X: randRange(60, ScreenWidth-60),
			Y: -60, Width: 70, Height: 70,
			Health: hp, MaxHealth: hp,
			VX: 0,
			VY: 0.6 * spdMult,
			Score: 90 + g.Wave*3, ExpReward: 14 + g.Wave*2, CoinDrop: 10,
			EnemyType: 3, FireCd: 90, FireTimer: 45,
			Hue: 30 + rand.Float64()*20,
		})
	case 5:
		// 四面来敌：后期从边缘刷出，形成包围压迫感。
		hp := int(30 * hpMult)
		side := rand.Intn(4)
		x, y := 0.0, 0.0
		vx, vy := 0.0, 0.0
		switch side {
		case 0:
			x = -40
			y = randRange(60, ScreenHeight-60)
			vx = (2.5 + rand.Float64()) * spdMult
		case 1:
			x = ScreenWidth + 40
			y = randRange(60, ScreenHeight-60)
			vx = -(2.5 + rand.Float64()) * spdMult
		case 2:
			x = randRange(60, ScreenWidth-60)
			y = -40
			vy = (2.0 + rand.Float64()) * spdMult
		default:
			x = randRange(60, ScreenWidth-60)
			y = ScreenHeight + 40
			vy = -(2.0 + rand.Float64()) * spdMult
		}
		g.Enemies = append(g.Enemies, Enemy{
			X: x, Y: y, Width: 36, Height: 36,
			Health: hp, MaxHealth: hp,
			VX: vx, VY: vy,
			Score: 70 + g.Wave*2, ExpReward: 8 + g.Wave*3/2, CoinDrop: 4,
			EnemyType: 5, FireCd: 70, FireTimer: 30,
			Hue: 180 + rand.Float64()*80,
		})
	}
}

func (g *Game) updateEnemies() {
	active := g.Enemies[:0]
	for _, e := range g.Enemies {
		if e.EnemyType == 1 {
			if e.X < e.Width/2 || e.X > ScreenWidth-e.Width/2 {
				e.VX = -e.VX
			}
		}
		if e.EnemyType == 2 {
			// 射手怪：简单定点追踪弹。
			e.VX = math.Sin(g.Time+e.X) * 1.5
			e.FireTimer--
			if e.FireTimer <= 0 && e.Y > 0 && e.Y < ScreenHeight*0.6 {
				e.FireTimer = e.FireCd
				dx := g.Player.X - e.X
				dy := g.Player.Y - e.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if g.Wave >= 8 && rand.Float64() < 0.35 {
					g.Bullets = append(g.Bullets, Bullet{X: e.X, Y: e.Y, VX: dx / dist * 8, VY: dy / dist * 8, Damage: 14, Player: false, Size: 12, Color: color.RGBA{255, 80, 180, 255}, Life: 38, Beam: true})
					continue
				}
				sp := 3.5
				g.Bullets = append(g.Bullets, Bullet{
					X: e.X, Y: e.Y,
					VX: dx / dist * sp,
					VY: dy / dist * sp,
					Damage: 8, Player: false,
					Size: 6, Color: color.RGBA{255, 150, 50, 255}, Life: 300,
				})
			}
		}
		if e.EnemyType == 5 {
			// 边缘冲锋怪：专门从四周夹击玩家。
			e.FireTimer--
			if e.FireTimer <= 0 {
				e.FireTimer = e.FireCd
				// 四向冲击弹，模拟雷霆战机式压迫。
				for i := 0; i < 4; i++ {
					angle := float64(i) * math.Pi / 2
					g.Bullets = append(g.Bullets, Bullet{
						X: e.X, Y: e.Y,
						VX: math.Cos(angle) * 3.2,
						VY: math.Sin(angle) * 3.2,
						Damage: 9, Player: false,
						Size: 7, Color: color.RGBA{255, 120, 220, 255}, Life: 420,
					})
				}
			}
		}
		if e.EnemyType == 3 {
			// 旋转怪：制造扇形/环形弹幕压力。
			e.VX = math.Cos(g.Time*0.8+e.X) * 2
			e.FireTimer--
			if e.FireTimer <= 0 && e.Y > 0 && e.Y < ScreenHeight*0.5 {
				e.FireTimer = e.FireCd
				for i := 0; i < 8; i++ {
					angle := float64(i) * math.Pi * 2 / 8
					sp := 2.5
					g.Bullets = append(g.Bullets, Bullet{
						X: e.X, Y: e.Y,
						VX: math.Cos(angle) * sp,
						VY: math.Sin(angle) * sp,
						Damage: 10, Player: false,
						Size: 7, Color: color.RGBA{255, 80, 80, 255}, Life: 400,
					})
				}
			}
		}
		if e.EnemyType == 4 {
			// 首领：边移动边持续弹幕压制。
			if e.Y < 120 {
				e.VY = 1.2
			} else {
				e.VY = 0.15
			}
			e.FireTimer--
			if e.FireTimer <= 0 {
				e.FireTimer = 60
				angles := 12
				for i := 0; i < angles; i++ {
					angle := float64(i)*math.Pi*2/float64(angles) + g.Time*0.3
					sp := 3.0
					g.Bullets = append(g.Bullets, Bullet{
						X: e.X, Y: e.Y,
						VX: math.Cos(angle) * sp,
						VY: math.Sin(angle) * sp,
						Damage: 12, Player: false,
						Size: 8, Color: color.RGBA{255, 50, 150, 255}, Life: 500,
					})
				}
			}
		}

		e.X += e.VX
		e.Y += e.VY

		if e.Y < ScreenHeight+80 {
			active = append(active, e)
		}
	}
	g.Enemies = active
}

func (g *Game) updateWave() {
	if g.BossActive {
		return
	}

	// 第一首领击败后开始明显膨胀，数值采用三次曲线。
	if g.BossDefeated == 0 {
		if g.WaveTimer >= g.WaveDuration {
			g.spawnBoss()
		}
		return
	}

	if g.WaveTimer >= g.WaveDuration {
		g.Wave++
		g.WaveTimer = 0
		g.WaveDuration = int(60 * (22 + math.Min(26, float64(g.Wave)*0.85)))
		g.BattlefieldScale = 1 + math.Min(1.8, float64(g.Wave-1)*0.10)
		if g.WaveDuration < 60*18 {
			g.WaveDuration = 60 * 18
		}
		if g.Wave%3 == 0 {
			g.spawnBoss()
		}
	}
}

func (g *Game) spawnBoss() {
	g.BossActive = true
	f := waveGrowth(g.Wave) * enemyGrowthMultiplier(g.Wave)
	waveFloor := int(900*f) + g.BossDefeated*500
	bossHp := waveFloor
	if g.BossDefeated >= 3 {
		bossHp = max(waveFloor, g.bossHPFromPlayerPower())
	}
	g.Enemies = append(g.Enemies, Enemy{
		X: ScreenWidth / 2, Y: -100,
		Width: 180, Height: 120,
		Health: bossHp, MaxHealth: bossHp,
		VX: 0, VY: 0,
		Score: 300 + g.Wave*20, ExpReward: 60 + g.Wave*3, CoinDrop: 60,
		EnemyType: 4, FireCd: 0, FireTimer: 60,
		Hue: 340,
	})
	g.ScreenShake = 20
}

func (g *Game) bossHPFromPlayerPower() int {
	power := g.playerPowerEstimate()
	return int(power * (12 + float64(g.Wave)*0.8))
}

func (g *Game) eliteHPFromPlayerPower(mult float64) int {
	power := g.playerPowerEstimate()
	return int(power * mult * 1.2)
}

func (g *Game) playerPowerEstimate() float64 {
	p := g.Player
	shots := float64(p.BulletCount + p.SideGuns*2)
	if p.BackGun {
		shots += 1
	}
	if p.TwinShip {
		shots += 2.5
	}
	if p.FrontArc {
		shots += 3
	}
	if p.LaserBeam {
		shots += 5
	}
	if p.ChainLightning {
		shots *= 1.25
	}
	crit := 1 + p.CritChance*(p.CritMult-1)
	return float64(p.Damage) * shots * crit
}

func (g *Game) updateCoins() {
	p := g.Player
	active := g.Coins[:0]
	for _, c := range g.Coins {
		dx := p.X - c.X
		dy := p.Y - c.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		// 拾取范围由磁力属性控制，堆两张磁力装置后基本覆盖全屏。
		if dist < p.PickupRange && dist > 0.1 {
			pull := 0.45
			if dist < p.PickupRange*2.5 {
				pull = 0.9
			}
			c.VX += dx / dist * 12 * pull
			c.VY += dy / dist * 12 * pull
		}

		c.X += c.VX
		c.Y += c.VY
		c.VX *= 0.9
		c.VY *= 0.9

		if dist < 20 {
			p.Coins += c.Value
			g.Meta.Coins += c.Value
			g.saveMeta()
			g.Score += c.Value
			continue
		}
		if c.Y < ScreenHeight+50 && c.X > -50 && c.X < ScreenWidth+50 {
			active = append(active, c)
		}
	}
	g.Coins = active
}

func (g *Game) updateExpOrbs() {
	p := g.Player
	active := g.ExpOrbs[:0]
	for _, e := range g.ExpOrbs {
		dx := p.X - e.X
		dy := p.Y - e.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		// 拾取范围由磁力属性控制，堆两张磁力装置后基本覆盖全屏。
		if dist < p.PickupRange && dist > 0.1 {
			pull := 0.55
			if dist < p.PickupRange*2.5 {
				pull = 1.0
			}
			e.VX += dx / dist * 14 * pull
			e.VY += dy / dist * 14 * pull
		}

		e.X += e.VX
		e.Y += e.VY
		e.VX *= 0.88
		e.VY *= 0.88

		if dist < 20 {
			g.addExp(e.Value)
			if g.State != StatePlaying {
				active = append(active, e)
			}
			continue
		}
		if e.Y < ScreenHeight+50 && e.X > -50 && e.X < ScreenWidth+50 {
			active = append(active, e)
		}
	}
	g.ExpOrbs = active
}

func (g *Game) updateSlashWaves() {
	active := g.SlashWaves[:0]
	for _, s := range g.SlashWaves {
		s.Radius += s.Expand
		s.Angle += s.Spin
		s.Life--
		if s.Life > 0 {
			active = append(active, s)
		}
	}
	g.SlashWaves = active
}

func (g *Game) updatePowerUps() {
	active := g.PowerUps[:0]
	for _, p := range g.PowerUps {
		p.Y += 1.5
		if p.Y < ScreenHeight+30 {
			active = append(active, p)
		}
	}
	g.PowerUps = active
}

func (g *Game) updateParticles() {
	active := g.Particles[:0]
	for _, p := range g.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.VY += 0.05
		p.VX *= 0.98
		p.Life--
		if p.Life > 0 {
			active = append(active, p)
		}
	}
	g.Particles = active
}
