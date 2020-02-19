package gameplay

import "fmt"

func HandleZombieAbility(info map[string]interface{}) (map[string]interface{}, bool) {
	monster := info["monster"].(Monster)
	if monster.Ability.AbilityID != 1 {
		return info, false
	}
	player := info["player"].(*Player)
	monster.MaxHealth = monster.MaxHealth / 2
	monster.Health = monster.MaxHealth
	monster.Ability.AbilityID = 0
	monster.Ability.Name = ""
	monster.Ability.Description = ""
	player.SpawnMonster(monster, info["position"].(int))
	fmt.Println("Spawned a monster with zombie ability")
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
