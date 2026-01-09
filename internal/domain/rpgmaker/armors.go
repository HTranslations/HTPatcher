package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// Armor represents an armor item
type Armor struct {
	ID          int                        `json:"id"`
	AtypeID     int                        `json:"atypeId"`
	Description string                     `json:"description"`
	EtypeID     int                        `json:"etypeId"`
	Traits      []Trait                    `json:"traits"`
	IconIndex   int                        `json:"iconIndex"`
	Name        string                     `json:"name"`
	Note        string                     `json:"note"`
	Params      []int                      `json:"params"`
	Price       int                        `json:"price"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (a *Armor) UnmarshalJSON(data []byte) error {
	type Alias Armor
	aux := (*Alias)(a)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(a))
	return nil
}

func (a Armor) MarshalJSON() ([]byte, error) {
	type Alias Armor
	return util.MarshalWithExtras((*Alias)(&a), a.Extras)
}

// ArmorsData is an array of armors
type ArmorsData []*Armor
