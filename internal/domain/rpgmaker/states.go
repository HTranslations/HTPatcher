package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// State represents a status effect/state
type State struct {
	ID                  int                        `json:"id"`
	AutoRemovalTiming   int                        `json:"autoRemovalTiming"`
	ChanceByDamage      int                        `json:"chanceByDamage"`
	IconIndex           int                        `json:"iconIndex"`
	MaxTurns            int                        `json:"maxTurns"`
	Message1            string                     `json:"message1"`
	Message2            string                     `json:"message2"`
	Message3            string                     `json:"message3"`
	Message4            string                     `json:"message4"`
	MinTurns            int                        `json:"minTurns"`
	Motion              int                        `json:"motion"`
	Name                string                     `json:"name"`
	Note                string                     `json:"note"`
	Overlay             int                        `json:"overlay"`
	Priority            int                        `json:"priority"`
	ReleaseByDamage     bool                       `json:"releaseByDamage"`
	RemoveAtBattleEnd   bool                       `json:"removeAtBattleEnd"`
	RemoveByDamage      bool                       `json:"removeByDamage"`
	RemoveByRestriction bool                       `json:"removeByRestriction"`
	RemoveByWalking     bool                       `json:"removeByWalking"`
	Restriction         int                        `json:"restriction"`
	StepsToRemove       int                        `json:"stepsToRemove"`
	Traits              []Trait                    `json:"traits"`
	MessageType         int                        `json:"messageType"`
	Extras              map[string]json.RawMessage `json:"-"`
}

func (s *State) UnmarshalJSON(data []byte) error {
	type Alias State
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(s))
	return nil
}

func (s State) MarshalJSON() ([]byte, error) {
	type Alias State
	return util.MarshalWithExtras((*Alias)(&s), s.Extras)
}

// StatesData is an array of states
type StatesData []*State
