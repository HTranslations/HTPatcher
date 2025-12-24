package main

import (
	"encoding/json"
)

type SkillDamage struct {
	Critical  bool   `json:"critical"`
	ElementId int    `json:"elementId"`
	Formula   string `json:"formula"`
	Type      int    `json:"type"`
	Variance  int    `json:"variance"`
}

type SkillEffect struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
}

type Skill struct {
	ID               int           `json:"id"`
	AnimationId      int           `json:"animationId"`
	Damage           SkillDamage   `json:"damage"`
	Description      string        `json:"description"`
	Effects          []SkillEffect `json:"effects"`
	HitType          int           `json:"hitType"`
	IconIndex        int           `json:"iconIndex"`
	Message1         string        `json:"message1"`
	Message2         string        `json:"message2"`
	MpCost           int           `json:"mpCost"`
	Name             string        `json:"name"`
	Note             string        `json:"note"`
	Occasion         int           `json:"occasion"`
	Repeats          int           `json:"repeats"`
	RequiredWtypeId1 int           `json:"requiredWtypeId1"`
	RequiredWtypeId2 int           `json:"requiredWtypeId2"`
	Scope            int           `json:"scope"`
	Speed            int           `json:"speed"`
	StypeId          int           `json:"stypeId"`
	SuccessRate      int           `json:"successRate"`
	TpCost           int           `json:"tpCost"`
	TpGain           int           `json:"tpGain"`
	MessageType      int           `json:"messageType"`
}

type SkillsData []*Skill

func PatchSkills(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var skills SkillsData

	if err := json.Unmarshal(data, &skills); err != nil {
		return nil, err
	}

	for _, skill := range skills {
		if skill == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[skill.Name]
		if ok {
			skill.Name = name
		}
		description, ok := patchInfo.Dictionary[skill.Description]
		if ok {
			skill.Description = Wrap(NoNewline(description), patchInfo.Config.WrapWidth)
		}
		message1, ok := patchInfo.Dictionary[skill.Message1]
		if ok {
			skill.Message1 = Wrap(NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		message2, ok := patchInfo.Dictionary[skill.Message2]
		if ok {
			skill.Message2 = Wrap(NoNewline(message2), patchInfo.Config.WrapWidth)
		}
	}

	mergedData, err := MergeJsonChanges(data, skills)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
