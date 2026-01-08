package rpgmaker

// SkillDamage represents damage calculation for a skill
type SkillDamage struct {
	Critical  bool   `json:"critical"`
	ElementId int    `json:"elementId"`
	Formula   string `json:"formula"`
	Type      int    `json:"type"`
	Variance  int    `json:"variance"`
}

// SkillEffect represents an effect of using a skill
type SkillEffect struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
}

// Skill represents a character skill/ability
type Skill struct {
	ID               int           `json:"id"`
	AnimationId      int           `json:"animationId"`
	Damage           SkillDamage   `json:"damage"`
	Description      string        `json:"description"`
	Effects          []SkillEffect `json:"effects"`
	HitType          int           `json:"hitType"`
	IconIndex        int           `json:"iconIndex"`
	Message1         string        `json:"message1"`
	Message2         string        `json:"message2"`
	MpCost           int           `json:"mpCost"`
	Name             string        `json:"name"`
	Note             string        `json:"note"`
	Occasion         int           `json:"occasion"`
	Repeats          int           `json:"repeats"`
	RequiredWtypeId1 int           `json:"requiredWtypeId1"`
	RequiredWtypeId2 int           `json:"requiredWtypeId2"`
	Scope            int           `json:"scope"`
	Speed            int           `json:"speed"`
	StypeId          int           `json:"stypeId"`
	SuccessRate      int           `json:"successRate"`
	TpCost           int           `json:"tpCost"`
	TpGain           int           `json:"tpGain"`
	MessageType      int           `json:"messageType"`
}

// SkillsData is an array of skills
type SkillsData []*Skill




