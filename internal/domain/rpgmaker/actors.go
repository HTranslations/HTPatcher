package rpgmaker

// Actor represents a playable character
type Actor struct {
	ID             int     `json:"id"`
	BattlerName    string  `json:"battlerName"`
	CharacterIndex int     `json:"characterIndex"`
	CharacterName  string  `json:"characterName"`
	ClassId        int     `json:"classId"`
	Equips         []int   `json:"equips"`
	FaceIndex      int     `json:"faceIndex"`
	FaceName       string  `json:"faceName"`
	Traits         []Trait `json:"traits"`
	InitialLevel   int     `json:"initialLevel"`
	MaxLevel       int     `json:"maxLevel"`
	Name           string  `json:"name"`
	Nickname       string  `json:"nickname"`
	Note           string  `json:"note"`
	Profile        string  `json:"profile"`
}

// ActorsData is an array of actors
type ActorsData []*Actor




