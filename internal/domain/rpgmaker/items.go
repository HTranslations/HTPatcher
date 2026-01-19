package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// ItemDamage represents damage calculation for an item
type ItemDamage struct {
	Critical  bool                       `json:"critical"`
	ElementId int                        `json:"elementId"`
	Formula   any                        `json:"formula"`
	Type      int                        `json:"type"`
	Variance  int                        `json:"variance"`
	Extras    map[string]json.RawMessage `json:"-"`
}

func (i *ItemDamage) UnmarshalJSON(data []byte) error {
	type Alias ItemDamage
	aux := (*Alias)(i)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	i.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(i))
	return nil
}

func (i ItemDamage) MarshalJSON() ([]byte, error) {
	type Alias ItemDamage
	return util.MarshalWithExtras((*Alias)(&i), i.Extras)
}

// ItemEffect represents an effect of using an item
type ItemEffect struct {
	Code   int                        `json:"code"`
	DataId int                        `json:"dataId"`
	Value1 float64                    `json:"value1"`
	Value2 float64                    `json:"value2"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (i *ItemEffect) UnmarshalJSON(data []byte) error {
	type Alias ItemEffect
	aux := (*Alias)(i)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	i.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(i))
	return nil
}

func (i ItemEffect) MarshalJSON() ([]byte, error) {
	type Alias ItemEffect
	return util.MarshalWithExtras((*Alias)(&i), i.Extras)
}

// Item represents a usable item
type Item struct {
	ID          int                        `json:"id"`
	AnimationId int                        `json:"animationId"`
	Consumable  bool                       `json:"consumable"`
	Damage      ItemDamage                 `json:"damage"`
	Description string                     `json:"description"`
	Effects     []ItemEffect               `json:"effects"`
	HitType     int                        `json:"hitType"`
	IconIndex   int                        `json:"iconIndex"`
	ItypeId     int                        `json:"itypeId"`
	Name        string                     `json:"name"`
	Note        string                     `json:"note"`
	Occasion    int                        `json:"occasion"`
	Price       int                        `json:"price"`
	Repeats     int                        `json:"repeats"`
	Scope       int                        `json:"scope"`
	Speed       int                        `json:"speed"`
	SuccessRate int                        `json:"successRate"`
	TpGain      int                        `json:"tpGain"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (i *Item) UnmarshalJSON(data []byte) error {
	type Alias Item
	aux := (*Alias)(i)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	i.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(i))
	return nil
}

func (i Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	return util.MarshalWithExtras((*Alias)(&i), i.Extras)
}

// ItemsData is an array of items
type ItemsData []*Item
