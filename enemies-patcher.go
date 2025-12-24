package main

import (
	"encoding/json"
)

type EnemyAction struct {
	ConditionParam1 float64 `json:"conditionParam1"`
	ConditionParam2 float64 `json:"conditionParam2"`
	ConditionType   int     `json:"conditionType"`
	Rating          int     `json:"rating"`
	SkillId         int     `json:"skillId"`
}

type DropItem struct {
	DataId      int `json:"dataId"`
	Denominator int `json:"denominator"`
	Kind        int `json:"kind"`
}

type Enemy struct {
	ID          int           `json:"id"`
	Actions     []EnemyAction `json:"actions"`
	BattlerHue  int           `json:"battlerHue"`
	BattlerName string        `json:"battlerName"`
	DropItems   []DropItem    `json:"dropItems"`
	Exp         int           `json:"exp"`
	Traits      []Trait       `json:"traits"`
	Gold        int           `json:"gold"`
	Name        string        `json:"name"`
	Note        string        `json:"note"`
	Params      []int         `json:"params"`
}

type EnemiesData []*Enemy

func PatchEnemies(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var enemies EnemiesData

	if err := json.Unmarshal(data, &enemies); err != nil {
		return nil, err
	}

	for _, enemy := range enemies {
		if enemy == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[enemy.Name]
		if ok {
			enemy.Name = name
		}
		note, ok := patchInfo.Dictionary[enemy.Note]
		if ok {
			enemy.Note = note
		}
	}

	mergedData, err := MergeJsonChanges(data, enemies)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
