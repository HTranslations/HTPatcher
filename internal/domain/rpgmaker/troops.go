package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// TroopMember represents an enemy member of a troop
type TroopMember struct {
	EnemyId int                        `json:"enemyId"`
	X       float64                    `json:"x"`
	Y       float64                    `json:"y"`
	Hidden  bool                       `json:"hidden"`
	Extras  map[string]json.RawMessage `json:"-"`
}

func (t *TroopMember) UnmarshalJSON(data []byte) error {
	type Alias TroopMember
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TroopMember) MarshalJSON() ([]byte, error) {
	type Alias TroopMember
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// TroopPageCondition defines conditions for a troop page to activate
type TroopPageCondition struct {
	ActorHp     int                        `json:"actorHp"`
	ActorId     int                        `json:"actorId"`
	ActorValid  bool                       `json:"actorValid"`
	EnemyHp     int                        `json:"enemyHp"`
	EnemyIndex  int                        `json:"enemyIndex"`
	EnemyValid  bool                       `json:"enemyValid"`
	SwitchId    int                        `json:"switchId"`
	SwitchValid bool                       `json:"switchValid"`
	TurnA       int                        `json:"turnA"`
	TurnB       int                        `json:"turnB"`
	TurnEnding  bool                       `json:"turnEnding"`
	TurnValid   bool                       `json:"turnValid"`
	Extras      map[string]json.RawMessage `json:"-"`
}

func (t *TroopPageCondition) UnmarshalJSON(data []byte) error {
	type Alias TroopPageCondition
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TroopPageCondition) MarshalJSON() ([]byte, error) {
	type Alias TroopPageCondition
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// TroopPage represents a page of troop events
type TroopPage struct {
	Conditions TroopPageCondition         `json:"conditions"`
	List       []*EventCommand            `json:"list"`
	Span       int                        `json:"span"`
	Extras     map[string]json.RawMessage `json:"-"`
}

func (t *TroopPage) UnmarshalJSON(data []byte) error {
	type Alias TroopPage
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TroopPage) MarshalJSON() ([]byte, error) {
	type Alias TroopPage
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// Troop represents an enemy troop/group
type Troop struct {
	ID      int                        `json:"id"`
	Members []TroopMember              `json:"members"`
	Name    string                     `json:"name"`
	Pages   []TroopPage                `json:"pages"`
	Extras  map[string]json.RawMessage `json:"-"`
}

func (t *Troop) UnmarshalJSON(data []byte) error {
	type Alias Troop
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t Troop) MarshalJSON() ([]byte, error) {
	type Alias Troop
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// TroopsData is an array of troops
type TroopsData []*Troop
