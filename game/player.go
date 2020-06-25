package game

import (
	"encoding/json"
	"grizzled/database"
	//"strings"
	//"fmt"
)

type player struct {
	Id         string
	Name       string
	Handcard   []database.Card
	Hero       database.Hero
	Support    supports
	SpeechTime int
	threat     []string
	status string
	/* Render */
	Process string
}

type supports struct {
	Left   int
	Right  int
	Left2  int
	Right2 int
}

func (this *player) Render() []byte {
	this.Process = "hand"

	data, _ := json.Marshal(this)
	return data
}

func (this *player) InitPlayer() {
	this.Handcard = []database.Card{}
	this.Hero = database.Hero{}
	this.Support = supports{Left: 1, Right: 1}
	this.SpeechTime = 0
	this.takeSpeech()
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

func (this *player) takeSpeech() {
	if Game.speech > 0 {
		Game.speech--
		this.SpeechTime++
	}
}

func (this *player) drawCard(card database.Card) {
	this.Handcard = append(this.Handcard, card)
}

func (this *player) playCard(choose int) {

	if this.Handcard[choose].HardKnock == true {
		this.hardKnock(this.Handcard[choose])
	} else {
		Game.admission(this.Handcard[choose])
	}

	if this.Handcard[choose].Trap == true && !Game.trials.cards.Empty() {
		this.Handcard = append(this.Handcard, Game.trials.cards.Pop().(database.Card))
		this.playCard(len(this.Handcard) - 1)
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

func (this *player) PlayHero() {
	id := this.Hero.ID + 6
	this.Hero = database.Hero{}
	database.DB.Where("ID=?", id).Find(&this.Hero)
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

func (this *player) checkHand(choose int, handle string) bool {

	switch handle {
	case "Mask":
		if this.Handcard[choose].Mask {
			return true
		}
	case "Rain":
		if this.Handcard[choose].Rain {
			return true
		}
	case "Snow":
		if this.Handcard[choose].Snow {
			return true
		}
	case "Bullet":
		if this.Handcard[choose].Bullet {
			return true
		}
	case "Night":
		if this.Handcard[choose].Night {
			return true
		}
	case "Whistle":
		if this.Handcard[choose].Whistle {
			return true
		}
	}
	return false
}

func (this *player) leaveCard(choose int) {
	for i := choose; i < len(this.Handcard)-1; i++ {
		this.Handcard[i] = this.Handcard[i+1]
	}
	this.Handcard = this.Handcard[:len(this.Handcard)-1]
}
