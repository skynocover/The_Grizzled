package game

import (
	//"encoding/json"
	"github.com/looplab/fsm"
	"grizzled/database"
	//"math/rand"
	//"time"
	DataStr "github.com/skynocover/GoStackQueue"
	"log"
	//"fmt"
)

type tRound struct {
	tmp          int //用來暫時紀錄抽牌的人數
	speechThreat string
	playerlist   DataStr.Queue
	Status       *fsm.FSM
}

func (this *tRound) Init() {
	for i := range Players {
		this.playerlist.Push(Players[i])
	}

	this.Status = fsm.NewFSM(
		"pending", //初始值
		fsm.Events{
			//{Name: "Start", Src: []string{"pending", "End"}, Dst: "Draw"}, //遊戲開始到抽牌前
			{Name: "Start", Src: []string{"pending"}, Dst: "Mission"}, //抽玩牌遊戲開始

			{Name: "LuckyClover", Src: []string{"Mission"}, Dst: "LuckyClover"},    //抽玩牌遊戲開始
			{Name: "LuckyCloverEnd", Src: []string{"LuckyClover"}, Dst: "Mission"}, //抽玩牌遊戲開始

			{Name: "Speech", Src: []string{"Mission"}, Dst: "Speech"},    //抽玩牌遊戲開始
			{Name: "SpeechEnd", Src: []string{"Speech"}, Dst: "Mission"}, //抽玩牌遊戲開始

			{Name: "PlayCard", Src: []string{"Mission", "Play"}, Dst: "Play"},
			{Name: "Support", Src: []string{"Mission", "Play"}, Dst: "End"},
		},
		fsm.Callbacks{ //成功設置後執行
			"enter_state": func(e *fsm.Event) {
				switch e.Event {
				case "Start":
					Game.Stage = "Mission Start"
					this.tmp = 0
				case "LuckyClover":
					Game.Stage = "幸運草"
				case "LuckyCloverEnd":
					this.playerlist.Push(this.playerlist.Get())
					Game.Stage = this.playerlist.Peek().(player).Name
				case "Speech":
					Game.Stage = "演說:" + this.speechThreat
					this.playerlist.Push(this.playerlist.Get())
					//Game.Stage = this.playerlist.Peek().(player).Name
				case "SpeechEnd":
					for i := range Players{
						Players[i].status =""
					}
					this.playerlist.Push(this.playerlist.Get())
					Game.Stage = this.playerlist.Peek().(player).Name
				}
			},
		},
	)
}

func (this *tRound) Speech(player *player, threat string) {
	this.tmp = 0
	player.SpeechTime--
	this.speechThreat = threat
	this.Status.Event("Speech")
}
func (this *tRound) SpeechCard(player *player, choose int) bool {
	if player.status=="choosed" {
		return false
	}

	if player.checkHand(choose, this.speechThreat) {
		player.leaveCard(choose)
	} else { // 判斷是否明明有可以選卻沒選到
		for i := range player.Handcard {
			if player.checkHand(i, this.speechThreat) {
				return false
			}
		}
	}
	//選擇正確或沒有可以選擇才會到這裡
	player.status = "choosed"
	this.tmp++
	if this.tmp == this.playerlist.Len() {
		this.Status.Event("SpeechEnd")
	}
	return true
}

func (this *tRound) HeroUse(player *player, choose int) bool {
	if Game.checkLand(choose, player.Hero.Handle) {
		this.Status.Event("LuckyCloverEnd")
		return true
	}
	return false
}

func (this *tRound) PlayHero(playnow *player) bool {
	if Game.noManStage()[playnow.Hero.Handle] > 0 {
		this.Status.Event("LuckyClover")
		playnow.PlayHero()
		return true
	}
	return false
}

func (this *tRound) PlayCard(playernow *player, num int) {
	playernow.playCard(num)
	this.playerlist.Push(this.playerlist.Get())
	Game.Stage = this.playerlist.Peek().(player).Name
}

func (this *tRound) Draw(player *player, num int) {

	for j := 0; j < num; j++ {
		if Game.trials.cards.Empty() {
			return
		}
		player.drawCard(Game.trials.cards.Pop().(database.Card))
	}
	Game.Stage = player.Name + "已抽牌"
	this.tmp++
	log.Println(this.playerlist.Len())
	if this.tmp == this.playerlist.Len() {
		this.Status.Event("Start")
	}
}
