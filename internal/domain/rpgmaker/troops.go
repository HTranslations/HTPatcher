package rpgmaker

// TroopMember represents an enemy member of a troop
type TroopMember struct {
	EnemyId int  `json:"enemyId"`
	X       int  `json:"x"`
	Y       int  `json:"y"`
	Hidden  bool `json:"hidden"`
}

// TroopPageCondition defines conditions for a troop page to activate
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

// TroopPage represents a page of troop events
type TroopPage struct {
	Conditions TroopPageCondition `json:"conditions"`
	List       []*EventCommand    `json:"list"`
	Span       int                `json:"span"`
}

// Troop represents an enemy troop/group
type Troop struct {
	ID      int           `json:"id"`
	Members []TroopMember `json:"members"`
	Name    string        `json:"name"`
	Pages   []TroopPage   `json:"pages"`
}

// TroopsData is an array of troops
type TroopsData []*Troop




