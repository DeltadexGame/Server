package scripts

import (
	"deltadex/gameplay"
	"deltadex/gameplay/events"
	"fmt"
)

func HandleMonsterDieEvent(event events.Event) {
	monster := event.EventInfo["monster"].(gameplay.Monster)
	if monster.Ability.AbilityID != 1 {
		return
	}
	player := event.EventInfo["player"].(*gameplay.Player)
	monster.MaxHealth = monster.MaxHealth / 2
	monster.Health = monster.MaxHealth
	monster.Ability.AbilityID = 0
	monster.Ability.Name = ""
	monster.Ability.Description = ""
	player.SpawnMonster(monster, event.EventInfo["position"].(int))
	fmt.Println("Spawned a monster with zombie ability")
}
