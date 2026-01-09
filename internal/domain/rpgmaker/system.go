package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// Bgm represents background music settings
type Bgm struct {
	Name   string                     `json:"name"`
	Pan    int                        `json:"pan"`
	Pitch  int                        `json:"pitch"`
	Volume int                        `json:"volume"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (b *Bgm) UnmarshalJSON(data []byte) error {
	type Alias Bgm
	aux := (*Alias)(b)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	b.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(b))
	return nil
}

func (b Bgm) MarshalJSON() ([]byte, error) {
	type Alias Bgm
	return util.MarshalWithExtras((*Alias)(&b), b.Extras)
}

// Vehicle represents a vehicle (boat, ship, airship)
type Vehicle struct {
	Bgm            Bgm                        `json:"bgm"`
	CharacterIndex int                        `json:"characterIndex"`
	CharacterName  string                     `json:"characterName"`
	StartMapId     int                        `json:"startMapId"`
	StartX         int                        `json:"startX"`
	StartY         int                        `json:"startY"`
	Extras         map[string]json.RawMessage `json:"-"`
}

func (v *Vehicle) UnmarshalJSON(data []byte) error {
	type Alias Vehicle
	aux := (*Alias)(v)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(v))
	return nil
}

func (v Vehicle) MarshalJSON() ([]byte, error) {
	type Alias Vehicle
	return util.MarshalWithExtras((*Alias)(&v), v.Extras)
}

// AttackMotion represents an attack animation
type AttackMotion struct {
	Type          int                        `json:"type"`
	WeaponImageId int                        `json:"weaponImageId"`
	Extras        map[string]json.RawMessage `json:"-"`
}

func (a *AttackMotion) UnmarshalJSON(data []byte) error {
	type Alias AttackMotion
	aux := (*Alias)(a)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(a))
	return nil
}

func (a AttackMotion) MarshalJSON() ([]byte, error) {
	type Alias AttackMotion
	return util.MarshalWithExtras((*Alias)(&a), a.Extras)
}

// TestBattler represents a test battle character
type TestBattler struct {
	ActorId int                        `json:"actorId"`
	Level   int                        `json:"level"`
	Equips  []int                      `json:"equips"`
	Extras  map[string]json.RawMessage `json:"-"`
}

func (t *TestBattler) UnmarshalJSON(data []byte) error {
	type Alias TestBattler
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TestBattler) MarshalJSON() ([]byte, error) {
	type Alias TestBattler
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// TitleCommandWindow represents title screen command window settings
type TitleCommandWindow struct {
	Background int                        `json:"background"`
	OffsetX    int                        `json:"offsetX"`
	OffsetY    int                        `json:"offsetY"`
	Extras     map[string]json.RawMessage `json:"-"`
}

func (t *TitleCommandWindow) UnmarshalJSON(data []byte) error {
	type Alias TitleCommandWindow
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TitleCommandWindow) MarshalJSON() ([]byte, error) {
	type Alias TitleCommandWindow
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// Advanced contains advanced game settings
type Advanced struct {
	GameId             int                        `json:"gameId"`
	ScreenWidth        int                        `json:"screenWidth"`
	ScreenHeight       int                        `json:"screenHeight"`
	UiAreaWidth        int                        `json:"uiAreaWidth"`
	UiAreaHeight       int                        `json:"uiAreaHeight"`
	NumberFontFilename string                     `json:"numberFontFilename"`
	FallbackFonts      string                     `json:"fallbackFonts"`
	FontSize           int                        `json:"fontSize"`
	MainFontFilename   string                     `json:"mainFontFilename"`
	WindowOpacity      int                        `json:"windowOpacity"`
	ScreenScale        int                        `json:"screenScale"`
	Extras             map[string]json.RawMessage `json:"-"`
}

func (a *Advanced) UnmarshalJSON(data []byte) error {
	type Alias Advanced
	aux := (*Alias)(a)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(a))
	return nil
}

func (a Advanced) MarshalJSON() ([]byte, error) {
	type Alias Advanced
	return util.MarshalWithExtras((*Alias)(&a), a.Extras)
}

// TermsMessages contains system message text
type TermsMessages struct {
	AlwaysDash      string                     `json:"alwaysDash"`
	CommandRemember string                     `json:"commandRemember"`
	TouchUI         string                     `json:"touchUI"`
	BgmVolume       string                     `json:"bgmVolume"`
	BgsVolume       string                     `json:"bgsVolume"`
	MeVolume        string                     `json:"meVolume"`
	SeVolume        string                     `json:"seVolume"`
	Possession      string                     `json:"possession"`
	ExpTotal        string                     `json:"expTotal"`
	ExpNext         string                     `json:"expNext"`
	SaveMessage     string                     `json:"saveMessage"`
	LoadMessage     string                     `json:"loadMessage"`
	File            string                     `json:"file"`
	Autosave        string                     `json:"autosave"`
	PartyName       string                     `json:"partyName"`
	Emerge          string                     `json:"emerge"`
	Preemptive      string                     `json:"preemptive"`
	Surprise        string                     `json:"surprise"`
	EscapeStart     string                     `json:"escapeStart"`
	EscapeFailure   string                     `json:"escapeFailure"`
	Victory         string                     `json:"victory"`
	Defeat          string                     `json:"defeat"`
	ObtainExp       string                     `json:"obtainExp"`
	ObtainGold      string                     `json:"obtainGold"`
	ObtainItem      string                     `json:"obtainItem"`
	LevelUp         string                     `json:"levelUp"`
	ObtainSkill     string                     `json:"obtainSkill"`
	UseItem         string                     `json:"useItem"`
	CriticalToEnemy string                     `json:"criticalToEnemy"`
	CriticalToActor string                     `json:"criticalToActor"`
	ActorDamage     string                     `json:"actorDamage"`
	ActorRecovery   string                     `json:"actorRecovery"`
	ActorGain       string                     `json:"actorGain"`
	ActorLoss       string                     `json:"actorLoss"`
	ActorDrain      string                     `json:"actorDrain"`
	ActorNoDamage   string                     `json:"actorNoDamage"`
	ActorNoHit      string                     `json:"actorNoHit"`
	EnemyDamage     string                     `json:"enemyDamage"`
	EnemyRecovery   string                     `json:"enemyRecovery"`
	EnemyGain       string                     `json:"enemyGain"`
	EnemyLoss       string                     `json:"enemyLoss"`
	EnemyDrain      string                     `json:"enemyDrain"`
	EnemyNoDamage   string                     `json:"enemyNoDamage"`
	EnemyNoHit      string                     `json:"enemyNoHit"`
	Evasion         string                     `json:"evasion"`
	MagicEvasion    string                     `json:"magicEvasion"`
	MagicReflection string                     `json:"magicReflection"`
	CounterAttack   string                     `json:"counterAttack"`
	Substitute      string                     `json:"substitute"`
	BuffAdd         string                     `json:"buffAdd"`
	DebuffAdd       string                     `json:"debuffAdd"`
	BuffRemove      string                     `json:"buffRemove"`
	ActionFailure   string                     `json:"actionFailure"`
	Extras          map[string]json.RawMessage `json:"-"`
}

func (t *TermsMessages) UnmarshalJSON(data []byte) error {
	type Alias TermsMessages
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t TermsMessages) MarshalJSON() ([]byte, error) {
	type Alias TermsMessages
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// Terms contains game terminology
type Terms struct {
	Basic    []string                   `json:"basic"`
	Commands []*string                  `json:"commands"`
	Params   []string                   `json:"params"`
	Messages TermsMessages              `json:"messages"`
	Extras   map[string]json.RawMessage `json:"-"`
}

func (t *Terms) UnmarshalJSON(data []byte) error {
	type Alias Terms
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t Terms) MarshalJSON() ([]byte, error) {
	type Alias Terms
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// System represents the system data (System.json)
type System struct {
	Advanced           Advanced                   `json:"advanced"`
	Airship            Vehicle                    `json:"airship"`
	ArmorTypes         []string                   `json:"armorTypes"`
	AttackMotions      []AttackMotion             `json:"attackMotions"`
	BattleBgm          Bgm                        `json:"battleBgm"`
	Battleback1Name    string                     `json:"battleback1Name"`
	Battleback2Name    string                     `json:"battleback2Name"`
	BattlerHue         int                        `json:"battlerHue"`
	BattlerName        string                     `json:"battlerName"`
	BattleSystem       int                        `json:"battleSystem"`
	Boat               Vehicle                    `json:"boat"`
	CurrencyUnit       string                     `json:"currencyUnit"`
	DefeatMe           Bgm                        `json:"defeatMe"`
	EditMapId          int                        `json:"editMapId"`
	Elements           []string                   `json:"elements"`
	EquipTypes         []string                   `json:"equipTypes"`
	GameTitle          string                     `json:"gameTitle"`
	GameoverMe         Bgm                        `json:"gameoverMe"`
	ItemCategories     []bool                     `json:"itemCategories"`
	Locale             string                     `json:"locale"`
	MagicSkills        []int                      `json:"magicSkills"`
	MenuCommands       []bool                     `json:"menuCommands"`
	OptAutosave        bool                       `json:"optAutosave"`
	OptDisplayTp       bool                       `json:"optDisplayTp"`
	OptDrawTitle       bool                       `json:"optDrawTitle"`
	OptExtraExp        bool                       `json:"optExtraExp"`
	OptFloorDeath      bool                       `json:"optFloorDeath"`
	OptFollowers       bool                       `json:"optFollowers"`
	OptKeyItemsNumber  bool                       `json:"optKeyItemsNumber"`
	OptSideView        bool                       `json:"optSideView"`
	OptSlipDeath       bool                       `json:"optSlipDeath"`
	OptTransparent     bool                       `json:"optTransparent"`
	PartyMembers       []int                      `json:"partyMembers"`
	Ship               Vehicle                    `json:"ship"`
	SkillTypes         []string                   `json:"skillTypes"`
	Sounds             []Bgm                      `json:"sounds"`
	StartMapId         int                        `json:"startMapId"`
	StartX             int                        `json:"startX"`
	StartY             int                        `json:"startY"`
	Switches           []string                   `json:"switches"`
	Terms              Terms                      `json:"terms"`
	TestBattlers       []TestBattler              `json:"testBattlers"`
	TestTroopId        int                        `json:"testTroopId"`
	Title1Name         string                     `json:"title1Name"`
	Title2Name         string                     `json:"title2Name"`
	TitleBgm           Bgm                        `json:"titleBgm"`
	TitleCommandWindow TitleCommandWindow         `json:"titleCommandWindow"`
	Variables          []string                   `json:"variables"`
	VersionId          int                        `json:"versionId"`
	VictoryMe          Bgm                        `json:"victoryMe"`
	WeaponTypes        []string                   `json:"weaponTypes"`
	WindowTone         []int                      `json:"windowTone"`
	TileSize           int                        `json:"tileSize"`
	HasEncryptedImages bool                       `json:"hasEncryptedImages"`
	HasEncryptedAudio  bool                       `json:"hasEncryptedAudio"`
	EncryptionKey      string                     `json:"encryptionKey"`
	Extras             map[string]json.RawMessage `json:"-"`
}

func (s *System) UnmarshalJSON(data []byte) error {
	type Alias System
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(s))
	return nil
}

func (s System) MarshalJSON() ([]byte, error) {
	type Alias System
	return util.MarshalWithExtras((*Alias)(&s), s.Extras)
}
