package patcher

import (
	"encoding/json"
	"fmt"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
)

func patchTextField(value *string, extras *map[string]json.RawMessage, fieldPath string, dict map[string]string, keyMode string, entryType string, transform func(string) string) {
	if translation, ok := util.DictLookup(dict, keyMode, entryType, *value); ok {
		setHTOriginal(extras, fieldPath, *value)
		if transform != nil {
			translation = transform(translation)
		}
		*value = translation
	}
}

func wrapNoNewline(width int) func(string) string {
	return func(text string) string {
		return util.Wrap(util.NoNewline(text), width)
	}
}

// patchActors patches actor data
func patchActors(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var actors rpgmaker.ActorsData
	if err := json.Unmarshal(data, &actors); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, actor := range actors {
		if actor == nil {
			continue
		}
		patchTextField(&actor.Name, &actor.Extras, "name", dict, km, "actor_name", nil)
		patchTextField(&actor.Nickname, &actor.Extras, "nickname", dict, km, "actor_nickname", nil)
		patchTextField(&actor.Profile, &actor.Extras, "profile", dict, km, "actor_profile", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&actor.Note, &actor.Extras, "note", dict, km, "actor_note", nil)
	}

	return json.Marshal(actors)
}

// patchArmors patches armor data
func patchArmors(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var armors rpgmaker.ArmorsData
	if err := json.Unmarshal(data, &armors); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, armor := range armors {
		if armor == nil {
			continue
		}
		patchTextField(&armor.Name, &armor.Extras, "name", dict, km, "armor_name", nil)
		patchTextField(&armor.Description, &armor.Extras, "description", dict, km, "armor_description", wrapNoNewline(patchInfo.Config.WrapWidth))
	}

	return json.Marshal(armors)
}

// patchClasses patches class data
func patchClasses(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var classes rpgmaker.ClassesData
	if err := json.Unmarshal(data, &classes); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, class := range classes {
		if class == nil {
			continue
		}
		patchTextField(&class.Name, &class.Extras, "name", dict, km, "class_name", nil)
		patchTextField(&class.Note, &class.Extras, "note", dict, km, "class_name", nil)
	}

	return json.Marshal(classes)
}

// patchCommonEvents patches common event data
func patchCommonEvents(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var commonEvents rpgmaker.CommonEventsData
	if err := json.Unmarshal(data, &commonEvents); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, commonEvent := range commonEvents {
		if commonEvent == nil {
			continue
		}
		patchTextField(&commonEvent.Name, &commonEvent.Extras, "name", dict, km, "common_event_name", nil)
		newCommands, err := patchCommands(commonEvent.List, patchInfo)
		if err != nil {
			return nil, err
		}
		commonEvent.List = newCommands
	}

	return json.Marshal(commonEvents)
}

// patchEnemies patches enemy data
func patchEnemies(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var enemies rpgmaker.EnemiesData
	if err := json.Unmarshal(data, &enemies); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, enemy := range enemies {
		if enemy == nil {
			continue
		}
		patchTextField(&enemy.Name, &enemy.Extras, "name", dict, km, "enemy_name", nil)
		patchTextField(&enemy.Note, &enemy.Extras, "note", dict, km, "enemy_note", nil)
	}

	return json.Marshal(enemies)
}

// patchItems patches item data
func patchItems(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var items rpgmaker.ItemsData
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, item := range items {
		if item == nil {
			continue
		}
		patchTextField(&item.Name, &item.Extras, "name", dict, km, "item_name", nil)
		patchTextField(&item.Description, &item.Extras, "description", dict, km, "item_description", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&item.Note, &item.Extras, "note", dict, km, "item_note", util.NoNewline)
	}

	return json.Marshal(items)
}

// patchMap patches map data
func patchMap(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var mapData rpgmaker.MapData
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	patchTextField(&mapData.DisplayName, &mapData.Extras, "displayName", dict, km, "map_display_name", nil)

	for _, event := range mapData.Events {
		if event == nil {
			continue
		}
		patchTextField(&event.Name, &event.Extras, "name", dict, km, "event_name", nil)
		patchTextField(&event.Note, &event.Extras, "note", dict, km, "event_note", nil)
		for i := range event.Pages {
			newCommands, err := patchCommands(event.Pages[i].List, patchInfo)
			if err != nil {
				return nil, err
			}
			event.Pages[i].List = newCommands
		}
	}

	return json.Marshal(mapData)
}

// patchMapInfos patches editor map names from MapInfos.json.
func patchMapInfos(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var mapInfos []*rpgmaker.MapInfo
	if err := json.Unmarshal(data, &mapInfos); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, info := range mapInfos {
		if info == nil {
			continue
		}
		patchTextField(&info.Name, &info.Extras, "name", dict, km, "map_name", nil)
	}

	return json.Marshal(mapInfos)
}

// patchSkills patches skill data
func patchSkills(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var skills rpgmaker.SkillsData
	if err := json.Unmarshal(data, &skills); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, skill := range skills {
		if skill == nil {
			continue
		}
		patchTextField(&skill.Name, &skill.Extras, "name", dict, km, "skill_name", nil)
		patchTextField(&skill.Description, &skill.Extras, "description", dict, km, "skill_description", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&skill.Message1, &skill.Extras, "message1", dict, km, "skill_message1", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&skill.Message2, &skill.Extras, "message2", dict, km, "skill_message2", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&skill.Note, &skill.Extras, "note", dict, km, "skill_note", nil)
	}

	return json.Marshal(skills)
}

// patchStates patches state data
func patchStates(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var states rpgmaker.StatesData
	if err := json.Unmarshal(data, &states); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, state := range states {
		if state == nil {
			continue
		}
		patchTextField(&state.Name, &state.Extras, "name", dict, km, "state_name", nil)
		patchTextField(&state.Message1, &state.Extras, "message1", dict, km, "state_message1", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&state.Message2, &state.Extras, "message2", dict, km, "state_message2", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&state.Message3, &state.Extras, "message3", dict, km, "state_message3", wrapNoNewline(patchInfo.Config.WrapWidth))
		patchTextField(&state.Message4, &state.Extras, "message4", dict, km, "state_message4", wrapNoNewline(patchInfo.Config.WrapWidth))
	}

	return json.Marshal(states)
}

// patchSystem patches system data
func patchSystem(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var system rpgmaker.System
	if err := json.Unmarshal(data, &system); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	// Set locale if specified in patch config
	if patchInfo.Config != nil && patchInfo.Config.Locale != "" {
		system.Locale = patchInfo.Config.Locale
	}

	// Patch game title
	patchTextField(&system.GameTitle, &system.Extras, "gameTitle", dict, km, "game_title", nil)

	// Patch currency unit
	patchTextField(&system.CurrencyUnit, &system.Extras, "currencyUnit", dict, km, "system_term", nil)

	// Patch armor types
	for i := range system.ArmorTypes {
		patchTextField(&system.ArmorTypes[i], &system.Extras, fmt.Sprintf("armorTypes[%d]", i), dict, km, "system_term", nil)
	}

	// Patch elements
	for i := range system.Elements {
		patchTextField(&system.Elements[i], &system.Extras, fmt.Sprintf("elements[%d]", i), dict, km, "system_term", nil)
	}

	// Patch equip types
	for i := range system.EquipTypes {
		patchTextField(&system.EquipTypes[i], &system.Extras, fmt.Sprintf("equipTypes[%d]", i), dict, km, "system_term", nil)
	}

	// Patch skill types
	for i := range system.SkillTypes {
		patchTextField(&system.SkillTypes[i], &system.Extras, fmt.Sprintf("skillTypes[%d]", i), dict, km, "system_term", nil)
	}

	// Patch weapon types
	for i := range system.WeaponTypes {
		patchTextField(&system.WeaponTypes[i], &system.Extras, fmt.Sprintf("weaponTypes[%d]", i), dict, km, "system_term", nil)
	}

	// Patch switches
	for i := range system.Switches {
		patchTextField(&system.Switches[i], &system.Extras, fmt.Sprintf("switches[%d]", i), dict, km, "system_switch", nil)
	}

	// Patch variables
	for i := range system.Variables {
		patchTextField(&system.Variables[i], &system.Extras, fmt.Sprintf("variables[%d]", i), dict, km, "system_variable", nil)
	}

	// Patch terms basic
	for i := range system.Terms.Basic {
		patchTextField(&system.Terms.Basic[i], &system.Terms.Extras, fmt.Sprintf("basic[%d]", i), dict, km, "system_term", nil)
	}

	// Patch terms commands
	for i := range system.Terms.Commands {
		if system.Terms.Commands[i] != nil {
			patchTextField(system.Terms.Commands[i], &system.Terms.Extras, fmt.Sprintf("commands[%d]", i), dict, km, "system_term", nil)
		}
	}

	// Patch terms params
	for i := range system.Terms.Params {
		patchTextField(&system.Terms.Params[i], &system.Terms.Extras, fmt.Sprintf("params[%d]", i), dict, km, "system_term", nil)
	}

	// Patch all term messages
	patchTermMessage(&system.Terms.Messages.AlwaysDash, &system.Terms.Messages.Extras, "alwaysDash", dict, km)
	patchTermMessage(&system.Terms.Messages.CommandRemember, &system.Terms.Messages.Extras, "commandRemember", dict, km)
	patchTermMessage(&system.Terms.Messages.TouchUI, &system.Terms.Messages.Extras, "touchUI", dict, km)
	patchTermMessage(&system.Terms.Messages.BgmVolume, &system.Terms.Messages.Extras, "bgmVolume", dict, km)
	patchTermMessage(&system.Terms.Messages.BgsVolume, &system.Terms.Messages.Extras, "bgsVolume", dict, km)
	patchTermMessage(&system.Terms.Messages.MeVolume, &system.Terms.Messages.Extras, "meVolume", dict, km)
	patchTermMessage(&system.Terms.Messages.SeVolume, &system.Terms.Messages.Extras, "seVolume", dict, km)
	patchTermMessage(&system.Terms.Messages.Possession, &system.Terms.Messages.Extras, "possession", dict, km)
	patchTermMessage(&system.Terms.Messages.ExpTotal, &system.Terms.Messages.Extras, "expTotal", dict, km)
	patchTermMessage(&system.Terms.Messages.ExpNext, &system.Terms.Messages.Extras, "expNext", dict, km)
	patchTermMessage(&system.Terms.Messages.SaveMessage, &system.Terms.Messages.Extras, "saveMessage", dict, km)
	patchTermMessage(&system.Terms.Messages.LoadMessage, &system.Terms.Messages.Extras, "loadMessage", dict, km)
	patchTermMessage(&system.Terms.Messages.File, &system.Terms.Messages.Extras, "file", dict, km)
	patchTermMessage(&system.Terms.Messages.Autosave, &system.Terms.Messages.Extras, "autosave", dict, km)
	patchTermMessage(&system.Terms.Messages.PartyName, &system.Terms.Messages.Extras, "partyName", dict, km)
	patchTermMessage(&system.Terms.Messages.Emerge, &system.Terms.Messages.Extras, "emerge", dict, km)
	patchTermMessage(&system.Terms.Messages.Preemptive, &system.Terms.Messages.Extras, "preemptive", dict, km)
	patchTermMessage(&system.Terms.Messages.Surprise, &system.Terms.Messages.Extras, "surprise", dict, km)
	patchTermMessage(&system.Terms.Messages.EscapeStart, &system.Terms.Messages.Extras, "escapeStart", dict, km)
	patchTermMessage(&system.Terms.Messages.EscapeFailure, &system.Terms.Messages.Extras, "escapeFailure", dict, km)
	patchTermMessage(&system.Terms.Messages.Victory, &system.Terms.Messages.Extras, "victory", dict, km)
	patchTermMessage(&system.Terms.Messages.Defeat, &system.Terms.Messages.Extras, "defeat", dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainExp, &system.Terms.Messages.Extras, "obtainExp", dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainGold, &system.Terms.Messages.Extras, "obtainGold", dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainItem, &system.Terms.Messages.Extras, "obtainItem", dict, km)
	patchTermMessage(&system.Terms.Messages.LevelUp, &system.Terms.Messages.Extras, "levelUp", dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainSkill, &system.Terms.Messages.Extras, "obtainSkill", dict, km)
	patchTermMessage(&system.Terms.Messages.UseItem, &system.Terms.Messages.Extras, "useItem", dict, km)
	patchTermMessage(&system.Terms.Messages.CriticalToEnemy, &system.Terms.Messages.Extras, "criticalToEnemy", dict, km)
	patchTermMessage(&system.Terms.Messages.CriticalToActor, &system.Terms.Messages.Extras, "criticalToActor", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorDamage, &system.Terms.Messages.Extras, "actorDamage", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorRecovery, &system.Terms.Messages.Extras, "actorRecovery", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorGain, &system.Terms.Messages.Extras, "actorGain", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorLoss, &system.Terms.Messages.Extras, "actorLoss", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorDrain, &system.Terms.Messages.Extras, "actorDrain", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorNoDamage, &system.Terms.Messages.Extras, "actorNoDamage", dict, km)
	patchTermMessage(&system.Terms.Messages.ActorNoHit, &system.Terms.Messages.Extras, "actorNoHit", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyDamage, &system.Terms.Messages.Extras, "enemyDamage", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyRecovery, &system.Terms.Messages.Extras, "enemyRecovery", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyGain, &system.Terms.Messages.Extras, "enemyGain", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyLoss, &system.Terms.Messages.Extras, "enemyLoss", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyDrain, &system.Terms.Messages.Extras, "enemyDrain", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyNoDamage, &system.Terms.Messages.Extras, "enemyNoDamage", dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyNoHit, &system.Terms.Messages.Extras, "enemyNoHit", dict, km)
	patchTermMessage(&system.Terms.Messages.Evasion, &system.Terms.Messages.Extras, "evasion", dict, km)
	patchTermMessage(&system.Terms.Messages.MagicEvasion, &system.Terms.Messages.Extras, "magicEvasion", dict, km)
	patchTermMessage(&system.Terms.Messages.MagicReflection, &system.Terms.Messages.Extras, "magicReflection", dict, km)
	patchTermMessage(&system.Terms.Messages.CounterAttack, &system.Terms.Messages.Extras, "counterAttack", dict, km)
	patchTermMessage(&system.Terms.Messages.Substitute, &system.Terms.Messages.Extras, "substitute", dict, km)
	patchTermMessage(&system.Terms.Messages.BuffAdd, &system.Terms.Messages.Extras, "buffAdd", dict, km)
	patchTermMessage(&system.Terms.Messages.DebuffAdd, &system.Terms.Messages.Extras, "debuffAdd", dict, km)
	patchTermMessage(&system.Terms.Messages.BuffRemove, &system.Terms.Messages.Extras, "buffRemove", dict, km)
	patchTermMessage(&system.Terms.Messages.ActionFailure, &system.Terms.Messages.Extras, "actionFailure", dict, km)

	return json.Marshal(system)
}

// patchTermMessage is a helper to patch a single term message
func patchTermMessage(message *string, extras *map[string]json.RawMessage, fieldPath string, dictionary map[string]string, keyMode string) {
	patchTextField(message, extras, fieldPath, dictionary, keyMode, "system_message", nil)
}

// patchTroops patches troop data
func patchTroops(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var troops rpgmaker.TroopsData
	if err := json.Unmarshal(data, &troops); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, troop := range troops {
		if troop == nil {
			continue
		}
		patchTextField(&troop.Name, &troop.Extras, "name", dict, km, "troop_name", nil)
		for i := range troop.Pages {
			newCommands, err := patchCommands(troop.Pages[i].List, patchInfo)
			if err != nil {
				return nil, err
			}
			troop.Pages[i].List = newCommands
		}
	}

	return json.Marshal(troops)
}

// patchWeapons patches weapon data
func patchWeapons(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var weapons rpgmaker.WeaponsData
	if err := json.Unmarshal(data, &weapons); err != nil {
		return nil, err
	}

	km := patchInfo.Config.KeyMode
	dict := patchInfo.Dictionary

	for _, weapon := range weapons {
		if weapon == nil {
			continue
		}
		patchTextField(&weapon.Name, &weapon.Extras, "name", dict, km, "weapon_name", nil)
		patchTextField(&weapon.Description, &weapon.Extras, "description", dict, km, "weapon_description", wrapNoNewline(patchInfo.Config.WrapWidth))
	}

	return json.Marshal(weapons)
}
