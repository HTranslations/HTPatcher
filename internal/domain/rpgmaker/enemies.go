package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// EnemyAction represents an enemy's action in battle
type EnemyAction struct {
	ConditionParam1 float64                    `json:"conditionParam1"`
	ConditionParam2 float64                    `json:"conditionParam2"`
	ConditionType   int                        `json:"conditionType"`
	Rating          int                        `json:"rating"`
	SkillId         int                        `json:"skillId"`
	Extras          map[string]json.RawMessage `json:"-"`
}

func (e *EnemyAction) UnmarshalJSON(data []byte) error {
	type Alias EnemyAction
	aux := (*Alias)(e)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(e))
	return nil
}

func (e EnemyAction) MarshalJSON() ([]byte, error) {
	type Alias EnemyAction
	return util.MarshalWithExtras((*Alias)(&e), e.Extras)
}

// DropItem represents an item that can be dropped by an enemy
type DropItem struct {
	DataId      int                        `json:"dataId"`
	Denominator int                        `json:"denominator"`
	Kind        int                        `json:"kind"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (d *DropItem) UnmarshalJSON(data []byte) error {
	type Alias DropItem
	aux := (*Alias)(d)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	d.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(d))
	return nil
}

func (d DropItem) MarshalJSON() ([]byte, error) {
	type Alias DropItem
	return util.MarshalWithExtras((*Alias)(&d), d.Extras)
}

// Enemy represents an enemy character
type Enemy struct {
	ID          int                        `json:"id"`
	Actions     []EnemyAction              `json:"actions"`
	BattlerHue  int                        `json:"battlerHue"`
	BattlerName string                     `json:"battlerName"`
	DropItems   []DropItem                 `json:"dropItems"`
	Exp         int                        `json:"exp"`
	Traits      []Trait                    `json:"traits"`
	Gold        int                        `json:"gold"`
	Name        string                     `json:"name"`
	Note        string                     `json:"note"`
	Params      []any                      `json:"params"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (e *Enemy) UnmarshalJSON(data []byte) error {
	type Alias Enemy
	aux := (*Alias)(e)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(e))
	return nil
}

func (e Enemy) MarshalJSON() ([]byte, error) {
	type Alias Enemy
	return util.MarshalWithExtras((*Alias)(&e), e.Extras)
}

// EnemiesData is an array of enemies
type EnemiesData []*Enemy
