package entitys

type EntityContainer struct {
	IsLeft bool
	nodes  []entityNode
}

type entityNode struct {
	entity *Entity
	nodes  []entityNode
}
