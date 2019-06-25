package entitys

// Builder is Entity Build Interface
type Builder interface {
	Build() Modeler
	Clone() Builder
}

// Mover is Entity Move Interface
type Mover interface {
	Move(frame uint64)
}

// Attacker is Entity Attacking interface
type Attacker interface {
	Attack(modeler Modeler, frame uint64) []Modeler
}

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func()
