package game

import (
	"grizzled/database"
	"log"
	"os"
	"strconv"

	"github.com/looplab/fsm"
	DataStr "github.com/skynocover/GoStackQueue"
)

type tRound struct {
	rounds       int        //回合數
	drawNum      int        //抽牌數目
	success      bool       //任務成功
	support      []Supports //用來紀錄每個人拿到的support
	tmp          int        //用來暫時紀錄抽牌的人數
	tmpEnd       int        //tmp的最終期望值
	speechThreat string     //演說的威脅
	loseReason   string     //輸遊戲的原因
	playerlist   DataStr.Queue
	Status       *fsm.FSM
}

func init() { //宣告狀態機
	Game.NoMansLand = []database.Card{}

	Round.Status = fsm.NewFSM(
		"pending", //初始值
		fsm.Events{
			//{Name: "Start", Src: []string{"pending", "End"}, Dst: "Draw"}, //遊戲開始到抽牌前
			{Name: "Start", Src: []string{"pending", "Mission", "End", "LuckyClover", "Speech", "Support", "WinLose"}, Dst: "Mission"}, //抽玩牌遊戲開始

			{Name: "LuckyClover", Src: []string{"Mission"}, Dst: "LuckyClover"},    //幸運草階段
			{Name: "LuckyCloverEnd", Src: []string{"LuckyClover"}, Dst: "Mission"}, //幸運草階段結束

			{Name: "Speech", Src: []string{"Mission"}, Dst: "Speech"},    //演講開始
			{Name: "SpeechEnd", Src: []string{"Speech"}, Dst: "Mission"}, //演講結束

			{Name: "Support", Src: []string{"Mission"}, Dst: "Support"}, //所有人都退出後進入結算階段
			{Name: "SupportEnd", Src: []string{"Support"}, Dst: "End"},  //所有人都確認結算的結果後進入此階段

			{Name: "WinGame", Src: []string{"Mission"}, Dst: "WinLose"},
			{Name: "LoseGame", Src: []string{"Support", "Mission", "End"}, Dst: "WinLose"},
		},
		fsm.Callbacks{ //成功設置後執行
			"enter_state": func(e *fsm.Event) {
				switch e.Event {
				case "Start": //四個人都抽完牌
					Round.tmp = 0
					Game.Stage = Round.playerlist.Peek().(player).Name

				case "LuckyClover":
					log.Println("幸運草階段")
					Game.Stage = Round.playerlist.Peek().(player).Name + ":幸運草"
				case "LuckyCloverEnd":
					log.Println("幸運草結束")
					Round.playerlist.Push(Round.playerlist.Get())
					Game.Stage = Round.playerlist.Peek().(player).Name
				case "Speech":
					log.Println("演說開始")
					Game.Stage = "演說:" + Round.speechThreat
				case "SpeechEnd":
					log.Println("演說結束")
					Round.tmp = 0
					for i := range Players {
						Players[i].status = ""
					}
					Round.playerlist.Push(Round.playerlist.Get())
					Game.Stage = Round.playerlist.Peek().(player).Name
				case "Support":
					log.Println("支援結算")
					log.Println(Round.support)
					maxsupport := 0            //統計最多數目的支援
					maxsupported := []player{} //統計得到最多支援的人,使用陣列是因為可能有多個人相同
					for i := range Round.support {
						nowsupport := 0 //用來暫存此人得到的支援數
						for j := range Round.support[i].Support {
							nowsupport = Round.support[i].Support[j] + nowsupport
							Players[i].Support[j] = Players[i].Support[j] + Round.support[i].Support[j]
						}
						if nowsupport > maxsupport {
							maxsupport = nowsupport
							maxsupported = []player{}
							maxsupported = append(maxsupported, Players[i])
						} else if nowsupport == maxsupport {
							maxsupported = append(maxsupported, Players[i])
						}
					}
					if len(maxsupported) > 1 {
						Game.Stage = "SupportSkip"
						Game.RoundWin = Round.success
					} else {
						Game.Stage = "Support:" + maxsupported[0].Name
						Game.RoundWin = Round.success
					}

				case "SupportEnd":
					log.Println("支援階段結束")
					Round.tmp = 0
					//確認下一位隊長
					leader := Round.rounds % len(Players)
					Game.Stage = "Leader:" + Players[leader].Name
					if leader == 0 {
						Players[len(Players)-1].takeSpeech()
					} else {
						Players[leader-1].takeSpeech()
					}

					trial := 0 //總共要失去的士氣
					for i := range Players {
						trial = trial + len(Players[i].Handcard)
					}
					if trial < 3 {
						trial = 3
					}
					if !Round.success {
						rand := randCard(len(Game.NoMansLand))
						for i := range rand {
							Game.trials.cards.Push(Game.NoMansLand[rand[i]-1])
						}
					}
					for i := 0; i < trial; i++ {
						if Game.morale.cards.Empty() {
							log.Println("遊戲失敗")
							Round.loseReason = "士氣歸零"
							break
						}
						Game.trials.cards.Push(Game.morale.cards.Pop())
					}
				case "WinGame":
					log.Println("遊戲成功")
					Game.Stage = "Winner,Winner,Chicken Dinner"
				case "LoseGame":
					log.Println("遊戲失敗")
					Game.Stage = "Loser,Loser,now who’s dinner?:" + Round.loseReason
				}
			},
		},
	)
}

func (this *tRound) Init(num int) { //開始新的回合
	this.drawNum = num //設定本回合的抽牌數

	Game.NoMansLand = []database.Card{}
	this.loseReason = ""
	this.tmp = 0

	this.playerlist = DataStr.Queue{}
	this.support = []Supports{} //清空用來紀錄的support
	for i := range Players {
		Players[i].WithDraw = false //重設所有人的撤退
		Players[i].status = ""
		this.playerlist.Push(Players[i])
		this.support = append(this.support, Supports{}) //暫時放入空的support物件
	}

	for i := 0; i < this.rounds; i++ { //根據回合數推進玩家順序
		this.playerlist.Push(this.playerlist.Get())
	}
	this.rounds++
	Game.Stage = "DrawCard" //通知所有人抽牌
}

func (this *tRound) Support(playnow *player, choose int) bool {
	if playnow.Id != this.playerlist.Peek().(player).Id {
		return false
	}

	if playnow.Support[choose] == 0 {
		if playnow.Support[0] != 0 || playnow.Support[1] != 0 || playnow.Support[2] != 0 || playnow.Support[3] != 0 {
			return false //若選擇錯誤並且還有選擇
		}
	} else {
		playnow.Support[choose]-- //選擇正確
		// 先暫時紀錄傳遞的結果,支援階段再結算
		pOrder := GetOrder(playnow.Id)
		var tOrder int
		switch choose {
		case 0: //左傳
			if pOrder == 0 {
				tOrder = len(Players) - 1
			} else {
				tOrder = pOrder - 1
			}
			this.support[tOrder].Support[0] = this.support[tOrder].Support[0] + 1
		case 1: //右傳
			if pOrder == len(this.support)-1 {
				tOrder = 0
			} else {
				tOrder = pOrder + 1
			}
			this.support[tOrder].Support[1] = this.support[tOrder].Support[1] + 1
		case 2: //左傳2
			if pOrder == 0 || pOrder == 1 {
				tOrder = len(Players) - 2 + pOrder
			} else {
				tOrder = pOrder - 2
			}
			this.support[tOrder].Support[2] = this.support[tOrder].Support[2] + 1
		case 3: //右傳2
			if pOrder == len(this.support)-1 {
				tOrder = 0
			} else if pOrder == len(this.support)-2 {
				tOrder = 1
			} else {
				tOrder = pOrder + 2
			}
			this.support[tOrder].Support[3] = this.support[tOrder].Support[3] + 1
		}
	}
	//若沒有可以選的則會直接進入撤退
	playnow.WithDraw = true
	this.playerlist.Get()

	if this.playerlist.Len() == 0 {
		this.success = true //所有人都撤退則任務成功並進入支援階段
		this.Status.Event("Support")
	} else {
		Game.Stage = this.playerlist.Peek().(player).Name
	}
	return true
}
func (this *tRound) SupportEnd(playnow *player, choose string) {
	switch choose {
	case "Threat":
		cancel := 0
		if this.success {
			cancel = 2
		} else {
			cancel = 1
		}
		for i := range playnow.threat {
			if i >= cancel {
				break
			}
			playnow.threat = playnow.threat[1:]
		}
	case "Lucky":
		id := playnow.Hero.ID - 6
		if id > 0 {
			playnow.Hero = database.Hero{}
			database.DB.Where("ID=?", id).Find(&playnow.Hero)
		}
	}
	this.tmp++
	Game.Stage = playnow.Name + "已確認"
	if this.tmp == len(Players) {
		this.Status.Event("SupportEnd")
	}
	if this.loseReason != "" {
		this.loseGame(this.loseReason)
	}

}

func (this *tRound) Speech(playnow *player, threat string) bool {
	if playnow.Id != this.playerlist.Peek().(player).Id || playnow.SpeechTime == 0 {
		return false
	}
	this.tmp = 0
	playnow.SpeechTime--
	this.speechThreat = threat
	this.Status.Event("Speech")

	this.tmpEnd = 0 //預期會有幾個人做演講的動作
	for i := range Players {
		if len(Players[i].Handcard) != 0 {
			this.tmpEnd++
		}
	}
	return true
}
func (this *tRound) SpeechCard(player *player, choose int) bool {
	if player.status == "choosed" {
		return false
	}

	if player.checkHand(choose, this.speechThreat) { //選擇正確
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
	log.Print("tmp of speech")
	log.Println(this.tmp)

	if this.tmp == this.tmpEnd {
		this.Status.Event("SpeechEnd")
	}
	return true
}

func (this *tRound) HeroUse(player *player, choose int) bool { //英雄能力選牌
	if Game.checkLand(choose, player.Hero.Handle) {
		this.Status.Event("LuckyCloverEnd")
		return true
	}
	return false
}
func (this *tRound) PlayHero(playnow *player) bool { //使用英雄能力
	if Game.noManStage()[playnow.Hero.Handle] == 0 || playnow.Id != this.playerlist.Peek().(player).Id {
		return false
	}
	this.Status.Event("LuckyClover")
	playnow.PlayHero()
	return true
}

func (this *tRound) Draw(player *player) {
	for j := 0; j < this.drawNum; j++ {
		if Game.trials.cards.Empty() {
			break
		}
		player.drawCard(Game.trials.cards.Pop().(database.Card))
	}
	this.tmp++
	if this.tmp == len(Players) { //每個人都抽完牌後進入遊戲開始的程序
		this.Status.Event("Start")
	} else {
		Game.Stage = player.Name + "已抽牌"
	}
}

func (this *tRound) PlayCard(playernow *player, num int) bool {
	if this.playerlist.Peek().(player).Id != playernow.Id {
		return false
	} else if num >= len(playernow.Handcard) {
		return false
	}
	playernow.playCard(num)

	if this.threatOver() {
		log.Println("威脅過多")
		Round.success = false //超過威脅則任務失敗並進入支援階段
		Round.Status.Event("Support")
	} else if len(playernow.Handcard) == 0 && Game.trials.cards.Len() == 0 {
		var hand = 0
		for i := range Players {
			hand = hand + len(Players[i].Handcard)
		}
		if hand == 0 {
			this.Status.Event("WinGame")
		}
	} else {
		this.playerlist.Push(this.playerlist.Get())
		Game.Stage = this.playerlist.Peek().(player).Name
	}
	return true
}

func (this *tRound) threatOver() bool {
	threat := map[string]int{}

	for i := range Game.NoMansLand {
		if Game.NoMansLand[i].Rain {
			threat["Rain"] = threat["Rain"] + 1
		}
		if Game.NoMansLand[i].Bullet {
			threat["Bullet"] = threat["Bullet"] + 1
		}
		if Game.NoMansLand[i].Mask {
			threat["Mask"] = threat["Mask"] + 1
		}
		if Game.NoMansLand[i].Night {
			threat["Night"] = threat["Night"] + 1
		}
		if Game.NoMansLand[i].Snow {
			threat["Snow"] = threat["Snow"] + 1
		}
		if Game.NoMansLand[i].Whistle {
			threat["Whistle"] = threat["Whistle"] + 1
		}
	}
	for i := range Players {
		for j := range Players[i].threat {
			if Players[i].WithDraw == false && Players[i].threat[j] != "HardKnock" {
				threat[Players[i].threat[j]] = threat[Players[i].threat[j]] + 1
			}
		}
		hardKnockLimit, _ := strconv.Atoi(os.Getenv("hardKnockLimit"))
		if len(Players[i].threat) >= hardKnockLimit {
			this.loseGame("戰友死亡")
			return true
		}
	}
	log.Print("當前場上威脅")
	log.Println(threat)
	for _, v := range threat {
		threatLimit, _ := strconv.Atoi(os.Getenv("threatLimit"))
		if v >= threatLimit {
			return true
		}
	}
	return false
}

func (this *tRound) loseGame(reason string) {
	this.loseReason = reason
	this.Status.Event("LoseGame")
}
