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

	for _, actor := range actors {
		if actor == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(actor.Name)]; ok {
			actor.Name = name
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

	for _, armor := range armors {
		if armor == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(armor.Name)]; ok {
			armor.Name = name
		}
		if description, ok := patchInfo.Dictionary[util.GetTranslationKey(armor.Description)]; ok {
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

	for _, class := range classes {
		if class == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(class.Name)]; ok {
			class.Name = name
		}
		if note, ok := patchInfo.Dictionary[util.GetTranslationKey(class.Note)]; ok {
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

	for _, commonEvent := range commonEvents {
		if commonEvent == nil {
			continue
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

	for _, enemy := range enemies {
		if enemy == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(enemy.Name)]; ok {
			enemy.Name = name
		}
		if note, ok := patchInfo.Dictionary[util.GetTranslationKey(enemy.Note)]; ok {
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

	for _, item := range items {
		if item == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(item.Name)]; ok {
			item.Name = name
		}
		if description, ok := patchInfo.Dictionary[util.GetTranslationKey(item.Description)]; ok {
			item.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
		if note, ok := patchInfo.Dictionary[util.GetTranslationKey(item.Note)]; ok {
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

	if displayName, ok := patchInfo.Dictionary[util.GetTranslationKey(mapData.DisplayName)]; ok {
		mapData.DisplayName = displayName
	}

	for _, event := range mapData.Events {
		if event == nil {
			continue
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

	for _, skill := range skills {
		if skill == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(skill.Name)]; ok {
			skill.Name = name
		}
		if description, ok := patchInfo.Dictionary[util.GetTranslationKey(skill.Description)]; ok {
			skill.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
		if message1, ok := patchInfo.Dictionary[util.GetTranslationKey(skill.Message1)]; ok {
			skill.Message1 = util.Wrap(util.NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		if message2, ok := patchInfo.Dictionary[util.GetTranslationKey(skill.Message2)]; ok {
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

	for _, state := range states {
		if state == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(state.Name)]; ok {
			state.Name = name
		}
		if message1, ok := patchInfo.Dictionary[util.GetTranslationKey(state.Message1)]; ok {
			state.Message1 = util.Wrap(util.NoNewline(message1), patchInfo.Config.WrapWidth)
		}
		if message2, ok := patchInfo.Dictionary[util.GetTranslationKey(state.Message2)]; ok {
			state.Message2 = util.Wrap(util.NoNewline(message2), patchInfo.Config.WrapWidth)
		}
		if message3, ok := patchInfo.Dictionary[util.GetTranslationKey(state.Message3)]; ok {
			state.Message3 = util.Wrap(util.NoNewline(message3), patchInfo.Config.WrapWidth)
		}
		if message4, ok := patchInfo.Dictionary[util.GetTranslationKey(state.Message4)]; ok {
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

	// Patch armor types
	for i := range system.ArmorTypes {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.ArmorTypes[i])]; ok {
			system.ArmorTypes[i] = translation
		}
	}

	// Patch elements
	for i := range system.Elements {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.Elements[i])]; ok {
			system.Elements[i] = translation
		}
	}

	// Patch equip types
	for i := range system.EquipTypes {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.EquipTypes[i])]; ok {
			system.EquipTypes[i] = translation
		}
	}

	// Patch skill types
	for i := range system.SkillTypes {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.SkillTypes[i])]; ok {
			system.SkillTypes[i] = translation
		}
	}

	// Patch weapon types
	for i := range system.WeaponTypes {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.WeaponTypes[i])]; ok {
			system.WeaponTypes[i] = translation
		}
	}

	// Patch switches
	for i := range system.Switches {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.Switches[i])]; ok {
			system.Switches[i] = translation
		}
	}

	// Patch variables
	for i := range system.Variables {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.Variables[i])]; ok {
			system.Variables[i] = translation
		}
	}

	// Patch terms basic
	for i := range system.Terms.Basic {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.Terms.Basic[i])]; ok {
			system.Terms.Basic[i] = translation
		}
	}

	// Patch terms commands
	for i := range system.Terms.Commands {
		if system.Terms.Commands[i] != nil {
			if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(*system.Terms.Commands[i])]; ok {
				system.Terms.Commands[i] = &translation
			}
		}
	}

	// Patch terms params
	for i := range system.Terms.Params {
		if translation, ok := patchInfo.Dictionary[util.GetTranslationKey(system.Terms.Params[i])]; ok {
			system.Terms.Params[i] = translation
		}
	}

	// Patch all term messages
	patchTermMessage(&system.Terms.Messages.AlwaysDash, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.CommandRemember, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.TouchUI, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.BgmVolume, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.BgsVolume, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.MeVolume, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.SeVolume, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Possession, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ExpTotal, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ExpNext, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.SaveMessage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.LoadMessage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.File, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Autosave, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.PartyName, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Emerge, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Preemptive, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Surprise, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EscapeStart, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EscapeFailure, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Victory, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Defeat, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ObtainExp, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ObtainGold, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ObtainItem, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.LevelUp, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ObtainSkill, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.UseItem, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.CriticalToEnemy, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.CriticalToActor, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorDamage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorRecovery, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorGain, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorLoss, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorDrain, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorNoDamage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActorNoHit, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyDamage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyRecovery, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyGain, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyLoss, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyDrain, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyNoDamage, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.EnemyNoHit, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Evasion, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.MagicEvasion, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.MagicReflection, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.CounterAttack, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.Substitute, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.BuffAdd, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.DebuffAdd, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.BuffRemove, patchInfo.Dictionary)
	patchTermMessage(&system.Terms.Messages.ActionFailure, patchInfo.Dictionary)

	return json.Marshal(system)
}

// patchTermMessage is a helper to patch a single term message
func patchTermMessage(message *string, dictionary map[string]string) {
	if translation, ok := dictionary[util.GetTranslationKey(*message)]; ok {
		*message = translation
	}
}

// patchTroops patches troop data
func patchTroops(data []byte, patchInfo *domain.PatchInfo) ([]byte, error) {
	var troops rpgmaker.TroopsData
	if err := json.Unmarshal(data, &troops); err != nil {
		return nil, err
	}

	for _, troop := range troops {
		if troop == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(troop.Name)]; ok {
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

	for _, weapon := range weapons {
		if weapon == nil {
			continue
		}
		if name, ok := patchInfo.Dictionary[util.GetTranslationKey(weapon.Name)]; ok {
			weapon.Name = name
		}
		if description, ok := patchInfo.Dictionary[util.GetTranslationKey(weapon.Description)]; ok {
			weapon.Description = util.Wrap(util.NoNewline(description), patchInfo.Config.WrapWidth)
		}
	}

	return json.Marshal(weapons)
}
