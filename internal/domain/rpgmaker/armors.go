package rpgmaker

// Armor represents an armor item
type Armor struct {
	ID          int     `json:"id"`
	AtypeID     int     `json:"atypeId"`
	Description string  `json:"description"`
	EtypeID     int     `json:"etypeId"`
	Traits      []Trait `json:"traits"`
	IconIndex   int     `json:"iconIndex"`
	Name        string  `json:"name"`
	Note        string  `json:"note"`
	Params      []int   `json:"params"`
	Price       int     `json:"price"`
}

// ArmorsData is an array of armors
type ArmorsData []*Armor




