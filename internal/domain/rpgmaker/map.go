package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// MapEventPageCondition defines conditions for a map event page
type MapEventPageCondition struct {
	ActorId         int                        `json:"actorId"`
	ActorValid      bool                       `json:"actorValid"`
	ItemId          int                        `json:"itemId"`
	ItemValid       bool                       `json:"itemValid"`
	SelfSwitchCh    string                     `json:"selfSwitchCh"`
	SelfSwitchValid bool                       `json:"selfSwitchValid"`
	Switch1Id       int                        `json:"switch1Id"`
	Switch1Valid    bool                       `json:"switch1Valid"`
	Switch2Id       int                        `json:"switch2Id"`
	Switch2Valid    bool                       `json:"switch2Valid"`
	VariableId      int                        `json:"variableId"`
	VariableValid   bool                       `json:"variableValid"`
	VariableValue   int                        `json:"variableValue"`
	Extras          map[string]json.RawMessage `json:"-"`
}

func (m *MapEventPageCondition) UnmarshalJSON(data []byte) error {
	type Alias MapEventPageCondition
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapEventPageCondition) MarshalJSON() ([]byte, error) {
	type Alias MapEventPageCondition
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}

// MapEventPageImage defines the visual appearance of an event
type MapEventPageImage struct {
	TileId        int                        `json:"tileId"`
	CharacterName string                     `json:"characterName"`
	Direction     int                        `json:"direction"`
	Pattern       int                        `json:"pattern"`
	Extras        map[string]json.RawMessage `json:"-"`
}

func (m *MapEventPageImage) UnmarshalJSON(data []byte) error {
	type Alias MapEventPageImage
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapEventPageImage) MarshalJSON() ([]byte, error) {
	type Alias MapEventPageImage
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}

// MapEventPageMoveRoute defines a custom move route
type MapEventPageMoveRoute struct {
	List      []any                      `json:"list"`
	Repeat    bool                       `json:"repeat"`
	Skippable bool                       `json:"skippable"`
	Wait      bool                       `json:"wait"`
	Extras    map[string]json.RawMessage `json:"-"`
}

func (m *MapEventPageMoveRoute) UnmarshalJSON(data []byte) error {
	type Alias MapEventPageMoveRoute
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapEventPageMoveRoute) MarshalJSON() ([]byte, error) {
	type Alias MapEventPageMoveRoute
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}

// MapEventPage represents a page of a map event
type MapEventPage struct {
	Conditions    MapEventPageCondition      `json:"conditions"`
	DirectionFix  bool                       `json:"directionFix"`
	Image         MapEventPageImage          `json:"image"`
	List          []*EventCommand            `json:"list"`
	MoveFrequency int                        `json:"moveFrequency"`
	MoveRoute     MapEventPageMoveRoute      `json:"moveRoute"`
	MoveSpeed     int                        `json:"moveSpeed"`
	MoveType      int                        `json:"moveType"`
	PriorityType  int                        `json:"priorityType"`
	StepAnime     bool                       `json:"stepAnime"`
	Through       bool                       `json:"through"`
	WalkAnime     bool                       `json:"walkAnime"`
	Extras        map[string]json.RawMessage `json:"-"`
}

func (m *MapEventPage) UnmarshalJSON(data []byte) error {
	type Alias MapEventPage
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapEventPage) MarshalJSON() ([]byte, error) {
	type Alias MapEventPage
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}

// MapEvent represents an event on a map
type MapEvent struct {
	ID     int                        `json:"id"`
	Name   string                     `json:"name"`
	Note   string                     `json:"note"`
	Pages  []MapEventPage             `json:"pages"`
	X      int                        `json:"x"`
	Y      int                        `json:"y"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (m *MapEvent) UnmarshalJSON(data []byte) error {
	type Alias MapEvent
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapEvent) MarshalJSON() ([]byte, error) {
	type Alias MapEvent
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}

// BgmBgs represents background music or sound settings
type BgmBgs struct {
	Name   string                     `json:"name"`
	Pan    int                        `json:"pan"`
	Pitch  int                        `json:"pitch"`
	Volume int                        `json:"volume"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (b *BgmBgs) UnmarshalJSON(data []byte) error {
	type Alias BgmBgs
	aux := (*Alias)(b)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	b.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(b))
	return nil
}

func (b BgmBgs) MarshalJSON() ([]byte, error) {
	type Alias BgmBgs
	return util.MarshalWithExtras((*Alias)(&b), b.Extras)
}

// MapData represents a game map
type MapData struct {
	AutoplayBgm       bool                       `json:"autoplayBgm"`
	AutoplayBgs       bool                       `json:"autoplayBgs"`
	Battleback1Name   string                     `json:"battleback1Name"`
	Battleback2Name   string                     `json:"battleback2Name"`
	Bgm               BgmBgs                     `json:"bgm"`
	Bgs               BgmBgs                     `json:"bgs"`
	DisableDashing    bool                       `json:"disableDashing"`
	DisplayName       string                     `json:"displayName"`
	EncounterList     []any                      `json:"encounterList"`
	EncounterStep     int                        `json:"encounterStep"`
	Events            []*MapEvent                `json:"events"`
	Height            int                        `json:"height"`
	Note              string                     `json:"note"`
	ParallaxLoopX     bool                       `json:"parallaxLoopX"`
	ParallaxLoopY     bool                       `json:"parallaxLoopY"`
	ParallaxName      string                     `json:"parallaxName"`
	ParallaxShow      bool                       `json:"parallaxShow"`
	ParallaxSx        int                        `json:"parallaxSx"`
	ParallaxSy        int                        `json:"parallaxSy"`
	ScrollType        int                        `json:"scrollType"`
	SpecifyBattleback bool                       `json:"specifyBattleback"`
	TilesetId         int                        `json:"tilesetId"`
	Width             int                        `json:"width"`
	Data              []int                      `json:"data"`
	Extras            map[string]json.RawMessage `json:"-"`
}

func (m *MapData) UnmarshalJSON(data []byte) error {
	type Alias MapData
	aux := (*Alias)(m)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(m))
	return nil
}

func (m MapData) MarshalJSON() ([]byte, error) {
	type Alias MapData
	return util.MarshalWithExtras((*Alias)(&m), m.Extras)
}
