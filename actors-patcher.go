package main

import (
	"encoding/json"
)

type Trait struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value  float64 `json:"value"`
}

type Actor struct {
	ID             int     `json:"id"`
	BattlerName    string  `json:"battlerName"`
	CharacterIndex int     `json:"characterIndex"`
	CharacterName  string  `json:"characterName"`
	ClassId        int     `json:"classId"`
	Equips         []int   `json:"equips"`
	FaceIndex      int     `json:"faceIndex"`
	FaceName       string  `json:"faceName"`
	Traits         []Trait `json:"traits"`
	InitialLevel   int     `json:"initialLevel"`
	MaxLevel       int     `json:"maxLevel"`
	Name           string  `json:"name"`
	Nickname       string  `json:"nickname"`
	Note           string  `json:"note"`
	Profile        string  `json:"profile"`
}

type ActorsData []*Actor

func PatchActors(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var actors ActorsData

	if err := json.Unmarshal(data, &actors); err != nil {
		return nil, err
	}

	for _, actor := range actors {
		if actor == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[GetTranslationKey(actor.Name)]
		if ok {
			actor.Name = name
		}
	}

	mergedData, err := MergeJsonChanges(data, actors)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
