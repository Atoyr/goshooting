package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type CharacterBuilder struct {
	*Character
}

type Character struct {
	*Entity
	characterSize common.CharacterSize
	value         string
}

func (n *Character) ID() uint64 {
	return n.ID()
}

func (n *Character) BasicEntity() ecs.BasicEntity {
	return n.basicEntity
}

func (n *Character) SetEntitySize(width, height float32) {
	n.spaceComponent.Width = width
	n.spaceComponent.Height = height
}

func (n *Character) SetZIndex(index float32) {
	n.SetZIndex(index)
}

func (n *Character) SetVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	n.virtualPosition = &engo.Point{X: point.X, Y: point.Y}
	n.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(*n.virtualPosition))
}

func (n *Character) AddVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	n.virtualPosition.Add(point)
	p := engo.Point{X: n.virtualPosition.X, Y: n.virtualPosition.Y}
	n.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(p))
}

func (n *Character) VirtualPosition() engo.Point {
	return *n.virtualPosition
}

func (n *Character) IsCollision(target Entity) bool {
	return n.IsCollision(target)
}

func (n *Character) SetMoveFunc(movefunc EntityMoveFunc) {
	n.Move = movefunc
}

func (n *Character) SetSpeed(speed float32) {
	n.Speed = speed
}

func (n *Character) SetAngle(angle float32) {
	n.Angle = angle
}

func (n *Character) SetSpeedRate(speedrate float32) {
	n.SpeedRate = speedrate
}

func (n *Character) SetAngleRate(anglerate float32) {
	n.AngleRate = anglerate
}

func (n *Character) AddedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&n.basicEntity, n.renderComponent, n.spaceComponent)
}

func (n *Character) RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64 {
	i := n.ID()
	rs.Remove(n.basicEntity)
	return i
}

func NewCharacterBuilder(size common.CharacterSize) (*CharacterBuilder, error) {
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	r := engoCommon.RenderComponent{Hidden: true}
	em := EntityModel{spaceComponent: &sc, renderComponent: &r}
	emover := EntityMove{}
	e := Entity{EntityModel: &em, EntityMove: &emover}
	c := Character{
		Entity:        &e,
		value:         "",
		characterSize: size,
	}
	return &CharacterBuilder{
		&c,
	}, nil
}

func (nb *CharacterBuilder) SetEntitySize(width, height float32) {
	nb.Entity.SetEntitySize(width, height)
}

func (nb *CharacterBuilder) SetZIndex(index float32) {
	nb.Entity.SetZIndex(index)
}

func (nb *CharacterBuilder) SetVirtualPosition(point engo.Point) {
	nb.Entity.SetVirtualPosition(point)
}

func (nb *CharacterBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
}

func (nb *CharacterBuilder) SetCollisionDetectionSize(size float32) {
}

func (nb *CharacterBuilder) SetMoveFunc(movefunc EntityMoveFunc) {
	nb.Entity.Move = movefunc
}

func (nb *CharacterBuilder) SetSpeed(speed float32) {
	nb.Entity.Speed = speed
}

func (nb *CharacterBuilder) SetAngle(angle float32) {
	nb.Entity.Angle = angle
}

func (nb *CharacterBuilder) SetSpeedRate(speedrate float32) {
	nb.Entity.SpeedRate = speedrate
}

func (nb *CharacterBuilder) SetAngleRate(anglerate float32) {
	nb.Entity.AngleRate = anglerate
}

func (nb *CharacterBuilder) Build() Character {
	n := *nb.Character
	n.basicEntity = ecs.NewBasic()

	return n
}

func (n *Character) SetCharacter(value string) error {
	n.value = value
	t, err := common.GetCharacterTexture(value, n.characterSize)
	if err != nil {
		return err
	}
	n.SetDrawable(t)
	n.renderComponent.Hidden = false
	return nil
}
