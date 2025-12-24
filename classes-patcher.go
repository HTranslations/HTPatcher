package main

import (
	"encoding/json"
)

type Learning struct {
	Level   int    `json:"level"`
	Note    string `json:"note"`
	SkillId int    `json:"skillId"`
}

type Class struct {
	ID        int        `json:"id"`
	ExpParams []int      `json:"expParams"`
	Traits    []Trait    `json:"traits"`
	Learnings []Learning `json:"learnings"`
	Name      string     `json:"name"`
	Note      string     `json:"note"`
	Params    [][]int    `json:"params"`
}

type ClassesData []*Class

func PatchClasses(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var classes ClassesData

	if err := json.Unmarshal(data, &classes); err != nil {
		return nil, err
	}

	for _, class := range classes {
		if class == nil {
			continue
		}
		name, ok := patchInfo.Dictionary[class.Name]
		if ok {
			class.Name = name
		}
		note, ok := patchInfo.Dictionary[class.Note]
		if ok {
			class.Note = note
		}
	}

	mergedData, err := MergeJsonChanges(data, classes)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
