package shooter

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miyahoyo/godanmaku/danmaku/internal/collision"
	"github.com/miyahoyo/godanmaku/danmaku/internal/effect"

	"github.com/miyahoyo/godanmaku/danmaku/internal/flyweight"

	"github.com/miyahoyo/godanmaku/danmaku/internal/field"
	"github.com/miyahoyo/godanmaku/danmaku/internal/weapon"

	"github.com/miyahoyo/godanmaku/danmaku/internal/sprite"
	"github.com/miyahoyo/godanmaku/danmaku/internal/util"
)

func init() {
	healthBar = NewHealthBar()
}

var healthBar *HealthBar

type Target interface {
	GetX() float64
	GetY() float64
	IsDead() bool
}

type Shooter struct {
	x, y          float64
	width, height float64
	field         *field.Field
	isActive      bool
	speed         float64
	vx            float64
	vy            float64
	degree        int
	name          string
	spr           *sprite.Sprite
	sprIndex      int
	life          int
	maxLife       int
	updateCount   int
	mainWeapon    weapon.Weapon
	target        Target
	destination   struct{ x, y float64 }
	shotsPool     *flyweight.Pool
	collisionBox  []*collision.Box
	controller    Controller
	funnel        []*Shooter
	owner         *Shooter
	funnelDegree  int
}

func NewShooter() *Shooter {
	sh := &Shooter{}

	return sh
}

func (sh *Shooter) init() {
	sh.controller.init(sh)
}

func (sh *Shooter) Update() {
	sh.updateCount++
	sh.controller.update(sh)
}

func (sh *Shooter) UpdatePlayer() {
	sh.updateCount++
	sh.controller.update(sh)
}

func (sh *Shooter) Draw(screen *ebiten.Image) {
	sh.controller.draw(sh, screen)
}

func (sh *Shooter) GetX() float64 {
	return sh.x
}

func (sh *Shooter) GetY() float64 {
	return sh.y
}

func (sh *Shooter) GetPosition() (float64, float64) {
	return sh.x, sh.y
}

func (sh *Shooter) GetWidth() float64 {
	return sh.width
}

func (sh *Shooter) GetHeight() float64 {
	return sh.height
}

func (sh *Shooter) GetName() string {
	return sh.name
}

// Implement quadtree.Collider interface
func (sh *Shooter) GetRect() (x0 float64, y0 float64, x1 float64, y1 float64) {
	return sh.x - sh.width/2, sh.y - sh.height/2, sh.x + sh.width/2, sh.y + sh.height/2
}

func (sh *Shooter) GetCollisionBox() []*collision.Box {
	return sh.collisionBox
}

func (sh *Shooter) GetDegree() int {
	return sh.degree
}

func (sh *Shooter) IsActive() bool {
	return sh.isActive
}

func (sh *Shooter) GetMainSpriteIndex() int {
	return sh.sprIndex
}

func (sh *Shooter) SetMainSpriteIndex(index int) {
	sh.sprIndex = index
}

func (sh *Shooter) SetSpeed(speed float64, degree int) {
	sh.speed = speed
	sh.degree = degree
	sh.vx = math.Cos(util.DegToRad(sh.degree)) * speed
	sh.vy = math.Sin(util.DegToRad(sh.degree)) * speed
}

func (sh *Shooter) AddDamage(damage int) {
	sh.life -= damage
	if sh.life <= 0 {
		sh.isActive = false
		effect.CreateExplosion(sh.x, sh.y)
	}
}

func (sh *Shooter) IsDead() bool {
	return sh.life <= 0
}

func (sh *Shooter) GetLife() int {
	return sh.life
}

func (sh *Shooter) SetWeapon(mainWeapon weapon.Weapon) {
	sh.mainWeapon = mainWeapon
}

func (sh *Shooter) SetPosition(x, y float64) {
	sh.x = x
	sh.y = y
}

func (sh *Shooter) Recovery() {
	sh.life = sh.maxLife
}

func (sh *Shooter) Fire() bool {
	var ret = false
	ret = sh.mainWeapon.Fire(sh, sh.x, sh.y, sh.degree)
	if len(sh.funnel) > 0 {
		for f := range sh.funnel {
			if sh.funnel[f].Fire() {
				ret = true
			}
		}
	}
	return ret
}

func (sh *Shooter) SetTarget(target Target) {
	sh.target = target
}

func (sh *Shooter) setSize(width, height float64) {
	sh.width = width
	sh.height = height
}
