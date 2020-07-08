package game

import (
	"encoding/json"
	"grizzled/database"
	"os"
)

type player struct {
	Id       string
	Name     string
	Handcard []database.Card
	Hero     database.Hero
	Supports
	SpeechTime int
	threat     []string
	WithDraw   bool
	status     string //判斷是否在speech有選擇卡片
	/* Render */
	Process string
}

type Supports struct {
	Support [4]int //左傳,右傳,左傳2,右傳2
}

func (this *player) Render() []byte {
	this.Process = "hand"

	data, _ := json.Marshal(this)
	return data
}

func (this *player) InitPlayer(hero int, support int) {
	this.Handcard = []database.Card{}
	this.Hero = database.Hero{}
	this.Support = [4]int{1, 1, 0, 0}
	this.SpeechTime = 0
	this.threat = []string{}
	/* take support*/
	if len(Players) > 3 {
		switch support {
		case 1, 2:
			this.Support[2]++
		case 3, 4:
			this.Support[3]++
		default:
			if support%2 == 0 {
				this.Support[1]++
			} else {
				this.Support[0]++
			}
		}
	} else {
		if support%2 == 0 {
			this.Support[1]++
		} else {
			this.Support[0]++
		}
	}
	/*take hero*/
	database.DB.Where("ID=?", hero).Find(&this.Hero)
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

	if this.Handcard[choose].Trap == true && !Game.trials.cards.Empty() && os.Getenv("trap") != "0" {
		this.Handcard = append(this.Handcard, Game.trials.cards.Pop().(database.Card))
		this.playCard(len(this.Handcard) - 1)
	}

	for i := choose; i < len(this.Handcard)-1; i++ {
		this.Handcard[i] = this.Handcard[i+1]
	}
	this.Handcard = this.Handcard[:len(this.Handcard)-1]

	return
}

func (this *player) PlayHero() {
	id := this.Hero.ID + 6
	this.Hero = database.Hero{}
	database.DB.Where("ID=?", id).Find(&this.Hero)
}

func (this *player) hardKnock(card database.Card) {
	if card.Rain == true {
		this.threat = append(this.threat, "Rain")
	} else if card.Snow == true {
		this.threat = append(this.threat, "Snow")
	} else if card.Night == true {
		this.threat = append(this.threat, "Night")
	} else if card.Bullet == true {
		this.threat = append(this.threat, "Bullet")
	} else if card.Mask == true {
		this.threat = append(this.threat, "Mask")
	} else if card.Whistle == true {
		this.threat = append(this.threat, "Whistle")
	} else {
		this.threat = append(this.threat, "HardKnock")
		this.threat = append(this.threat, "HardKnock")
	}
}

func (this *player) checkHand(choose int, handle string) bool { //確認演說的是否和手牌相符

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
