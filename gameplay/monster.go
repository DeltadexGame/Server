package gameplay

// Monster keeps track of monsters currently place on the board
type Monster struct {
	Name      string
	Attack    int
	Health    int
	MaxHealth int
	Ability   Ability
}

// Damage decreases a monster's health by the given amount
func (m *Monster) Damage(damage int) {
	m.Health -= damage
}
