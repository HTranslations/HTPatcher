package main

import (
	"encoding/json"
)

type Bgm struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
}

type Vehicle struct {
	Bgm            Bgm    `json:"bgm"`
	CharacterIndex int    `json:"characterIndex"`
	CharacterName  string `json:"characterName"`
	StartMapId     int    `json:"startMapId"`
	StartX         int    `json:"startX"`
	StartY         int    `json:"startY"`
}

type AttackMotion struct {
	Type          int `json:"type"`
	WeaponImageId int `json:"weaponImageId"`
}

type TestBattler struct {
	ActorId int   `json:"actorId"`
	Level   int   `json:"level"`
	Equips  []int `json:"equips"`
}

type TitleCommandWindow struct {
	Background int `json:"background"`
	OffsetX    int `json:"offsetX"`
	OffsetY    int `json:"offsetY"`
}

type Advanced struct {
	GameId             int    `json:"gameId"`
	ScreenWidth        int    `json:"screenWidth"`
	ScreenHeight       int    `json:"screenHeight"`
	UiAreaWidth        int    `json:"uiAreaWidth"`
	UiAreaHeight       int    `json:"uiAreaHeight"`
	NumberFontFilename string `json:"numberFontFilename"`
	FallbackFonts      string `json:"fallbackFonts"`
	FontSize           int    `json:"fontSize"`
	MainFontFilename   string `json:"mainFontFilename"`
	WindowOpacity      int    `json:"windowOpacity"`
	ScreenScale        int    `json:"screenScale"`
}

type TermsMessages struct {
	AlwaysDash      string `json:"alwaysDash"`
	CommandRemember string `json:"commandRemember"`
	TouchUI         string `json:"touchUI"`
	BgmVolume       string `json:"bgmVolume"`
	BgsVolume       string `json:"bgsVolume"`
	MeVolume        string `json:"meVolume"`
	SeVolume        string `json:"seVolume"`
	Possession      string `json:"possession"`
	ExpTotal        string `json:"expTotal"`
	ExpNext         string `json:"expNext"`
	SaveMessage     string `json:"saveMessage"`
	LoadMessage     string `json:"loadMessage"`
	File            string `json:"file"`
	Autosave        string `json:"autosave"`
	PartyName       string `json:"partyName"`
	Emerge          string `json:"emerge"`
	Preemptive      string `json:"preemptive"`
	Surprise        string `json:"surprise"`
	EscapeStart     string `json:"escapeStart"`
	EscapeFailure   string `json:"escapeFailure"`
	Victory         string `json:"victory"`
	Defeat          string `json:"defeat"`
	ObtainExp       string `json:"obtainExp"`
	ObtainGold      string `json:"obtainGold"`
	ObtainItem      string `json:"obtainItem"`
	LevelUp         string `json:"levelUp"`
	ObtainSkill     string `json:"obtainSkill"`
	UseItem         string `json:"useItem"`
	CriticalToEnemy string `json:"criticalToEnemy"`
	CriticalToActor string `json:"criticalToActor"`
	ActorDamage     string `json:"actorDamage"`
	ActorRecovery   string `json:"actorRecovery"`
	ActorGain       string `json:"actorGain"`
	ActorLoss       string `json:"actorLoss"`
	ActorDrain      string `json:"actorDrain"`
	ActorNoDamage   string `json:"actorNoDamage"`
	ActorNoHit      string `json:"actorNoHit"`
	EnemyDamage     string `json:"enemyDamage"`
	EnemyRecovery   string `json:"enemyRecovery"`
	EnemyGain       string `json:"enemyGain"`
	EnemyLoss       string `json:"enemyLoss"`
	EnemyDrain      string `json:"enemyDrain"`
	EnemyNoDamage   string `json:"enemyNoDamage"`
	EnemyNoHit      string `json:"enemyNoHit"`
	Evasion         string `json:"evasion"`
	MagicEvasion    string `json:"magicEvasion"`
	MagicReflection string `json:"magicReflection"`
	CounterAttack   string `json:"counterAttack"`
	Substitute      string `json:"substitute"`
	BuffAdd         string `json:"buffAdd"`
	DebuffAdd       string `json:"debuffAdd"`
	BuffRemove      string `json:"buffRemove"`
	ActionFailure   string `json:"actionFailure"`
}

type Terms struct {
	Basic    []string      `json:"basic"`
	Commands []*string     `json:"commands"`
	Params   []string      `json:"params"`
	Messages TermsMessages `json:"messages"`
}

type System struct {
	Advanced           Advanced           `json:"advanced"`
	Airship            Vehicle            `json:"airship"`
	ArmorTypes         []string           `json:"armorTypes"`
	AttackMotions      []AttackMotion     `json:"attackMotions"`
	BattleBgm          Bgm                `json:"battleBgm"`
	Battleback1Name    string             `json:"battleback1Name"`
	Battleback2Name    string             `json:"battleback2Name"`
	BattlerHue         int                `json:"battlerHue"`
	BattlerName        string             `json:"battlerName"`
	BattleSystem       int                `json:"battleSystem"`
	Boat               Vehicle            `json:"boat"`
	CurrencyUnit       string             `json:"currencyUnit"`
	DefeatMe           Bgm                `json:"defeatMe"`
	EditMapId          int                `json:"editMapId"`
	Elements           []string           `json:"elements"`
	EquipTypes         []string           `json:"equipTypes"`
	GameTitle          string             `json:"gameTitle"`
	GameoverMe         Bgm                `json:"gameoverMe"`
	ItemCategories     []bool             `json:"itemCategories"`
	Locale             string             `json:"locale"`
	MagicSkills        []int              `json:"magicSkills"`
	MenuCommands       []bool             `json:"menuCommands"`
	OptAutosave        bool               `json:"optAutosave"`
	OptDisplayTp       bool               `json:"optDisplayTp"`
	OptDrawTitle       bool               `json:"optDrawTitle"`
	OptExtraExp        bool               `json:"optExtraExp"`
	OptFloorDeath      bool               `json:"optFloorDeath"`
	OptFollowers       bool               `json:"optFollowers"`
	OptKeyItemsNumber  bool               `json:"optKeyItemsNumber"`
	OptSideView        bool               `json:"optSideView"`
	OptSlipDeath       bool               `json:"optSlipDeath"`
	OptTransparent     bool               `json:"optTransparent"`
	PartyMembers       []int              `json:"partyMembers"`
	Ship               Vehicle            `json:"ship"`
	SkillTypes         []string           `json:"skillTypes"`
	Sounds             []Bgm              `json:"sounds"`
	StartMapId         int                `json:"startMapId"`
	StartX             int                `json:"startX"`
	StartY             int                `json:"startY"`
	Switches           []string           `json:"switches"`
	Terms              Terms              `json:"terms"`
	TestBattlers       []TestBattler      `json:"testBattlers"`
	TestTroopId        int                `json:"testTroopId"`
	Title1Name         string             `json:"title1Name"`
	Title2Name         string             `json:"title2Name"`
	TitleBgm           Bgm                `json:"titleBgm"`
	TitleCommandWindow TitleCommandWindow `json:"titleCommandWindow"`
	Variables          []string           `json:"variables"`
	VersionId          int                `json:"versionId"`
	VictoryMe          Bgm                `json:"victoryMe"`
	WeaponTypes        []string           `json:"weaponTypes"`
	WindowTone         []int              `json:"windowTone"`
	TileSize           int                `json:"tileSize"`
	HasEncryptedImages bool               `json:"hasEncryptedImages"`
	HasEncryptedAudio  bool               `json:"hasEncryptedAudio"`
	EncryptionKey      string             `json:"encryptionKey"`
}

func PatchSystem(data []byte, patchInfo PatchInfo) ([]byte, error) {
	var system System

	if err := json.Unmarshal(data, &system); err != nil {
		return nil, err
	}

	for i := range system.ArmorTypes {
		if translation, ok := patchInfo.Dictionary[system.ArmorTypes[i]]; ok {
			system.ArmorTypes[i] = translation
		}
	}

	for i := range system.Elements {
		if translation, ok := patchInfo.Dictionary[system.Elements[i]]; ok {
			system.Elements[i] = translation
		}
	}

	for i := range system.EquipTypes {
		if translation, ok := patchInfo.Dictionary[system.EquipTypes[i]]; ok {
			system.EquipTypes[i] = translation
		}
	}

	for i := range system.SkillTypes {
		if translation, ok := patchInfo.Dictionary[system.SkillTypes[i]]; ok {
			system.SkillTypes[i] = translation
		}
	}

	for i := range system.WeaponTypes {
		if translation, ok := patchInfo.Dictionary[system.WeaponTypes[i]]; ok {
			system.WeaponTypes[i] = translation
		}
	}

	for i := range system.Switches {
		if translation, ok := patchInfo.Dictionary[system.Switches[i]]; ok {
			system.Switches[i] = translation
		}
	}

	for i := range system.Variables {
		if translation, ok := patchInfo.Dictionary[system.Variables[i]]; ok {
			system.Variables[i] = translation
		}
	}

	for i := range system.Terms.Basic {
		if translation, ok := patchInfo.Dictionary[system.Terms.Basic[i]]; ok {
			system.Terms.Basic[i] = translation
		}
	}

	for i := range system.Terms.Commands {
		if system.Terms.Commands[i] != nil {
			if translation, ok := patchInfo.Dictionary[*system.Terms.Commands[i]]; ok {
				system.Terms.Commands[i] = &translation
			}
		}
	}

	for i := range system.Terms.Params {
		if translation, ok := patchInfo.Dictionary[system.Terms.Params[i]]; ok {
			system.Terms.Params[i] = translation
		}
	}

	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.AlwaysDash]; ok {
		system.Terms.Messages.AlwaysDash = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.CommandRemember]; ok {
		system.Terms.Messages.CommandRemember = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.TouchUI]; ok {
		system.Terms.Messages.TouchUI = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.BgmVolume]; ok {
		system.Terms.Messages.BgmVolume = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.BgsVolume]; ok {
		system.Terms.Messages.BgsVolume = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.MeVolume]; ok {
		system.Terms.Messages.MeVolume = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.SeVolume]; ok {
		system.Terms.Messages.SeVolume = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Possession]; ok {
		system.Terms.Messages.Possession = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ExpTotal]; ok {
		system.Terms.Messages.ExpTotal = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ExpNext]; ok {
		system.Terms.Messages.ExpNext = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.SaveMessage]; ok {
		system.Terms.Messages.SaveMessage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.LoadMessage]; ok {
		system.Terms.Messages.LoadMessage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.File]; ok {
		system.Terms.Messages.File = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Autosave]; ok {
		system.Terms.Messages.Autosave = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.PartyName]; ok {
		system.Terms.Messages.PartyName = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Emerge]; ok {
		system.Terms.Messages.Emerge = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Preemptive]; ok {
		system.Terms.Messages.Preemptive = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Surprise]; ok {
		system.Terms.Messages.Surprise = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EscapeStart]; ok {
		system.Terms.Messages.EscapeStart = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EscapeFailure]; ok {
		system.Terms.Messages.EscapeFailure = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Victory]; ok {
		system.Terms.Messages.Victory = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Defeat]; ok {
		system.Terms.Messages.Defeat = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ObtainExp]; ok {
		system.Terms.Messages.ObtainExp = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ObtainGold]; ok {
		system.Terms.Messages.ObtainGold = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ObtainItem]; ok {
		system.Terms.Messages.ObtainItem = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.LevelUp]; ok {
		system.Terms.Messages.LevelUp = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ObtainSkill]; ok {
		system.Terms.Messages.ObtainSkill = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.UseItem]; ok {
		system.Terms.Messages.UseItem = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.CriticalToEnemy]; ok {
		system.Terms.Messages.CriticalToEnemy = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.CriticalToActor]; ok {
		system.Terms.Messages.CriticalToActor = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorDamage]; ok {
		system.Terms.Messages.ActorDamage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorRecovery]; ok {
		system.Terms.Messages.ActorRecovery = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorGain]; ok {
		system.Terms.Messages.ActorGain = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorLoss]; ok {
		system.Terms.Messages.ActorLoss = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorDrain]; ok {
		system.Terms.Messages.ActorDrain = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorNoDamage]; ok {
		system.Terms.Messages.ActorNoDamage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActorNoHit]; ok {
		system.Terms.Messages.ActorNoHit = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyDamage]; ok {
		system.Terms.Messages.EnemyDamage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyRecovery]; ok {
		system.Terms.Messages.EnemyRecovery = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyGain]; ok {
		system.Terms.Messages.EnemyGain = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyLoss]; ok {
		system.Terms.Messages.EnemyLoss = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyDrain]; ok {
		system.Terms.Messages.EnemyDrain = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyNoDamage]; ok {
		system.Terms.Messages.EnemyNoDamage = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.EnemyNoHit]; ok {
		system.Terms.Messages.EnemyNoHit = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Evasion]; ok {
		system.Terms.Messages.Evasion = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.MagicEvasion]; ok {
		system.Terms.Messages.MagicEvasion = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.MagicReflection]; ok {
		system.Terms.Messages.MagicReflection = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.CounterAttack]; ok {
		system.Terms.Messages.CounterAttack = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.Substitute]; ok {
		system.Terms.Messages.Substitute = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.BuffAdd]; ok {
		system.Terms.Messages.BuffAdd = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.DebuffAdd]; ok {
		system.Terms.Messages.DebuffAdd = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.BuffRemove]; ok {
		system.Terms.Messages.BuffRemove = translation
	}
	if translation, ok := patchInfo.Dictionary[system.Terms.Messages.ActionFailure]; ok {
		system.Terms.Messages.ActionFailure = translation
	}

	mergedData, err := MergeJsonChanges(data, system)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
