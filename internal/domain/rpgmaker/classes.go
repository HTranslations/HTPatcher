package rpgmaker

// Learning represents a skill learned at a level
type Learning struct {
	Level   int    `json:"level"`
	Note    string `json:"note"`
	SkillId int    `json:"skillId"`
}

// Class represents a character class
type Class struct {
	ID        int        `json:"id"`
	ExpParams []int      `json:"expParams"`
	Traits    []Trait    `json:"traits"`
	Learnings []Learning `json:"learnings"`
	Name      string     `json:"name"`
	Note      string     `json:"note"`
	Params    [][]int    `json:"params"`
}

// ClassesData is an array of classes
type ClassesData []*Class




