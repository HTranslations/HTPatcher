package patcher

import (
	"encoding/json"
	"htpatcher/internal/domain"
	"htpatcher/internal/domain/rpgmaker"
	"htpatcher/internal/util"
)

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
		if name, ok := util.DictLookup(dict, km, "actor_name", actor.Name); ok {
			actor.Name = name
		}
		if nickname, ok := util.DictLookup(dict, km, "actor_nickname", actor.Nickname); ok {
			actor.Nickname = nickname
		}
		if profile, ok := util.DictLookup(dict, km, "actor_profile", actor.Profile); ok {
			actor.Profile = util.Wrap(util.NoNewline(profile), patchInfo.Config.WrapWidth)
		}
		if note, ok := util.DictLookup(dict, km, "actor_note", actor.Note); ok {
			actor.Note = note
		}
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
		if name, ok := util.DictLookup(dict, km, "armor_name", armor.Name); ok {
			armor.Name = name
		}
		if description, ok := util.DictLookup(dict, km, "armor_description", armor.Description); ok {
			armor.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
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
		if name, ok := util.DictLookup(dict, km, "class_name", class.Name); ok {
			class.Name = name
		}
		if note, ok := util.DictLookup(dict, km, "class_name", class.Note); ok {
			class.Note = note
		}
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
		if name, ok := util.DictLookup(dict, km, "common_event_name", commonEvent.Name); ok {
			commonEvent.Name = name
		}
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
		if name, ok := util.DictLookup(dict, km, "enemy_name", enemy.Name); ok {
			enemy.Name = name
		}
		if note, ok := util.DictLookup(dict, km, "enemy_note", enemy.Note); ok {
			enemy.Note = note
		}
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
		if name, ok := util.DictLookup(dict, km, "item_name", item.Name); ok {
			item.Name = name
		}
		if description, ok := util.DictLookup(dict, km, "item_description", item.Description); ok {
			item.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
		if note, ok := util.DictLookup(dict, km, "item_note", item.Note); ok {
			item.Note = util.NoNewline(note)
		}
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

	if displayName, ok := util.DictLookup(dict, km, "map_display_name", mapData.DisplayName); ok {
		mapData.DisplayName = displayName
	}

	for _, event := range mapData.Events {
		if event == nil {
			continue
		}
		if name, ok := util.DictLookup(dict, km, "event_name", event.Name); ok {
			event.Name = name
		}
		if note, ok := util.DictLookup(dict, km, "event_note", event.Note); ok {
			event.Note = note
		}
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
		if name, ok := util.DictLookup(dict, km, "skill_name", skill.Name); ok {
			skill.Name = name
		}
		if description, ok := util.DictLookup(dict, km, "skill_description", skill.Description); ok {
			skill.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
		if message1, ok := util.DictLookup(dict, km, "skill_message1", skill.Message1); ok {
			skill.Message1 = util.Wrap(util.NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		if message2, ok := util.DictLookup(dict, km, "skill_message2", skill.Message2); ok {
			skill.Message2 = util.Wrap(util.NoNewline(message2), patchInfo.Config.WrapWidth)
		}
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
		if name, ok := util.DictLookup(dict, km, "state_name", state.Name); ok {
			state.Name = name
		}
		if message1, ok := util.DictLookup(dict, km, "state_message1", state.Message1); ok {
			state.Message1 = util.Wrap(util.NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		if message2, ok := util.DictLookup(dict, km, "state_message2", state.Message2); ok {
			state.Message2 = util.Wrap(util.NoNewline(message2), patchInfo.Config.WrapWidth)
		}
		if message3, ok := util.DictLookup(dict, km, "state_message3", state.Message3); ok {
			state.Message3 = util.Wrap(util.NoNewline(message3), patchInfo.Config.WrapWidth)
		}
		if message4, ok := util.DictLookup(dict, km, "state_message4", state.Message4); ok {
			state.Message4 = util.Wrap(util.NoNewline(message4), patchInfo.Config.WrapWidth)
		}
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
	if translation, ok := util.DictLookup(dict, km, "game_title", system.GameTitle); ok {
		system.GameTitle = translation
	}

	// Patch currency unit
	if translation, ok := util.DictLookup(dict, km, "system_term", system.CurrencyUnit); ok {
		system.CurrencyUnit = translation
	}

	// Patch armor types
	for i := range system.ArmorTypes {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.ArmorTypes[i]); ok {
			system.ArmorTypes[i] = translation
		}
	}

	// Patch elements
	for i := range system.Elements {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.Elements[i]); ok {
			system.Elements[i] = translation
		}
	}

	// Patch equip types
	for i := range system.EquipTypes {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.EquipTypes[i]); ok {
			system.EquipTypes[i] = translation
		}
	}

	// Patch skill types
	for i := range system.SkillTypes {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.SkillTypes[i]); ok {
			system.SkillTypes[i] = translation
		}
	}

	// Patch weapon types
	for i := range system.WeaponTypes {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.WeaponTypes[i]); ok {
			system.WeaponTypes[i] = translation
		}
	}

	// Patch switches
	for i := range system.Switches {
		if translation, ok := util.DictLookup(dict, km, "system_switch", system.Switches[i]); ok {
			system.Switches[i] = translation
		}
	}

	// Patch variables
	for i := range system.Variables {
		if translation, ok := util.DictLookup(dict, km, "system_variable", system.Variables[i]); ok {
			system.Variables[i] = translation
		}
	}

	// Patch terms basic
	for i := range system.Terms.Basic {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.Terms.Basic[i]); ok {
			system.Terms.Basic[i] = translation
		}
	}

	// Patch terms commands
	for i := range system.Terms.Commands {
		if system.Terms.Commands[i] != nil {
			if translation, ok := util.DictLookup(dict, km, "system_term", *system.Terms.Commands[i]); ok {
				system.Terms.Commands[i] = &translation
			}
		}
	}

	// Patch terms params
	for i := range system.Terms.Params {
		if translation, ok := util.DictLookup(dict, km, "system_term", system.Terms.Params[i]); ok {
			system.Terms.Params[i] = translation
		}
	}

	// Patch all term messages
	patchTermMessage(&system.Terms.Messages.AlwaysDash, dict, km)
	patchTermMessage(&system.Terms.Messages.CommandRemember, dict, km)
	patchTermMessage(&system.Terms.Messages.TouchUI, dict, km)
	patchTermMessage(&system.Terms.Messages.BgmVolume, dict, km)
	patchTermMessage(&system.Terms.Messages.BgsVolume, dict, km)
	patchTermMessage(&system.Terms.Messages.MeVolume, dict, km)
	patchTermMessage(&system.Terms.Messages.SeVolume, dict, km)
	patchTermMessage(&system.Terms.Messages.Possession, dict, km)
	patchTermMessage(&system.Terms.Messages.ExpTotal, dict, km)
	patchTermMessage(&system.Terms.Messages.ExpNext, dict, km)
	patchTermMessage(&system.Terms.Messages.SaveMessage, dict, km)
	patchTermMessage(&system.Terms.Messages.LoadMessage, dict, km)
	patchTermMessage(&system.Terms.Messages.File, dict, km)
	patchTermMessage(&system.Terms.Messages.Autosave, dict, km)
	patchTermMessage(&system.Terms.Messages.PartyName, dict, km)
	patchTermMessage(&system.Terms.Messages.Emerge, dict, km)
	patchTermMessage(&system.Terms.Messages.Preemptive, dict, km)
	patchTermMessage(&system.Terms.Messages.Surprise, dict, km)
	patchTermMessage(&system.Terms.Messages.EscapeStart, dict, km)
	patchTermMessage(&system.Terms.Messages.EscapeFailure, dict, km)
	patchTermMessage(&system.Terms.Messages.Victory, dict, km)
	patchTermMessage(&system.Terms.Messages.Defeat, dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainExp, dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainGold, dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainItem, dict, km)
	patchTermMessage(&system.Terms.Messages.LevelUp, dict, km)
	patchTermMessage(&system.Terms.Messages.ObtainSkill, dict, km)
	patchTermMessage(&system.Terms.Messages.UseItem, dict, km)
	patchTermMessage(&system.Terms.Messages.CriticalToEnemy, dict, km)
	patchTermMessage(&system.Terms.Messages.CriticalToActor, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorDamage, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorRecovery, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorGain, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorLoss, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorDrain, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorNoDamage, dict, km)
	patchTermMessage(&system.Terms.Messages.ActorNoHit, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyDamage, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyRecovery, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyGain, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyLoss, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyDrain, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyNoDamage, dict, km)
	patchTermMessage(&system.Terms.Messages.EnemyNoHit, dict, km)
	patchTermMessage(&system.Terms.Messages.Evasion, dict, km)
	patchTermMessage(&system.Terms.Messages.MagicEvasion, dict, km)
	patchTermMessage(&system.Terms.Messages.MagicReflection, dict, km)
	patchTermMessage(&system.Terms.Messages.CounterAttack, dict, km)
	patchTermMessage(&system.Terms.Messages.Substitute, dict, km)
	patchTermMessage(&system.Terms.Messages.BuffAdd, dict, km)
	patchTermMessage(&system.Terms.Messages.DebuffAdd, dict, km)
	patchTermMessage(&system.Terms.Messages.BuffRemove, dict, km)
	patchTermMessage(&system.Terms.Messages.ActionFailure, dict, km)

	return json.Marshal(system)
}

// patchTermMessage is a helper to patch a single term message
func patchTermMessage(message *string, dictionary map[string]string, keyMode string) {
	if translation, ok := util.DictLookup(dictionary, keyMode, "system_message", *message); ok {
		*message = translation
	}
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
		if name, ok := util.DictLookup(dict, km, "troop_name", troop.Name); ok {
			troop.Name = name
		}
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
		if name, ok := util.DictLookup(dict, km, "weapon_name", weapon.Name); ok {
			weapon.Name = name
		}
		if description, ok := util.DictLookup(dict, km, "weapon_description", weapon.Description); ok {
			weapon.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
	}

	return json.Marshal(weapons)
}
