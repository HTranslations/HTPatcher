package rpgmaker

// MapEventPageCondition defines conditions for a map event page
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

// MapEventPageImage defines the visual appearance of an event
type MapEventPageImage struct {
	TileId        int    `json:"tileId"`
	CharacterName string `json:"characterName"`
	Direction     int    `json:"direction"`
	Pattern       int    `json:"pattern"`
}

// MapEventPageMoveRoute defines a custom move route
type MapEventPageMoveRoute struct {
	List      []any `json:"list"`
	Repeat    bool  `json:"repeat"`
	Skippable bool  `json:"skippable"`
	Wait      bool  `json:"wait"`
}

// MapEventPage represents a page of a map event
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

// MapEvent represents an event on a map
type MapEvent struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Note  string         `json:"note"`
	Pages []MapEventPage `json:"pages"`
	X     int            `json:"x"`
	Y     int            `json:"y"`
}

// BgmBgs represents background music or sound settings
type BgmBgs struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
}

// MapData represents a game map
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




