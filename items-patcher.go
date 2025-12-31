package main

import (
	"encoding/json"
)

type ItemDamage struct {
	Critical  bool   `json:"critical"`
	ElementId int    `json:"elementId"`
	Formula   string `json:"formula"`
	Type      int    `json:"type"`
	Variance  int    `json:"variance"`
}

type ItemEffect struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
}

type Item struct {
	ID          int          `json:"id"`
	AnimationId int          `json:"animationId"`
	Consumable  bool         `json:"consumable"`
	Damage      ItemDamage   `json:"damage"`
	Description string       `json:"description"`
	Effects     []ItemEffect `json:"effects"`
	HitType     int          `json:"hitType"`
	IconIndex   int          `json:"iconIndex"`
	ItypeId     int          `json:"itypeId"`
	Name        string       `json:"name"`
	Note        string       `json:"note"`
	Occasion    int          `json:"occasion"`
	Price       int          `json:"price"`
	Repeats     int          `json:"repeats"`
	Scope       int          `json:"scope"`
	Speed       int          `json:"speed"`
	SuccessRate int          `json:"successRate"`
	TpGain      int          `json:"tpGain"`
}

type ItemsData []*Item

func PatchItems(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var items ItemsData

	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	for _, item := range items {
		if item == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[GetTranslationKey(item.Name)]
		if ok {
			item.Name = name
		}
		description, ok := patchInfo.Dictionary[GetTranslationKey(item.Description)]
		if ok {
			item.Description = Wrap(NoNewline(description), patchInfo.Config.WrapWidth)
		}
		note, ok := patchInfo.Dictionary[GetTranslationKey(item.Note)]
		if ok {
			item.Note = NoNewline(note)
		}
	}

	mergedData, err := MergeJsonChanges(data, items)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
