package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// Weapon represents a weapon item
type Weapon struct {
	ID          int                        `json:"id"`
	AnimationId int                        `json:"animationId"`
	Description string                     `json:"description"`
	EtypeId     int                        `json:"etypeId"`
	Traits      []Trait                    `json:"traits"`
	IconIndex   int                        `json:"iconIndex"`
	Name        string                     `json:"name"`
	Note        string                     `json:"note"`
	Params      []int                      `json:"params"`
	Price       int                        `json:"price"`
	WtypeId     int                        `json:"wtypeId"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (w *Weapon) UnmarshalJSON(data []byte) error {
	type Alias Weapon
	aux := (*Alias)(w)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	w.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(w))
	return nil
}

func (w Weapon) MarshalJSON() ([]byte, error) {
	type Alias Weapon
	return util.MarshalWithExtras((*Alias)(&w), w.Extras)
}

// WeaponsData is an array of weapons
type WeaponsData []*Weapon
