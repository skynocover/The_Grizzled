package game

import (
	"grizzled/database"
	"encoding/json"
	//"strings"
	//"fmt"
)

type player struct {
	Id       string
	Name     string
	Handcard []database.Card
	Hero     database.Hero
	Support  supports
	SpeechTime   int
	threat   []string
	/* Render */
	Process string
}

type supports struct {
	Left   int
	Right  int
	Left2  int
	Right2 int
}

func (this *player)Render()[]byte{
	this.Process = "hand"

	data, _ := json.Marshal(this)
	return data
}

func (this *player) InitPlayer() {
	this.Handcard = []database.Card{}
	this.Hero = database.Hero{}
	this.Support = supports{Left:1,Right:1}
	this.SpeechTime = 0
	this.threat = []string{}
}

func (this *player) takeSupport(support int) {
	if len(Players) > 3 {
		switch support {
		case 1, 2:
			this.Support.Right2++
		case 3, 4:
			this.Support.Left2++
		default:
			if support%2 == 0 {
				this.Support.Right++
			} else {
				this.Support.Left++
			}
		}
	} else {
		if support%2 == 0 {
			this.Support.Right++
		} else {
			this.Support.Left++
		}
	}
}

func (this *player) drawCard(card database.Card) {
	this.Handcard = append(this.Handcard, card)
}

func (this *player) PlayCard(choose int) {
	/*
		if len(this.handcard) == 0 || len(this.handcard) < choose+1 {
			return
		}
	*/
	if this.Handcard[choose].HardKnock == true {
		this.hardKnock(this.Handcard[choose])
	} else {
		Game.admission(this.Handcard[choose])
	}

	if this.Handcard[choose].Trap == true && !Game.trials.cards.Empty() {
		this.Handcard = append(this.Handcard, Game.trials.cards.Pop().(database.Card))
		this.PlayCard(len(this.Handcard) - 1)
	}

	for i := choose; i < len(this.Handcard)-1; i++ {
		this.Handcard[i] = this.Handcard[i+1]
	}
	this.Handcard = this.Handcard[:len(this.Handcard)-1]

	return
}

func (this *player) TakeHero(num int) {
	database.DB.Where("ID=?", num).Find(&this.Hero)
}

func (this *player) PlayHero() bool {
	if Game.noManStage()[this.Hero.Handle] > 0 {
		id := this.Hero.ID + 6
		this.Hero = database.Hero{}
		database.DB.Where("ID=?", id).Find(&this.Hero)

		Game.Stage = "幸運草"
		return true
	}
	return false
}

func (this *player) HeroPower(choose int) bool {
	if Game.checkLand(choose, this.Hero.Handle) {
		Game.Stage = "幸運草結束"
		return true
	}
	return false
}

func (this *player) hardKnock(card database.Card) {
	if card.Rain == true {
		this.threat = append(this.threat, "Rain")
	}
	if card.Snow == true {
		this.threat = append(this.threat, "Snow")
	}
	if card.Night == true {
		this.threat = append(this.threat, "Night")
	}
	if card.Bullet == true {
		this.threat = append(this.threat, "Bullet")
	}
	if card.Mask == true {
		this.threat = append(this.threat, "Mask")
	}
	if card.Whistle == true {
		this.threat = append(this.threat, "Whistle")
	}
	this.threat = append(this.threat, "HardKnock")
}

func (this *player) Speech() bool {
	if this.SpeechTime == 0 {
		return false
	}
	this.SpeechTime--

	Game.Stage = "演說"

	return true
}
