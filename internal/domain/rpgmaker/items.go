package rpgmaker

// ItemDamage represents damage calculation for an item
type ItemDamage struct {
	Critical  bool   `json:"critical"`
	ElementId int    `json:"elementId"`
	Formula   string `json:"formula"`
	Type      int    `json:"type"`
	Variance  int    `json:"variance"`
}

// ItemEffect represents an effect of using an item
type ItemEffect struct {
	Code   int     `json:"code"`
	DataId int     `json:"dataId"`
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
}

// Item represents a usable item
type Item struct {
	ID          int          `json:"id"`
	AnimationId int          `json:"animationId"`
	Consumable  bool         `json:"consumable"`
	Damage      ItemDamage   `json:"damage"`
	Description string       `json:"description"`
	Effects     []ItemEffect `json:"effects"`
	HitType     int          `json:"hitType"`
	IconIndex   int          `json:"iconIndex"`
	ItypeId     int          `json:"itypeId"`
	Name        string       `json:"name"`
	Note        string       `json:"note"`
	Occasion    int          `json:"occasion"`
	Price       int          `json:"price"`
	Repeats     int          `json:"repeats"`
	Scope       int          `json:"scope"`
	Speed       int          `json:"speed"`
	SuccessRate int          `json:"successRate"`
	TpGain      int          `json:"tpGain"`
}

// ItemsData is an array of items
type ItemsData []*Item




