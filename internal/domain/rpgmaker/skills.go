package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// SkillDamage represents damage calculation for a skill
type SkillDamage struct {
	Critical  bool                       `json:"critical"`
	ElementId int                        `json:"elementId"`
	Formula   any                        `json:"formula"`
	Type      int                        `json:"type"`
	Variance  int                        `json:"variance"`
	Extras    map[string]json.RawMessage `json:"-"`
}

func (s *SkillDamage) UnmarshalJSON(data []byte) error {
	type Alias SkillDamage
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(s))
	return nil
}

func (s SkillDamage) MarshalJSON() ([]byte, error) {
	type Alias SkillDamage
	return util.MarshalWithExtras((*Alias)(&s), s.Extras)
}

// SkillEffect represents an effect of using a skill
type SkillEffect struct {
	Code   int                        `json:"code"`
	DataId int                        `json:"dataId"`
	Value1 float64                    `json:"value1"`
	Value2 float64                    `json:"value2"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (s *SkillEffect) UnmarshalJSON(data []byte) error {
	type Alias SkillEffect
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(s))
	return nil
}

func (s SkillEffect) MarshalJSON() ([]byte, error) {
	type Alias SkillEffect
	return util.MarshalWithExtras((*Alias)(&s), s.Extras)
}

// Skill represents a character skill/ability
type Skill struct {
	ID               int                        `json:"id"`
	AnimationId      int                        `json:"animationId"`
	Damage           SkillDamage                `json:"damage"`
	Description      string                     `json:"description"`
	Effects          []SkillEffect              `json:"effects"`
	HitType          int                        `json:"hitType"`
	IconIndex        int                        `json:"iconIndex"`
	Message1         string                     `json:"message1"`
	Message2         string                     `json:"message2"`
	MpCost           int                        `json:"mpCost"`
	Name             string                     `json:"name"`
	Note             string                     `json:"note"`
	Occasion         int                        `json:"occasion"`
	Repeats          int                        `json:"repeats"`
	RequiredWtypeId1 int                        `json:"requiredWtypeId1"`
	RequiredWtypeId2 int                        `json:"requiredWtypeId2"`
	Scope            int                        `json:"scope"`
	Speed            int                        `json:"speed"`
	StypeId          int                        `json:"stypeId"`
	SuccessRate      int                        `json:"successRate"`
	TpCost           int                        `json:"tpCost"`
	TpGain           int                        `json:"tpGain"`
	MessageType      int                        `json:"messageType"`
	Extras           map[string]json.RawMessage `json:"-"`
}

func (s *Skill) UnmarshalJSON(data []byte) error {
	type Alias Skill
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(s))
	return nil
}

func (s Skill) MarshalJSON() ([]byte, error) {
	type Alias Skill
	return util.MarshalWithExtras((*Alias)(&s), s.Extras)
}

// SkillsData is an array of skills
type SkillsData []*Skill
