package main

import (
	"encoding/json"
)

type State struct {
	ID                  int     `json:"id"`
	AutoRemovalTiming   int     `json:"autoRemovalTiming"`
	ChanceByDamage      int     `json:"chanceByDamage"`
	IconIndex           int     `json:"iconIndex"`
	MaxTurns            int     `json:"maxTurns"`
	Message1            string  `json:"message1"`
	Message2            string  `json:"message2"`
	Message3            string  `json:"message3"`
	Message4            string  `json:"message4"`
	MinTurns            int     `json:"minTurns"`
	Motion              int     `json:"motion"`
	Name                string  `json:"name"`
	Note                string  `json:"note"`
	Overlay             int     `json:"overlay"`
	Priority            int     `json:"priority"`
	ReleaseByDamage     bool    `json:"releaseByDamage"`
	RemoveAtBattleEnd   bool    `json:"removeAtBattleEnd"`
	RemoveByDamage      bool    `json:"removeByDamage"`
	RemoveByRestriction bool    `json:"removeByRestriction"`
	RemoveByWalking     bool    `json:"removeByWalking"`
	Restriction         int     `json:"restriction"`
	StepsToRemove       int     `json:"stepsToRemove"`
	Traits              []Trait `json:"traits"`
	MessageType         int     `json:"messageType"`
}

type StatesData []*State

func PatchStates(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var states StatesData

	if err := json.Unmarshal(data, &states); err != nil {
		return nil, err
	}

	for _, state := range states {
		if state == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[GetTranslationKey(state.Name)]
		if ok {
			state.Name = name
		}
		message1, ok := patchInfo.Dictionary[GetTranslationKey(state.Message1)]
		if ok {
			state.Message1 = Wrap(NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		message2, ok := patchInfo.Dictionary[GetTranslationKey(state.Message2)]
		if ok {
			state.Message2 = Wrap(NoNewline(message2), patchInfo.Config.WrapWidth)
		}
		message3, ok := patchInfo.Dictionary[GetTranslationKey(state.Message3)]
		if ok {
			state.Message3 = Wrap(NoNewline(message3), patchInfo.Config.WrapWidth)
		}
		message4, ok := patchInfo.Dictionary[GetTranslationKey(state.Message4)]
		if ok {
			state.Message4 = Wrap(NoNewline(message4), patchInfo.Config.WrapWidth)
		}
	}

	mergedData, err := MergeJsonChanges(data, states)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
