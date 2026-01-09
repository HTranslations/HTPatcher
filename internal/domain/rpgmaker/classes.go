package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// Learning represents a skill learned at a level
type Learning struct {
	Level   int                        `json:"level"`
	Note    string                     `json:"note"`
	SkillId int                        `json:"skillId"`
	Extras  map[string]json.RawMessage `json:"-"`
}

func (l *Learning) UnmarshalJSON(data []byte) error {
	type Alias Learning
	aux := (*Alias)(l)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	l.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(l))
	return nil
}

func (l Learning) MarshalJSON() ([]byte, error) {
	type Alias Learning
	return util.MarshalWithExtras((*Alias)(&l), l.Extras)
}

// Class represents a character class
type Class struct {
	ID        int                        `json:"id"`
	ExpParams []int                      `json:"expParams"`
	Traits    []Trait                    `json:"traits"`
	Learnings []Learning                 `json:"learnings"`
	Name      string                     `json:"name"`
	Note      string                     `json:"note"`
	Params    [][]int                    `json:"params"`
	Extras    map[string]json.RawMessage `json:"-"`
}

func (c *Class) UnmarshalJSON(data []byte) error {
	type Alias Class
	aux := (*Alias)(c)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(c))
	return nil
}

func (c Class) MarshalJSON() ([]byte, error) {
	type Alias Class
	return util.MarshalWithExtras((*Alias)(&c), c.Extras)
}

// ClassesData is an array of classes
type ClassesData []*Class
