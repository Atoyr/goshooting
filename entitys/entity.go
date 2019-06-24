package entitys

// Builder is Entity Build Interface
type Builder interface {
	Build() Modeler
}

// Mover is Entity Move Interface
type Mover interface {
	Move(vx, vy, speed float32)
}

// Attacker is Entity Attacking interface
type Attacker interface {
	Attack()
}

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func()
