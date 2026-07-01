package game

import (
	"math/rand"
	"sort"
)

var upgradePool = []Upgrade{
	{0, "攻击强化", "伤害 +25%", 1, UpgradeTagStat, func(p *Player) { p.Damage = int(float64(p.Damage) * 1.25) }},
	{1, "快速射击", "射速 +20%", 1, UpgradeTagRate, func(p *Player) { p.FireRate = max(3, int(float64(p.FireRate)*0.80)) }},
	{2, "弹速提升", "子弹速度 +25%", 1, UpgradeTagStat, func(p *Player) { p.BulletSpeed *= 1.25 }},
	{3, "多重射击", "前方 +1 发子弹", 2, UpgradeTagCount, func(p *Player) { p.BulletCount++ }},
	{4, "扩散射击", "三向散射更宽", 2, UpgradeTagWeapon, func(p *Player) { p.Spread += 10 }},
	{5, "侧翼炮", "左右各 +1 发子弹", 2, UpgradeTagCount, func(p *Player) { p.SideGuns++ }},
	{6, "双重战机", "生成双机协同射击", 3, UpgradeTagUltimate, func(p *Player) { p.TwinShip = true }},
	{7, "激光穿透", "高伤害长射线", 3, UpgradeTagExclusive, func(p *Player) { p.LaserBeam = true }},
	{8, "生命强化", "最大生命 +25", 1, UpgradeTagStat, func(p *Player) { p.MaxHealth += 25; p.Health += 25 }},
	{9, "紧急维修", "恢复 40 生命", 1, UpgradeTagStat, func(p *Player) { if p.Health < p.MaxHealth { p.Health = min(p.Health+40, p.MaxHealth) } }},
	{10, "推进器", "移动速度 +18%", 1, UpgradeTagStat, func(p *Player) { p.MoveSpeed *= 1.18 }},
	{11, "磁力装置", "拾取范围 +60%", 1, UpgradeTagStat, func(p *Player) { p.PickupRange *= 1.6 }},
	{12, "致命一击", "暴击率 +12%", 2, UpgradeTagStat, func(p *Player) { p.CritChance += 0.12 }},
	{13, "暴击伤害", "暴击伤害 +60%", 2, UpgradeTagStat, func(p *Player) { p.CritMult += 0.6 }},
	{14, "纳米修复", "每秒恢复 0.8 生命", 2, UpgradeTagStat, func(p *Player) { p.Regen += 0.8 }},
	{15, "巨型弹头", "子弹伤害 +2", 2, UpgradeTagStat, func(p *Player) { p.Damage += 2 }},
	{16, "装甲镀层", "最大生命 +40", 2, UpgradeTagStat, func(p *Player) { p.MaxHealth += 40; p.Health += 40 }},
	{17, "连锁闪电", "稀有连锁伤害", 3, UpgradeTagUltimate, func(p *Player) { p.ChainLightning = true }},
	{18, "前向激光", "双战机可用的长射线", 3, UpgradeTagExclusive, func(p *Player) { p.FrontArc = true }},
}

func (g *Game) rollUpgrades() {
	shuffled := make([]Upgrade, len(upgradePool))
	copy(shuffled, upgradePool)
	filtered := shuffled[:0]
	for _, u := range shuffled {
		if u.Tag == UpgradeTagExclusive && g.hasExclusiveUpgrade() {
			continue
		}
		if u.ID == 7 && g.Player.FrontArc {
			continue
		}
		filtered = append(filtered, u)
	}
	shuffled = filtered
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	count := 3
	if count > len(shuffled) {
		count = len(shuffled)
	}
	g.UpgradeChoices = weightedPickUpgrades(g.Player.Level, shuffled, count)
	sort.Slice(g.UpgradeChoices, func(i, j int) bool {
		return g.UpgradeChoices[i].Rarity < g.UpgradeChoices[j].Rarity
	})
}

func (g *Game) hasExclusiveUpgrade() bool {
	for _, u := range g.Upgrades {
		if u.Tag == UpgradeTagExclusive || u.ID == 6 || u.ID == 7 || u.ID == 17 || u.ID == 18 {
			return true
		}
	}
	return false
}

func weightedPickUpgrades(level int, pool []Upgrade, count int) []Upgrade {
	if len(pool) == 0 {
		return nil
	}
	result := make([]Upgrade, 0, count)
	working := make([]Upgrade, len(pool))
	copy(working, pool)
	for len(result) < count && len(working) > 0 {
		total := 0
		weights := make([]int, len(working))
		for i, u := range working {
			w := 1
			switch level {
			case 1, 2:
				if u.Tag == UpgradeTagRate || u.Tag == UpgradeTagCount {
					w = 8
				} else if u.Rarity == 3 {
					w = 1
				} else {
					w = 4
				}
			case 3, 4:
				if u.Tag == UpgradeTagRate || u.Tag == UpgradeTagCount {
					w = 6
				} else if u.Rarity == 3 {
					w = 2
				} else {
					w = 4
				}
			default:
				if u.Rarity == 3 {
					w = 1
				} else {
					w = 3
				}
			}
			if u.Tag == UpgradeTagExclusive && level < 6 {
				w = 1
			}
			weights[i] = w
			total += w
		}
		r := rand.Intn(total)
		pick := 0
		for i, w := range weights {
			r -= w
			if r < 0 {
				pick = i
				break
			}
		}
		result = append(result, working[pick])
		working = append(working[:pick], working[pick+1:]...)
	}
	return result
}

func (g *Game) addExp(amount int) {
	p := g.Player
	p.Exp += amount
	for p.Exp >= p.ExpToNext {
		p.Exp -= p.ExpToNext
		p.Level++
		p.ExpToNext = expThreshold(p.Level)
		g.rollUpgrades()
		g.State = StateLevelUp
	}
}
