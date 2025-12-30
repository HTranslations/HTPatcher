package main

import (
	"encoding/json"
)

type Weapon struct {
	ID          int     `json:"id"`
	AnimationId int     `json:"animationId"`
	Description string  `json:"description"`
	EtypeId     int     `json:"etypeId"`
	Traits      []Trait `json:"traits"`
	IconIndex   int     `json:"iconIndex"`
	Name        string  `json:"name"`
	Note        string  `json:"note"`
	Params      []int   `json:"params"`
	Price       int     `json:"price"`
	WtypeId     int     `json:"wtypeId"`
}

type WeaponsData []*Weapon

func PatchWeapons(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var weapons WeaponsData

	if err := json.Unmarshal(data, &weapons); err != nil {
		return nil, err
	}

	for _, weapon := range weapons {
		if weapon == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[GetTranslationKey(weapon.Name)]
		if ok {
			weapon.Name = name
		}
		description, ok := patchInfo.Dictionary[GetTranslationKey(weapon.Description)]
		if ok {
			weapon.Description = Wrap(NoNewline(description), patchInfo.Config.WrapWidth)
		}
	}

	return json.Marshal(weapons)
}

