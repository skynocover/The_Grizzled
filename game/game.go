package game

import (
	"encoding/json"
	"grizzled/database"
	"math/rand"
	"os"
	"strconv"
	"time"

	DataStr "github.com/skynocover/GoStackQueue"
	//"fmt"
)

var (
	Round   tRound
	Game    Tgame
	Players []player
)

type pile struct {
	cards DataStr.Stack
}

type Tgame struct {
	NoMansLand []database.Card
	trials     pile
	morale     pile
	order      int
	Stage      string
	speech     int
	/* render用*/
	Process  string
	Players  []string
	Threats  [][]string
	RoundWin bool
	//PlayerNow string
	TM [2]int //兩個牌組的數量
}

func (this *Tgame) Render() []byte {
	this.Process = "game"

	this.Threats = [][]string{}
	this.Players = []string{}

	for i := 0; i < len(Players); i++ {
		if Players[i].WithDraw {
			this.Players = append(this.Players, Players[i].Name+"已撤退")
		} else {
			this.Players = append(this.Players, Players[i].Name)
		}

		this.Threats = append(this.Threats, []string{})
		for j := 0; j < len(Players[i].threat); j++ {
			this.Threats[i] = append(this.Threats[i], Players[i].threat[j])
		}
	}

	this.TM = [2]int{this.trials.cards.Len(), this.morale.cards.Len()}

	data, _ := json.Marshal(this)
	return data
}

func (this *Tgame) noManStage() map[string]int { //回傳每種威脅數量的map
	stage := map[string]int{
		"Rain": 0, "Night": 0, "Snow": 0, "Bullet": 0, "Mask": 0, "Whistle": 0,
	}

	for _, land := range this.NoMansLand {
		if land.Bullet {
			stage["Bullet"] = stage["Bullet"] + 1
		}
		if land.Rain {
			stage["Rain"] = stage["Rain"] + 1
		}
		if land.Night {
			stage["Night"] = stage["Night"] + 1
		}
		if land.Snow {
			stage["Snow"] = stage["Snow"] + 1
		}
		if land.Mask {
			stage["Mask"] = stage["Mask"] + 1
		}
		if land.Whistle {
			stage["Whistle"] = stage["Whistle"] + 1
		}
	}
	return stage
}

func (this *Tgame) InitGame() {
	this.speech = 5
	this.trials = pile{}
	this.morale = pile{}
	this.NoMansLand = []database.Card{}

	/*  洗牌  */
	allCard := 48
	trialsCard, _ := strconv.Atoi(os.Getenv("trial"))
	newCards := randCard(allCard)

	for i := 0; i < trialsCard; i++ {
		findcard := database.Card{}
		database.DB.Where("ID=?", newCards[i]).Find(&findcard)
		this.trials.cards.Push(findcard)
	}

	for j := trialsCard; j < allCard; j++ {
		findcard := database.Card{}
		database.DB.Where("ID=?", newCards[j]).Find(&findcard)
		this.morale.cards.Push(findcard)
	}

	/*  抽英雄  */
	hero := randCard(6)
	support := randCard(16 - len(Players)*2)
	for i := range Players {
		Players[i].InitPlayer(hero[i], support[i])
	}
	Round.rounds = 0
	Round.Status.SetState("pending")
	drawNum, _ := strconv.Atoi(os.Getenv("initcard")) //設定第一回合的抽牌數
	Round.Init(drawNum)                               //開始新的回合並且抽三張
}

func (this *Tgame) NewPlayer(id string, name string) bool {
	for i := range Players {
		if Players[i].Name == name {
			return false
		}
	}
	Players = append(Players, player{Id: id, Name: name})
	return true
}

func (this *Tgame) checkLand(choose int, handle string) bool {
	if this.NoMansLand[choose].Mask && handle == "Mask" {
		this.leaveCard(choose)
		return true
	}
	if this.NoMansLand[choose].Rain && handle == "Rain" {
		this.leaveCard(choose)
		return true
	}
	if this.NoMansLand[choose].Snow && handle == "Snow" {
		this.leaveCard(choose)
		return true
	}
	if this.NoMansLand[choose].Bullet && handle == "Bullet" {
		this.leaveCard(choose)
		return true
	}
	if this.NoMansLand[choose].Night && handle == "Night" {
		this.leaveCard(choose)
		return true
	}
	if this.NoMansLand[choose].Whistle && handle == "Whistle" {
		this.leaveCard(choose)
		return true
	}
	return false
}

func (this *Tgame) admission(card database.Card) {
	this.NoMansLand = append(this.NoMansLand, card)
}

func (this *Tgame) leaveCard(choose int) {
	for i := choose; i < len(this.NoMansLand)-1; i++ {
		this.NoMansLand[i] = this.NoMansLand[i+1]
	}
	this.NoMansLand = this.NoMansLand[:len(this.NoMansLand)-1]
}

func randCard(num int) (pile []int) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(pile) < num {
		r := rand.Intn((num)) + 1

		exist := false
		for _, v := range pile {
			if v == r {
				exist = true
				break
			}
		}

		if !exist {
			pile = append(pile, r)
		}
	}
	return
}

func GetOrder(id string) int { //回傳真正的玩家順序
	for i, p := range Players {
		if p.Id == id {
			return i
		}
	}
	return -1
}
