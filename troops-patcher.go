package main

import (
	"encoding/json"
)

type TroopMember struct {
	EnemyId int  `json:"enemyId"`
	X       int  `json:"x"`
	Y       int  `json:"y"`
	Hidden  bool `json:"hidden"`
}

type TroopPageCondition struct {
	ActorHp     int  `json:"actorHp"`
	ActorId     int  `json:"actorId"`
	ActorValid  bool `json:"actorValid"`
	EnemyHp     int  `json:"enemyHp"`
	EnemyIndex  int  `json:"enemyIndex"`
	EnemyValid  bool `json:"enemyValid"`
	SwitchId    int  `json:"switchId"`
	SwitchValid bool `json:"switchValid"`
	TurnA       int  `json:"turnA"`
	TurnB       int  `json:"turnB"`
	TurnEnding  bool `json:"turnEnding"`
	TurnValid   bool `json:"turnValid"`
}

type TroopPage struct {
	Conditions TroopPageCondition `json:"conditions"`
	List       []*EventCommand    `json:"list"`
	Span       int                `json:"span"`
}

type Troop struct {
	ID      int           `json:"id"`
	Members []TroopMember `json:"members"`
	Name    string        `json:"name"`
	Pages   []TroopPage   `json:"pages"`
}

type TroopsData []*Troop

func PatchTroops(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var troops TroopsData

	if err := json.Unmarshal(data, &troops); err != nil {
		return nil, err
	}

	for _, troop := range troops {
		if troop == nil {
			continue
		}

		if name, ok := patchInfo.Dictionary[GetTranslationKey(troop.Name)]; ok {
			troop.Name = name
		}

		for i := range troop.Pages {
			newCommands, err := PatchCommands(troop.Pages[i].List, patchInfo)
			if err != nil {
				return nil, err
			}
			troop.Pages[i].List = newCommands
		}
	}

	mergedData, err := MergeJsonChanges(data, troops)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
