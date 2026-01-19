package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// Actor represents a playable character
type Actor struct {
	ID             int                        `json:"id"`
	BattlerName    string                     `json:"battlerName"`
	CharacterIndex int                        `json:"characterIndex"`
	CharacterName  string                     `json:"characterName"`
	ClassId        int                        `json:"classId"`
	Equips         []int                      `json:"equips"`
	FaceIndex      int                        `json:"faceIndex"`
	FaceName       string                     `json:"faceName"`
	Traits         []Trait                    `json:"traits"`
	InitialLevel   int                        `json:"initialLevel"`
	MaxLevel       int                        `json:"maxLevel"`
	Name           string                     `json:"name"`
	Nickname       string                     `json:"nickname"`
	Note           string                     `json:"note"`
	Profile        string                     `json:"profile"`
	Extras         map[string]json.RawMessage `json:"-"`
}

func (a *Actor) UnmarshalJSON(data []byte) error {
	type Alias Actor
	aux := (*Alias)(a)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(a))
	return nil
}

func (a Actor) MarshalJSON() ([]byte, error) {
	type Alias Actor
	return util.MarshalWithExtras((*Alias)(&a), a.Extras)
}

// ActorsData is an array of actors
type ActorsData []*Actor
