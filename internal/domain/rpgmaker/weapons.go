package rpgmaker

// Weapon represents a weapon item
type Weapon struct {
	ID          int     `json:"id"`
	AnimationId int     `json:"animationId"`
	Description string  `json:"description"`
	EtypeId     int     `json:"etypeId"`
	Traits      []Trait `json:"traits"`
	IconIndex   int     `json:"iconIndex"`
	Name        string  `json:"name"`
	Note        string  `json:"note"`
	Params      []int   `json:"params"`
	Price       int     `json:"price"`
	WtypeId     int     `json:"wtypeId"`
}

// WeaponsData is an array of weapons
type WeaponsData []*Weapon




