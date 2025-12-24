package main

import (
	"encoding/json"
)

type Armor struct {
	ID          int     `json:"id"`
	AtypeID     int     `json:"atypeId"`
	Description string  `json:"description"`
	EtypeID     int     `json:"etypeId"`
	Traits      []Trait `json:"traits"`
	IconIndex   int     `json:"iconIndex"`
	Name        string  `json:"name"`
	Note        string  `json:"note"`
	Params      []int   `json:"params"`
	Price       int     `json:"price"`
}

type ArmorsData []*Armor

func PatchArmors(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var armors ArmorsData

	if err := json.Unmarshal(data, &armors); err != nil {
		return nil, err
	}

	for _, armor := range armors {
		if armor == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[armor.Name]
		if ok {
			armor.Name = name
		}
		description, ok := patchInfo.Dictionary[armor.Description]
		if ok {
			armor.Description = Wrap(NoNewline(description), patchInfo.Config.WrapWidth)
		}
	}

	mergedData, err := MergeJsonChanges(data, armors)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
