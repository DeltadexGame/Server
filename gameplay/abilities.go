package gameplay

func HandleZombieAbility(info map[string]interface{}) (map[string]interface{}, bool) {
	monster := info["monster"].(Monster)
	if monster.Ability.AbilityID != 1 {
		return info, false
	}
	player := info["player"].(*Player)
	player.Monsters[info["position"].(int)].MaxHealth = player.Monsters[info["position"].(int)].MaxHealth / 2
	player.Monsters[info["position"].(int)].Health = player.Monsters[info["position"].(int)].MaxHealth
	player.Monsters[info["position"].(int)].Ability.AbilityID = 0
	player.Monsters[info["position"].(int)].Ability.Name = ""
	player.Monsters[info["position"].(int)].Ability.Description = ""
	info["cancelled"] = true
	return info, false
}

func HandleHeavyAbility(info map[string]interface{}) (map[string]interface{}, bool) {
	monster := info["monster"].(Monster)
	if monster.Ability.AbilityID != 2 {
		return info, false
	}
	info["damage"] = info["damage"].(int) / 2
	return info, true
}
