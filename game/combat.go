package game

import (
	"image/color"
	"math"
	"math/rand"
)

func (g *Game) checkCollisions() {
	p := g.Player

	for i := len(g.Bullets) - 1; i >= 0; i-- {
		b := g.Bullets[i]
		if !b.Player {
			continue
		}
		hit := false
		for j := len(g.Enemies) - 1; j >= 0; j-- {
			e := g.Enemies[j]
			if math.Abs(b.X-e.X) < (e.Width/2+b.Size) && math.Abs(b.Y-e.Y) < (e.Height/2+b.Size) {
				g.Enemies[j].Health -= b.Damage
				g.createHitParticles(b.X, b.Y, b.Color, 4)
				if g.Player.ChainLightning {
					g.chainLightningHit(e.X, e.Y, max(1, b.Damage/2))
				}

				if b.Pierce > 0 {
					b.Pierce--
					g.Bullets[i] = b
				} else {
					g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
					hit = true
				}

				if g.Enemies[j].Health <= 0 {
					g.onEnemyDeath(j)
				}
				if hit {
					break
				}
			}
		}
	}

	if p.Invincible <= 0 {
		for i := len(g.Bullets) - 1; i >= 0; i-- {
			b := g.Bullets[i]
			if b.Player {
				continue
			}
			if math.Abs(b.X-p.X) < p.Width/2 && math.Abs(b.Y-p.Y) < p.Height/2 {
				g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
				g.damagePlayer(float64(b.Damage))
			}
		}
	}

	if p.Invincible <= 0 {
		for j := len(g.Enemies) - 1; j >= 0; j-- {
			e := g.Enemies[j]
			if math.Abs(e.X-p.X) < (e.Width+p.Width)/2-5 && math.Abs(e.Y-p.Y) < (e.Height+p.Height)/2-5 {
				dmg := 10.0
				if e.EnemyType >= 2 {
					dmg = 20
				}
				if e.EnemyType == 4 {
					dmg = 30
				}
				if e.EnemyType < 4 {
					g.createEnemyDeathEffect(e)
					g.onEnemyDeath(j)
				}
				g.damagePlayer(dmg)
				break
			}
		}
	}

	for i := len(g.SlashWaves) - 1; i >= 0; i-- {
		s := g.SlashWaves[i]
		dx := p.X - s.X
		dy := p.Y - s.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > s.Radius-25 && dist < s.Radius+20 {
			ang := math.Atan2(dy, dx)
			delta := math.Mod(math.Abs(ang-s.Angle), math.Pi*2)
			if delta < 0.65 || math.Abs(delta-math.Pi) < 0.65 {
				g.damagePlayer(float64(s.Damage))
				g.SlashWaves[i].HitOnce = true
			}
		}
		if s.HitOnce {
			g.SlashWaves[i].Life = 0
		}
	}

	for i := len(g.PowerUps) - 1; i >= 0; i-- {
		pu := g.PowerUps[i]
		if math.Abs(pu.X-p.X) < (pu.Size+p.Width)/2 && math.Abs(pu.Y-p.Y) < (pu.Size+p.Height)/2 {
			g.applyPowerUp(pu.Type)
			g.createHitParticles(pu.X, pu.Y, color.RGBA{100, 255, 200, 255}, 12)
			g.PowerUps = append(g.PowerUps[:i], g.PowerUps[i+1:]...)
		}
	}
}

func (g *Game) onEnemyDeath(idx int) {
	e := g.Enemies[idx]
	g.Player.Kills++
	g.Score += e.Score

	for i := 0; i < e.CoinDrop; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 1 + rand.Float64()*3
		g.Coins = append(g.Coins, Coin{
			X: e.X, Y: e.Y,
			VX: math.Cos(angle) * speed,
			VY: math.Sin(angle) * speed - 1,
			Size: 8, Value: 1,
		})
	}

	g.ExpOrbs = append(g.ExpOrbs, ExpOrb{
		X: e.X, Y: e.Y,
		VX: (rand.Float64()-0.5)*2,
		VY: -rand.Float64()*2,
		Size: 8 + float64(e.ExpReward), Value: e.ExpReward,
	})

	if rand.Float64() < 0.03+float64(e.EnemyType)*0.05 {
		types := []int{0, 1, 2, 3}
		t := types[rand.Intn(len(types))]
		g.PowerUps = append(g.PowerUps, PowerUp{
			X: e.X, Y: e.Y, Type: t, Size: 22,
		})
	}

	if g.Wave >= 10 && !g.Player.ChainLightning && rand.Float64() < 0.05 {
		g.PowerUps = append(g.PowerUps, PowerUp{X: e.X, Y: e.Y, Type: 4, Size: 24})
	}

	g.createEnemyDeathEffect(e)

	if e.EnemyType == 4 {
		g.BossActive = false
		g.BossDefeated++
		g.Wave++
		g.WaveTimer = 0
		g.WaveDuration += 60 * 5
		g.ScreenShake = 25
		for i := 0; i < 5; i++ {
			types := []int{0, 1, 2, 3}
			t := types[rand.Intn(len(types))]
			g.PowerUps = append(g.PowerUps, PowerUp{
				X: e.X + (rand.Float64()-0.5)*80,
				Y: e.Y + (rand.Float64()-0.5)*40,
				Type: t, Size: 22,
			})
		}
	}

	if e.EnemyType == 1 || e.EnemyType == 3 {
		if rand.Float64() < 0.25 {
			g.SlashWaves = append(g.SlashWaves, SlashWave{
				X: e.X, Y: e.Y,
				Radius: 24,
				Angle: g.Time * 2,
				Spin: 0.18,
				Expand: 2.2 + float64(g.Wave)*0.05,
				Life: 40,
				Damage: 8 + g.Wave/3,
				Color: color.RGBA{255, 120, 220, 255},
			})
		}
	}

	if e.EnemyType == 4 {
		g.SlashWaves = append(g.SlashWaves, SlashWave{
			X: e.X, Y: e.Y,
			Radius: 45,
			Angle: g.Time * 1.5,
			Spin: 0.28,
			Expand: 3.2,
			Life: 50,
			Damage: 14 + g.Wave/2,
			Color: color.RGBA{255, 80, 140, 255},
		})
	}

	g.Enemies = append(g.Enemies[:idx], g.Enemies[idx+1:]...)
}

func (g *Game) damagePlayer(dmg float64) {
	p := g.Player
	if p.Shield > 0 {
		p.Shield = 0
		p.Invincible = 30
		g.createHitParticles(p.X, p.Y, color.RGBA{100, 200, 255, 255}, 20)
		return
	}
	p.Health -= dmg
	p.Invincible = 60
	g.ScreenShake = int(dmg / 2)
	g.createExplosion(p.X, p.Y, 1)
	if p.Health <= 0 {
		g.gameOver()
	}
}

func (g *Game) gameOver() {
	g.State = StateGameOver
	g.saveHighScore()
	g.createExplosion(g.Player.X, g.Player.Y, 3)
	g.ScreenShake = 30
}

func (g *Game) applyPowerUp(t int) {
	p := g.Player
	switch t {
	case 0:
		p.Health = math.Min(p.Health+25, p.MaxHealth)
	case 1:
		p.Shield = 600
	case 2:
		g.Score += 100
	case 3:
		p.Damage += 2
		p.FireRate = max(3, p.FireRate-1)
	case 4:
		p.ChainLightning = true
	}
}

func (g *Game) chainLightningHit(x, y float64, dmg int) {
	if dmg < 1 {
		dmg = 1
	}
	chain := 0
	for i := range g.Enemies {
		e := &g.Enemies[i]
		dx := e.X - x
		dy := e.Y - y
		if math.Sqrt(dx*dx+dy*dy) < 120 {
			e.Health -= dmg
			g.createLightningArc(x, y, e.X, e.Y)
			g.createHitParticles(e.X, e.Y, color.RGBA{120, 220, 255, 255}, 4)
			chain++
			if chain >= 3 {
				break
			}
		}
	}
}

func (g *Game) createLightningArc(x1, y1, x2, y2 float64) {
	segments := 8
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		jx := (rand.Float64() - 0.5) * 10
		jy := (rand.Float64() - 0.5) * 10
		g.Particles = append(g.Particles, Particle{
			X: x1 + (x2-x1)*t + jx,
			Y: y1 + (y2-y1)*t + jy,
			VX: 0, VY: 0,
			Life: 10, MaxLife: 10,
			Color: color.RGBA{120, 220, 255, 255},
			Size: 3,
		})
	}
}

func (g *Game) createExplosion(x, y float64, tier int) {
	count := 12 + tier*8
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 1 + rand.Float64()*(3+float64(tier))
		hue := rand.Float64()*40 + 10
		col := hsvToRgb(hue, 0.8, 1.0)
		g.Particles = append(g.Particles, Particle{
			X: x, Y: y,
			VX: math.Cos(angle) * speed,
			VY: math.Sin(angle) * speed,
			Life: 25 + rand.Intn(20), MaxLife: 45,
			Color: col, Size: 2 + rand.Float64()*3,
		})
	}
	for i := 0; i < count/2; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 0.5 + rand.Float64()*2
		g.Particles = append(g.Particles, Particle{
			X: x, Y: y,
			VX: math.Cos(angle) * speed,
			VY: math.Sin(angle) * speed,
			Life: 15 + rand.Intn(15), MaxLife: 30,
			Color: color.RGBA{255, 255, 220, 255},
			Size:  3 + rand.Float64()*4,
		})
	}
	g.ScreenShake = max(g.ScreenShake, 2+tier*3)
}

func (g *Game) createHitParticles(x, y float64, col color.RGBA, count int) {
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 1 + rand.Float64()*3
		g.Particles = append(g.Particles, Particle{
			X: x, Y: y,
			VX: math.Cos(angle) * speed,
			VY: math.Sin(angle) * speed,
			Life: 12 + rand.Intn(8), MaxLife: 20,
			Color: col, Size: 1 + rand.Float64()*2,
		})
	}
}
