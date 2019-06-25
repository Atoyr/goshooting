package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
)

type Enemy struct {
	*EntityModel
}

type EnemyBuilder struct {
	*Enemy
}

func NewEnemyBuilder() EnemyBuilder {
	s := common.NewSetting()
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	rc := engoCommon.RenderComponent{Scale: s.Scale()}
	model := EntityModel{
		spaceComponent:  sc,
		renderComponent: rc,
		virtualPosition: engo.Point{X: 0, Y: 0},
		scale:           0.5,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)
	model.SetPosition(engo.Point{X: 0, Y: 0})

	enemy := new(Enemy)
	enemy.EntityModel = &model

	return EnemyBuilder{enemy}
}

func (eb *EnemyBuilder) Build() Modeler {
	entityModel := new(EntityModel)
	enemy := new(Enemy)
	copier.Copy(&entityModel, eb.EntityModel)
	copier.Copy(&enemy, eb.Enemy)
	enemy.EntityModel = entityModel
	enemy.basicEntity = ecs.NewBasic()

	return *enemy
}
