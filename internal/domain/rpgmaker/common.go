package rpgmaker

// Trait represents a character/item trait in RPG Maker
type Trait struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value  float64 `json:"value"`
}

// EventCommand represents a command in an event or troop page
type EventCommand struct {
	Code       int   `json:"code"`
	Indent     int   `json:"indent"`
	Parameters []any `json:"parameters"`
}




