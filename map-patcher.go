package main

import (
	"encoding/json"
)

type MapEventPageCondition struct {
	ActorId         int    `json:"actorId"`
	ActorValid      bool   `json:"actorValid"`
	ItemId          int    `json:"itemId"`
	ItemValid       bool   `json:"itemValid"`
	SelfSwitchCh    string `json:"selfSwitchCh"`
	SelfSwitchValid bool   `json:"selfSwitchValid"`
	Switch1Id       int    `json:"switch1Id"`
	Switch1Valid    bool   `json:"switch1Valid"`
	Switch2Id       int    `json:"switch2Id"`
	Switch2Valid    bool   `json:"switch2Valid"`
	VariableId      int    `json:"variableId"`
	VariableValid   bool   `json:"variableValid"`
	VariableValue   int    `json:"variableValue"`
}

type MapEventPageImage struct {
	TileId        int    `json:"tileId"`
	CharacterName string `json:"characterName"`
	Direction     int    `json:"direction"`
	Pattern       int    `json:"pattern"`
}

type MapEventPageMoveRoute struct {
	List      []any `json:"list"`
	Repeat    bool  `json:"repeat"`
	Skippable bool  `json:"skippable"`
	Wait      bool  `json:"wait"`
}

type MapEventPage struct {
	Conditions    MapEventPageCondition `json:"conditions"`
	DirectionFix  bool                  `json:"directionFix"`
	Image         MapEventPageImage     `json:"image"`
	List          []*EventCommand       `json:"list"`
	MoveFrequency int                   `json:"moveFrequency"`
	MoveRoute     MapEventPageMoveRoute `json:"moveRoute"`
	MoveSpeed     int                   `json:"moveSpeed"`
	MoveType      int                   `json:"moveType"`
	PriorityType  int                   `json:"priorityType"`
	StepAnime     bool                  `json:"stepAnime"`
	Through       bool                  `json:"through"`
	WalkAnime     bool                  `json:"walkAnime"`
}

type MapEvent struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Note  string         `json:"note"`
	Pages []MapEventPage `json:"pages"`
	X     int            `json:"x"`
	Y     int            `json:"y"`
}

type BgmBgs struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
}

type MapData struct {
	AutoplayBgm       bool        `json:"autoplayBgm"`
	AutoplayBgs       bool        `json:"autoplayBgs"`
	Battleback1Name   string      `json:"battleback1Name"`
	Battleback2Name   string      `json:"battleback2Name"`
	Bgm               BgmBgs      `json:"bgm"`
	Bgs               BgmBgs      `json:"bgs"`
	DisableDashing    bool        `json:"disableDashing"`
	DisplayName       string      `json:"displayName"`
	EncounterList     []any       `json:"encounterList"`
	EncounterStep     int         `json:"encounterStep"`
	Events            []*MapEvent `json:"events"`
	Height            int         `json:"height"`
	Note              string      `json:"note"`
	ParallaxLoopX     bool        `json:"parallaxLoopX"`
	ParallaxLoopY     bool        `json:"parallaxLoopY"`
	ParallaxName      string      `json:"parallaxName"`
	ParallaxShow      bool        `json:"parallaxShow"`
	ParallaxSx        int         `json:"parallaxSx"`
	ParallaxSy        int         `json:"parallaxSy"`
	ScrollType        int         `json:"scrollType"`
	SpecifyBattleback bool        `json:"specifyBattleback"`
	TilesetId         int         `json:"tilesetId"`
	Width             int         `json:"width"`
	Data              []int       `json:"data"`
}

func PatchMap(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var mapData MapData

	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}

	if displayName, ok := patchInfo.Dictionary[mapData.DisplayName]; ok {
		mapData.DisplayName = displayName
	}

	for _, event := range mapData.Events {
		if event == nil {
			continue
		}

		for i := range event.Pages {
			newCommands, err := PatchCommands(event.Pages[i].List, patchInfo)
			if err != nil {
				return nil, err
			}
			event.Pages[i].List = newCommands
		}
	}

	mergedData, err := MergeJsonChanges(data, mapData)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
