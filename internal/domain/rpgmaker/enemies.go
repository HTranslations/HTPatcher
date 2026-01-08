package rpgmaker

// EnemyAction represents an enemy's action in battle
type EnemyAction struct {
	ConditionParam1 float64 `json:"conditionParam1"`
	ConditionParam2 float64 `json:"conditionParam2"`
	ConditionType   int     `json:"conditionType"`
	Rating          int     `json:"rating"`
	SkillId         int     `json:"skillId"`
}

// DropItem represents an item that can be dropped by an enemy
type DropItem struct {
	DataId      int `json:"dataId"`
	Denominator int `json:"denominator"`
	Kind        int `json:"kind"`
}

// Enemy represents an enemy character
type Enemy struct {
	ID          int           `json:"id"`
	Actions     []EnemyAction `json:"actions"`
	BattlerHue  int           `json:"battlerHue"`
	BattlerName string        `json:"battlerName"`
	DropItems   []DropItem    `json:"dropItems"`
	Exp         int           `json:"exp"`
	Traits      []Trait       `json:"traits"`
	Gold        int           `json:"gold"`
	Name        string        `json:"name"`
	Note        string        `json:"note"`
	Params      []int         `json:"params"`
}

// EnemiesData is an array of enemies
type EnemiesData []*Enemy




