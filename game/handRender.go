package game

import (
	//"encoding/json"
)

type tOnhand struct {
	Process string
	Hands   []string
	Hero    string
	Speech  int
}
/*
func (this *tOnhand) Render(player int) []byte {
	this.Process = "hand"
	this.Hands = []string{}

	for i := 0; i < 8; i++ {
		if i < len(Players[player].handcard) {
			this.Hands = append(this.Hands, Players[player].handcard[i].Name)
		} else {
			this.Hands = append(this.Hands, "cardBack")
		}
	}

	this.Hero = Players[player].Hero.Name
	this.Speech = Players[player].Speech

	data, _ := json.Marshal(this)
	return data
}
*/