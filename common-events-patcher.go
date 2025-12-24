package main

import (
	"encoding/json"
)

type CommonEvent struct {
	ID       int             `json:"id"`
	List     []*EventCommand `json:"list"`
	Name     string          `json:"name"`
	SwitchId int             `json:"switchId"`
	Trigger  int             `json:"trigger"`
}

type CommonEventsData []*CommonEvent

func PatchCommonEvents(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var commonEvents CommonEventsData

	if err := json.Unmarshal(data, &commonEvents); err != nil {
		return nil, err
	}

	for _, commonEvent := range commonEvents {
		if commonEvent == nil {
			continue
		}

		new_commands, err := PatchCommands(commonEvent.List, patchInfo)
		if err != nil {
			return nil, err
		}
		commonEvent.List = new_commands
	}

	mergedData, err := MergeJsonChanges(data, commonEvents)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
