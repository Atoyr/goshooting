package entitys

type CollisionBuilder struct {
	*EntityModel
}

func NewCollisionBuilder() CollisionBuilder {
	s := common.NewSetting()
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	rc := engoCommon.RenderComponent{Scale: s.Scale()}
	model := EntityModel{
		spaceComponent:  sc,
		renderComponent: rc,
		virtualPosition: engo.Point{X: 0, Y: 0},
		scale:           0.5,
		hitPoint:        0,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)

	model.SetPosition(engo.Point{X: 0, Y: 0})
	return CollisionBuilder{model}
}

func (cb *CollisionBuilder) Build() Modeler {
	entityModel := new(EntityModel)
	copier.Copy(&entityModel, bb.EntityModel)
	entityModel.basicEntity = ecs.NewBasic()

	return *entityModel
}

func (cb *CollisionBuilder) Clone() Builder {
	builder := new(CollisionBuilder)
	entityModel := new(EntityModel)
	copier.Copy(&entityModel, bb.EntityModel)

	return builder
}
