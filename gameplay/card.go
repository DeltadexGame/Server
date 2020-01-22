package gameplay

// Card represents a physical card held by a player
type Card struct {
	Name       string
	Type       int
	Attack     int
	Health     int
	EnergyCost int
	Ability    Ability
}

// Ability keeps track of the ability/effect of the card
type Ability struct {
	AbilityID   int
	Name        string
	Description string
	Targeted    bool
}
