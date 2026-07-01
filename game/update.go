package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	metaHPCost      = 180
	metaBulletCost  = 320
	metaBoostCost   = 520
	metaCardBaseCost = 260
)

func waveGrowth(wave int) float64 {
	if wave <= 1 {
		return 1
	}
	f := float64(wave - 1)
	return 1 + 0.18*f + 0.018*f*f + 0.0012*f*f*f
}

func expThreshold(level int) int {
	if level <= 1 {
		return 6
	}
	f := float64(level - 1)
	value := 6 + 4*f + 1.6*f*f + 0.45*f*f*f
	if value > 240 {
		value = 240
	}
	return int(value)
}

// 第5波后进入爆发膨胀区间，每两波整体属性翻倍。
func enemyGrowthMultiplier(wave int) float64 {
	if wave < 5 {
		return 1
	}
	steps := (wave-5)/2 + 1
	return math.Pow(2, float64(steps))
}

func NewGame() *Game {
	g := &Game{
		State: StateMenu,
		Player: &Player{
			X: ScreenWidth / 2, Y: ScreenHeight - 100,
			Width: 36, Height: 36,
			Health: 100, MaxHealth: 100,
			Damage: 6, FireRate: 15,
			BulletSpeed: 10,
			BulletCount: 1, Spread: 0,
			SideGuns: 0, BackGun: false, LaserLevel: 0,
		MoveSpeed: 1.0, PickupRange: 60,
		CritChance: 0.05, CritMult: 2.0,
		Regen: 0,
		Level: 1, Exp: 0, ExpToNext: expThreshold(1),
	},
		Score:            0,
		Wave:             1,
		WaveDuration:     60 * 28,
		EnemySpawnTimer:  0,
		BattlefieldScale: 1,
	}
	g.initStars()
	g.loadHighScore()
	g.loadMeta()
	return g
}

func (g *Game) initStars() {
	g.Stars = make([]Star, 120)
	for i := range g.Stars {
		g.Stars[i] = Star{
			X:     rand.Float64() * ScreenWidth,
			Y:     rand.Float64() * ScreenHeight,
			Speed: 0.3 + rand.Float64()*1.5,
			Size:  1 + rand.Float64()*2,
			Alpha: 0.3 + float32(rand.Float64()*0.7),
		}
	}
}

func (g *Game) Update() error {
	g.Time += 1.0 / 60.0
	if g.ScreenShake > 0 {
		g.ScreenShake--
	}
	switch g.State {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			g.State = StatePlaying
		}
	case StateLevelUp:
		g.updateLevelUp()
	case StateGameOver:
		g.updateGameOver()
	}
	return nil
}

func (g *Game) updateMenu() {
	g.updateCheatCode()
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		g.buyMetaHP()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		g.buyMetaBullet()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		g.buyMetaBoost()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit4) {
		g.buyMetaInitialCard()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.startGame()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startGame()
	}
}

func (g *Game) spendMetaCoins(cost int) bool {
	if g.Meta.InfiniteGold {
		return true
	}
	if g.Meta.Coins < cost {
		return false
	}
	g.Meta.Coins -= cost
	return true
}

func (g *Game) buyMetaHP() {
	if g.spendMetaCoins(metaHPCost) {
		g.Meta.HPLevel++
		g.saveMeta()
	}
}

func (g *Game) buyMetaBullet() {
	if g.Meta.BulletLevel >= 2 {
		return
	}
	if g.spendMetaCoins(metaBulletCost) {
		g.Meta.BulletLevel++
		g.saveMeta()
	}
}

func (g *Game) buyMetaBoost() {
	if g.Meta.BoostCard {
		return
	}
	if g.spendMetaCoins(metaBoostCost) {
		g.Meta.BoostCard = true
		g.saveMeta()
	}
}

func (g *Game) buyMetaInitialCard() {
	if g.Meta.InitialCards >= 3 {
		return
	}
	cost := metaCardBaseCost * (g.Meta.InitialCards + 1)
	if g.spendMetaCoins(cost) {
		g.Meta.InitialCards++
		g.saveMeta()
	}
}

func (g *Game) updateCheatCode() {
	seq := []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyB, ebiten.KeyA}
	for _, key := range []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyB, ebiten.KeyA} {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}
		if key == seq[g.CheatIndex] {
			g.CheatIndex++
			if g.CheatIndex >= len(seq) {
				g.Meta.InfiniteGold = true
				g.Meta.Coins = 999999
				g.saveMeta()
				g.CheatIndex = 0
			}
		} else {
			g.CheatIndex = 0
		}
	}
}

func (g *Game) startGame() {
	g.State = StatePlaying
	p := g.Player
	p.X = ScreenWidth / 2
	p.Y = ScreenHeight - 100
	p.Health = 100
	p.MaxHealth = 100 + float64(g.Meta.HPLevel*25)
	p.Health = p.MaxHealth
	p.Damage = 6
	p.FireRate = 15
	p.FireCooldown = 0
	p.BulletSpeed = 10
	p.BulletCount = min(3, 1+g.Meta.BulletLevel)
	p.Spread = 0
	p.SideGuns = 0
	p.FrontArc = false
	p.LaserBeam = false
	p.ChainLightning = false
	p.TwinShip = false
	p.BackGun = false
	p.LaserLevel = 0
	p.MoveSpeed = 1.0
	p.PickupRange = 60
	p.CritChance = 0.05
	p.CritMult = 2.0
	p.Regen = 0
	p.Shield = 0
	p.Invincible = 60
	p.Level = 1
	p.Exp = 0
	p.ExpToNext = expThreshold(1)
	p.Kills = 0
	p.Coins = 0
	g.Bullets = nil
	g.Enemies = nil
	g.PowerUps = nil
	g.Coins = nil
	g.ExpOrbs = nil
	g.Particles = nil
	g.Score = 0
	g.Wave = 1
	if g.Meta.BoostCard {
		g.Wave = 5
		p.Damage += 18
		p.FireRate = max(4, p.FireRate-5)
		p.BulletSpeed *= 1.35
		p.MaxHealth += 80
		p.Health = p.MaxHealth
	}
	g.WaveTimer = 0
	g.EnemySpawnTimer = 0
	g.BossActive = false
	g.BossDefeated = 0
	g.Upgrades = nil
	g.BattlefieldScale = 1
	g.PendingInitialCards = g.Meta.InitialCards
	if g.PendingInitialCards > 0 {
		g.rollUpgrades()
		g.State = StateLevelUp
	}
}

func (g *Game) updatePlaying() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.State = StatePaused
		return
	}
	g.WaveTimer++
	g.updateStars()
	g.updatePlayer()
	g.updateBullets()
	g.updateEnemies()
	g.updateCoins()
	g.updateExpOrbs()
	g.updateSlashWaves()
	g.updatePowerUps()
	g.updateParticles()
	g.spawnEnemies()
	g.checkCollisions()
	g.updateWave()
}

func (g *Game) updateLevelUp() {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		btnW := 180
		gap := 24
		totalW := btnW*3 + gap*2
		startX := ScreenWidth/2 - totalW/2
		for i, u := range g.UpgradeChoices {
			btnX := startX + i*(btnW+gap)
			btnY := ScreenHeight/2 - 90
			if mx >= btnX && mx <= btnX+btnW && my >= btnY && my <= btnY+180 {
				u.Apply(g.Player)
				g.Upgrades = append(g.Upgrades, u)
				g.Player.ExpToNext = expThreshold(g.Player.Level)
				if g.PendingInitialCards > 0 {
					g.PendingInitialCards--
					if g.PendingInitialCards > 0 {
						g.rollUpgrades()
						return
					}
				}
				if g.State == StateLevelUp {
					g.State = StatePlaying
				}
				break
			}
		}
	}
}

func (g *Game) updateGameOver() {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.startGame()
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mx >= ScreenWidth/2-80 && mx <= ScreenWidth/2+80 &&
			my >= ScreenHeight/2+40 && my <= ScreenHeight/2+80 {
			g.startGame()
		}
		if mx >= ScreenWidth/2-80 && mx <= ScreenWidth/2+80 &&
			my >= ScreenHeight/2+100 && my <= ScreenHeight/2+140 {
			g.State = StateMenu
		}
	}
}

func (g *Game) updateStars() {
	for i := range g.Stars {
		s := &g.Stars[i]
		s.Y += s.Speed
		if s.Y > ScreenHeight {
			s.Y = 0
			s.X = rand.Float64() * ScreenWidth
		}
	}
}

func (g *Game) updatePlayer() {
	p := g.Player
	mx, my := ebiten.CursorPosition()
	targetX := float64(mx)
	targetY := float64(my)

	smooth := 0.25 * p.MoveSpeed
	if smooth > 0.9 {
		smooth = 0.9
	}
	p.X += (targetX - p.X) * smooth
	p.Y += (targetY - p.Y) * smooth
	p.X = math.Max(p.Width/2, math.Min(ScreenWidth-p.Width/2, p.X))
	p.Y = math.Max(p.Height/2, math.Min(ScreenHeight-p.Height/2, p.Y))

	if p.Invincible > 0 {
		p.Invincible--
	}
	if p.Shield > 0 {
		p.Shield--
	}

	if p.Regen > 0 {
		p.RegenTimer++
		if p.RegenTimer >= 60 {
			p.RegenTimer = 0
			if p.Health < p.MaxHealth {
				p.Health = math.Min(p.Health+p.Regen, p.MaxHealth)
			}
		}
	}

	p.FireCooldown--
	if p.FireCooldown <= 0 {
		p.FireCooldown = p.FireRate
		g.fireWeapons()
	}
}
