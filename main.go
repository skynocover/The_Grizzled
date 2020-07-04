package main

import (
	"encoding/json"
	"flag"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	//"fmt"
	"grizzled/database"
	"grizzled/dotenv"
	. "grizzled/game"
	"log"
	"os"
	"strconv"
)

var (
	initdb bool
	err    error
)

func init() {
	flag.BoolVar(&initdb, "initdb", false, "Init the database")
}

type Request struct {
	Order  string `json:"order"`
	Choose string `json:"choose"`
}

func main() {
	flag.Parse()
	{ //基本載入設定
		if err := dotenv.Config(); err != nil {
			return
		} //設定檔載入

		database.DB, err = gorm.Open("sqlite3", os.Getenv("dbroute"))
		if err != nil {
			log.Fatal(err)
			return
		}
		defer database.DB.Close()

		if initdb == true {
			database.InitCard()
		}
	}

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			var request Request
			json.Unmarshal(msg.Body, &request)

			receive := request.Order
			if receive == "login" {
				if Game.NewPlayer(nsConn.Conn.ID(), request.Choose) {
					Game.Stage = "Waiting"
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
					nsConn.Conn.Server().Broadcast(nsConn, msg)
				} else {
					Game.Stage = "ID重複"
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
				}
			} else {
				player := &Players[GetOrder(nsConn.Conn.ID())]

				switch receive {

				case "restart": //重新開始遊戲
					Game.InitGame()
					//Game.NoMansLand = []database.Card{}
					Game.Stage = "DrawCard"

					msg.Body = player.Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
					nsConn.Conn.Server().Broadcast(nsConn, msg)

				case "draw": //每個人抽牌
					Round.Draw(player)

					msg.Body = player.Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
					nsConn.Conn.Server().Broadcast(nsConn, msg)

				case "newRound": //新的回合
					num, _ := strconv.Atoi(request.Choose)
					Round.Init(num)

					msg.Body = player.Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
					nsConn.Conn.Server().Broadcast(nsConn, msg)

				case "playCard":
					num, _ := strconv.Atoi(request.Choose)
					if Round.PlayCard(player, num) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						msg.Body = Game.Render()
						nsConn.Conn.Write(msg)
						nsConn.Conn.Server().Broadcast(nsConn, msg)
					}

				case "heroUse":
					if Round.PlayHero(player) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						msg.Body = Game.Render()
						nsConn.Conn.Write(msg)
						nsConn.Conn.Server().Broadcast(nsConn, msg)
					}
				case "luckyClover":
					choose, _ := strconv.Atoi(request.Choose)

					if Round.HeroUse(player, choose) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						msg.Body = Game.Render()
						nsConn.Conn.Write(msg)
						nsConn.Conn.Server().Broadcast(nsConn, msg)
					}

				case "speech":
					if Round.Speech(player, request.Choose) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						msg.Body = Game.Render()
						nsConn.Conn.Write(msg)
						nsConn.Conn.Server().Broadcast(nsConn, msg)
					}

				case "speechCard":
					choose, _ := strconv.Atoi(request.Choose)

					if Round.SpeechCard(player, choose) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						if Round.Status.Current() == "Mission" {
							msg.Body = Game.Render()
							nsConn.Conn.Write(msg)
							nsConn.Conn.Server().Broadcast(nsConn, msg)
						}

					}

				case "support": //點選支援
					choose, _ := strconv.Atoi(request.Choose)

					if Round.Support(player, choose) {
						msg.Body = player.Render()
						nsConn.Conn.Write(msg)
						msg.Body = Game.Render()
						nsConn.Conn.Write(msg)
						nsConn.Conn.Server().Broadcast(nsConn, msg)
					}

				case "supportEnd": //支援結束後通知支援的結果
					Round.SupportEnd(player, request.Choose)

					msg.Body = player.Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
					nsConn.Conn.Server().Broadcast(nsConn, msg)

				default:

				}
			}
			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server", c.ID())
	}

	app := iris.New()
	app.RegisterView(iris.HTML("./v", ".html"))
	app.HandleDir("/js", "./v/js") // serve our custom javascript code.
	app.HandleDir("/", "./v")
	app.Get("/my_endpoint", websocket.Handler(ws))

	//主頁
	app.Get("/", func(ctx iris.Context) {
		//綁定數據
		ctx.ViewData("Host", os.Getenv("host")+":"+os.Getenv("port"))
		// 渲染視圖文件: ./v/index.html
		ctx.View("index.html")

	})

	app.Run(iris.Addr(":" + os.Getenv("port")))
}
