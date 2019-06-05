package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type NumberBuilder struct {
	*Number
}

type Number struct {
	*Entity
	texture []engoCommon.Texture
	value   int
}

func (n *Number) ID() uint64 {
	return n.ID()
}

func (n *Number) BasicEntity() ecs.BasicEntity {
	return n.basicEntity
}

func (n *Number) SetEntitySize(width, height float32) {
	n.spaceComponent.Width = width
	n.spaceComponent.Height = height
}

func (n *Number) SetZIndex(index float32) {
	n.SetZIndex(index)
}

func (n *Number) SetVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	n.virtualPosition = engo.Point{X: point.X, Y: point.Y}
	n.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(n.virtualPosition))
}

func (n *Number) AddVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	n.virtualPosition.Add(point)
	p := engo.Point{X: n.virtualPosition.X, Y: n.virtualPosition.Y}
	n.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(p))
}

func (n *Number) VirtualPosition() engo.Point {
	return n.virtualPosition
}

func (n *Number) IsCollision(target Entity) bool {
	return n.IsCollision(target)
}

func (n *Number) SetSpeed(speed float32) {
	n.Speed = speed
}

func (n *Number) SetAngle(angle float32) {
	n.Angle = angle
}

func (n *Number) SetSpeedRate(speedrate float32) {
	n.SpeedRate = speedrate
}

func (n *Number) SetAngleRate(anglerate float32) {
	n.AngleRate = anglerate
}

func (n *Number) AddedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&n.basicEntity, &n.renderComponent, &n.spaceComponent)
}

func (n *Number) RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64 {
	i := n.ID()
	rs.Remove(n.basicEntity)
	return i
}

func NewNumberBuilder(size common.NumberSize) (*NumberBuilder, error) {
	t, err := common.GetNumberTextures(size)
	if err != nil {
		return nil, err
	}

	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	r := engoCommon.RenderComponent{Drawable: t[0], Hidden: true}
	em := EntityModel{spaceComponent: sc, renderComponent: r}
	emover := EntityMove{}
	e := Entity{EntityModel: &em, EntityMove: &emover}
	n := Number{
		Entity:  &e,
		value:   -1,
		texture: t,
	}
	return &NumberBuilder{
		&n,
	}, nil
}

func (nb *NumberBuilder) SetEntitySize(width, height float32) {
	nb.Entity.SetEntitySize(width, height)
}

func (nb *NumberBuilder) SetZIndex(index float32) {
	nb.Entity.SetZIndex(index)
}

func (nb *NumberBuilder) SetVirtualPosition(point engo.Point) {
	nb.Entity.SetVirtualPosition(point)
}

func (nb *NumberBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
}

func (nb *NumberBuilder) SetCollisionDetectionSize(size float32) {
}

func (nb *NumberBuilder) SetSpeed(speed float32) {
	nb.Entity.Speed = speed
}

func (nb *NumberBuilder) SetAngle(angle float32) {
	nb.Entity.Angle = angle
}

func (nb *NumberBuilder) SetSpeedRate(speedrate float32) {
	nb.Entity.SpeedRate = speedrate
}

func (nb *NumberBuilder) SetAngleRate(anglerate float32) {
	nb.Entity.AngleRate = anglerate
}

func (nb *NumberBuilder) Build() Number {
	n := *nb.Number
	n.basicEntity = ecs.NewBasic()

	return n
}

func (n *Number) SetNumber(value int) {
	if value < 0 {
		n.renderComponent.Hidden = true
	} else {
		num := value % 10
		n.value = num
		n.SetDrawable(n.texture[num])
		n.renderComponent.Hidden = false
	}
}

func (n *Number) Add(value int) {
	if n.value == -1 {
		return
	}

	num := (n.value + value) % 10
	n.SetNumber(num)
}
