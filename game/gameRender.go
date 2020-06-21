package game

import (
	"encoding/json"
)

type tBoard struct {
	Process     string
	NoMansLands []string
	Players     []string
	Threats     [][]string
	Stage       string
	PlayerNow   string
}

func (this *tBoard) Render() []byte {
	this.Process = "game"

	this.NoMansLands = []string{}
	this.Threats = [][]string{}
	this.Players = []string{}

	for i := 0; i < 7; i++ {
		if i < len(Game.NoMansLand) {
			this.NoMansLands = append(this.NoMansLands, Game.NoMansLand[i].Name)
		} else {
			this.NoMansLands = append(this.NoMansLands, "")
		}
	}

	for i := 0; i < len(Players); i++ {
		this.Players = append(this.Players, Players[i].Name)
		this.Threats = append(this.Threats, []string{})
		for j := 0; j < len(Players[i].threat); j++ {
			this.Threats[i] = append(this.Threats[i], Players[i].threat[j])
		}
	}
	this.Stage = Game.Stage
	this.PlayerNow = Players[Game.order].Name

	data, _ := json.Marshal(this)
	return data
}
