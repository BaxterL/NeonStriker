package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 800
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateLevelUp
	StateGameOver
)

type Player struct {
	X, Y          float64
	Width, Height float64
	Health        float64
	MaxHealth     float64
	Damage        int
	FireRate      int
	FireCooldown  int
	BulletSpeed   float64
	BulletCount   int
	Spread        int
	SideGuns      int
	FrontArc      bool
	LaserBeam     bool
	ChainLightning bool
	TwinShip      bool
	BackGun       bool
	LaserLevel    int
	MoveSpeed     float64
	PickupRange   float64
	CritChance    float64
	CritMult      float64
	Regen         float64
	RegenTimer    int
	Shield        int
	ShieldSkill   bool
	ShieldTimer   int
	Lifesteal     float64
	CritDefense   bool
	GiantSlayer   bool
	FlyingKick    bool
	CritRhythm    bool
	CritRhythmStacks int
	CritRhythmTimer  int
	UltimateMark  bool
	UltimateTimer int
	UltimateWindow int
	BaronHand     bool
	Invincible    int
	Level         int
	Exp           int
	ExpToNext     int
	Kills         int
	Coins         int
}

type Bullet struct {
	X, Y   float64
	VX, VY float64
	Damage int
	Player bool
	Size   float64
	Life   int
	Color  color.RGBA
	Laser  bool
	Beam   bool
	Pierce int
}

type Enemy struct {
	X, Y       float64
	Width      float64
	Height     float64
	Health     int
	MaxHealth  int
	VX, VY     float64
	Score      int
	ExpReward  int
	CoinDrop   int
	EnemyType  int
	FireCd     int
	FireTimer  int
	BulletType int
	Hue        float64
	MarkedDamage int
}

type PowerUp struct {
	X, Y float64
	Type int
	Size float64
}

type Coin struct {
	X, Y   float64
	VX, VY float64
	Size   float64
	Value  int
}

type ExpOrb struct {
	X, Y   float64
	VX, VY float64
	Size   float64
	Value  int
}

type SlashWave struct {
	X, Y     float64
	Radius   float64
	Angle    float64
	Spin     float64
	Expand   float64
	Life     int
	Damage   int
	Color    color.RGBA
	HitOnce  bool
}

type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
	Color   color.RGBA
	Size    float64
	Type    int
}

type Star struct {
	X, Y  float64
	Speed float64
	Size  float64
	Alpha float32
}

type Upgrade struct {
	ID     int
	Name   string
	Desc   string
	Rarity int
	Tag    int
	Apply  func(*Player)
}

type MetaProgress struct {
	Coins        int  `json:"coins"`
	HPLevel      int  `json:"hp_level"`
	BulletLevel  int  `json:"bullet_level"`
	BoostCard    bool `json:"boost_card"`
	InitialCards int  `json:"initial_cards"`
	InfiniteGold bool `json:"infinite_gold"`
}

const (
	UpgradeTagStat = iota
	UpgradeTagRate
	UpgradeTagCount
	UpgradeTagWeapon
	UpgradeTagUltimate
	UpgradeTagExclusive
)

type Game struct {
	State            GameState
	Player           *Player
	Bullets          []Bullet
	Enemies          []Enemy
	PowerUps         []PowerUp
	Coins            []Coin
	ExpOrbs          []ExpOrb
	SlashWaves       []SlashWave
	Particles        []Particle
	Stars            []Star
	Score            int
	HighScore        int
	Wave             int
	WaveTimer        int
	WaveDuration     int
	EnemySpawnTimer  int
	BossActive       bool
	BossDefeated     int
	BattlefieldScale  float64
	Time             float64
	ScreenShake      int
	offscreen        *ebiten.Image
	UpgradeChoices   []Upgrade
	Upgrades         []Upgrade
	Meta             MetaProgress
	CheatIndex       int
	PendingInitialCards int
}

func colorScale(c color.Color) ebiten.ColorScale {
	r, g, b, a := c.RGBA()
	var cs ebiten.ColorScale
	cs.Scale(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	return cs
}

func hsvToRgb(h, s, v float64) color.RGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return color.RGBA{
		uint8((r + m) * 255),
		uint8((g + m) * 255),
		uint8((b + m) * 255),
		255,
	}
}

func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (g *Game) loadHighScore() {
	data, err := os.ReadFile("highscore.txt")
	if err != nil {
		g.HighScore = 0
		return
	}
	fmt.Sscanf(string(data), "%d", &g.HighScore)
}

func (g *Game) saveHighScore() {
	if g.Score > g.HighScore {
		g.HighScore = g.Score
		os.WriteFile("highscore.txt", []byte(fmt.Sprintf("%d", g.HighScore)), 0644)
	}
}

func (g *Game) loadMeta() {
	data, err := os.ReadFile("meta_progress.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &g.Meta)
}

func (g *Game) saveMeta() {
	data, err := json.MarshalIndent(g.Meta, "", "  ")
	if err == nil {
		os.WriteFile("meta_progress.json", data, 0644)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
